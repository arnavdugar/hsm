package parser

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/arnavdugar/hsm/codegen/config"
	"github.com/arnavdugar/hsm/codegen/util/orderedset"
	"go.yaml.in/yaml/v3"
)

type Machine struct {
	config.Machine
	ActionsMap map[string]*ActionData
	GroupsMap  map[string]*GroupData
	StatesMap  map[string]*StateData
}

type ActionData struct {
	Action      *config.Action
	Guards      orderedset.OrderedSet[string]
	Transitions orderedset.OrderedSet[string]
}

type GroupData struct {
	Group       *config.Group
	Descendants orderedset.OrderedSet[string]
	States      orderedset.OrderedSet[string]
	Initial     *config.State
}

type StateData struct {
	State *config.State
	// All containing groups, in declaration order.
	Groups orderedset.OrderedSet[string]
	// Effective transitions for each action, in evaluation order.
	Transitions map[string][]TransitionData
}

type TransitionData struct {
	Destination *config.State
	Guard       string
	Transition  string
	// Group boundaries in execution order, with duplicate paths removed.
	ExitGroups  []*config.Group
	EnterGroups []*config.Group
}

func Parse(reader io.Reader) (*Machine, error) {
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	machine := Machine{
		ActionsMap: map[string]*ActionData{},
		GroupsMap:  map[string]*GroupData{},
		StatesMap:  map[string]*StateData{},
	}

	err := decoder.Decode(&machine.Machine)
	if err != nil {
		return nil, fmt.Errorf("error parsing input file: %v", err)
	}

	actions := &machine.Machine.Actions
	if actions.Type == "" {
		actions.Type = "ActionType"
	}

	for index := range actions.Values {
		action := &actions.Values[index]
		if action.Name == "" {
			return nil, fmt.Errorf("action name must not be empty")
		}
		_, ok := machine.ActionsMap[action.Name]
		if ok {
			return nil, fmt.Errorf("duplicate action name: %s", action.Name)
		}
		machine.ActionsMap[action.Name] = &ActionData{
			Action:      action,
			Guards:      orderedset.Create[string](),
			Transitions: orderedset.Create[string](),
		}

		if action.Symbol == "" {
			action.Symbol = fmt.Sprintf("Action%s", action.Name)
		}
	}

	for index := range machine.Machine.Groups.Values {
		group := &machine.Machine.Groups.Values[index]
		if group.Name == "" {
			return nil, fmt.Errorf("group name must not be empty")
		}
		_, ok := machine.GroupsMap[group.Name]
		if ok {
			return nil, fmt.Errorf("duplicate group name: %s", group.Name)
		}
		machine.GroupsMap[group.Name] = &GroupData{
			Group: group,
		}
	}

	states := &machine.Machine.States
	if states.Type == "" {
		states.Type = "StateType"
	}

	for index := range states.Values {
		state := &states.Values[index]
		if state.Name == "" {
			return nil, fmt.Errorf("state name must not be empty")
		}
		_, ok := machine.StatesMap[state.Name]
		if ok {
			return nil, fmt.Errorf("duplicate state name: %s", state.Name)
		}
		if _, ok := machine.GroupsMap[state.Name]; ok {
			return nil, fmt.Errorf("name used by both a state and a group: %s", state.Name)
		}
		machine.StatesMap[state.Name] = &StateData{
			State:       state,
			Transitions: map[string][]TransitionData{},
		}

		if state.Symbol == "" {
			state.Symbol = fmt.Sprintf("State%s", state.Name)
		}
	}

	if err := machine.resolveGroups(); err != nil {
		return nil, err
	}

	for _, state := range machine.Machine.States.Values {
		err := machine.validateTransitions("state", state.Name, state.TransitionActions)
		if err != nil {
			return nil, err
		}
	}
	for _, group := range machine.Machine.Groups.Values {
		err := machine.validateTransitions("group", group.Name, group.TransitionActions)
		if err != nil {
			return nil, err
		}
	}
	for _, state := range machine.States.Values {
		machine.collectActionMetadata(state.TransitionActions)
	}
	for _, group := range machine.Groups.Values {
		machine.collectActionMetadata(group.TransitionActions)
	}
	if err := machine.resolveTransitions(); err != nil {
		return nil, err
	}

	return &machine, nil
}

func (machine *Machine) resolveGroups() error {
	positions := map[string]int{}
	for index, group := range machine.Groups.Values {
		positions[group.Name] = index
		for _, name := range group.Groups {
			if _, ok := machine.GroupsMap[name]; !ok {
				return fmt.Errorf("unknown group %q in group %q", name, group.Name)
			}
		}
		for _, name := range group.States {
			if _, ok := machine.StatesMap[name]; !ok {
				return fmt.Errorf("unknown state %q in group %q", name, group.Name)
			}
		}
	}

	// Resolve membership through the DAG before checking declaration order, so
	// cycles can be reported with their containment path.
	visiting, resolved := map[string]bool{}, map[string]bool{}
	var visit func(string, []string) error
	visit = func(name string, path []string) error {
		if visiting[name] {
			return fmt.Errorf("group containment cycle: %s", strings.Join(append(path, name), " -> "))
		}
		if resolved[name] {
			return nil
		}
		visiting[name] = true
		group := machine.GroupsMap[name]
		for _, state := range group.Group.States {
			group.States.Add(state)
		}
		for _, name := range group.Group.Groups {
			if err := visit(name, append(path, group.Group.Name)); err != nil {
				return err
			}
			subgroup := machine.GroupsMap[name]
			group.Descendants.Add(name)
			for _, descendant := range subgroup.Descendants.Values() {
				group.Descendants.Add(descendant)
			}
			for _, state := range subgroup.States.Values() {
				group.States.Add(state)
			}
		}
		visiting[name] = false
		resolved[name] = true
		return nil
	}
	for _, group := range machine.Groups.Values {
		if err := visit(group.Name, nil); err != nil {
			return err
		}
	}

	for _, group := range machine.Groups.Values {
		for _, subgroup := range group.Groups {
			if positions[group.Name] > positions[subgroup] {
				return fmt.Errorf("containing groups must be declared before their subgroups: group %q must be declared before group %q", group.Name, subgroup)
			}
		}
		for _, state := range machine.GroupsMap[group.Name].States.Values() {
			machine.StatesMap[state].Groups.Add(group.Name)
		}
	}

	for _, group := range machine.Groups.Values {
		if group.Initial == "" {
			continue
		}
		if err := machine.resolveInitial(machine.GroupsMap[group.Name]); err != nil {
			return err
		}
	}
	return nil
}

func (machine *Machine) resolveInitial(group *GroupData) error {
	if group.Initial != nil {
		return nil
	}
	name := group.Group.Initial
	if name == "" {
		return fmt.Errorf("group %q has no initial value", group.Group.Name)
	}
	if state, ok := machine.StatesMap[name]; ok {
		if !group.States.Contains(name) {
			return fmt.Errorf("initial state %q is not a member of group %q", name, group.Group.Name)
		}
		group.Initial = state.State
		return nil
	}
	if subgroup, ok := machine.GroupsMap[name]; ok {
		if !group.Descendants.Contains(name) {
			return fmt.Errorf("initial group %q is not a subgroup of group %q", name, group.Group.Name)
		}
		// Initial group references follow strict descendants of the validated DAG,
		// so initial resolution cannot cycle.
		if err := machine.resolveInitial(subgroup); err != nil {
			return fmt.Errorf("initial value for group %q: %w", group.Group.Name, err)
		}
		group.Initial = subgroup.Initial
		return nil
	}
	return fmt.Errorf("unknown initial target %q in group %q", name, group.Group.Name)
}

func (machine *Machine) validateTransitions(kind, name string, actions []config.TransitionAction) error {
	for _, action := range actions {
		_, ok := machine.ActionsMap[action.Action]
		if !ok {
			return fmt.Errorf("unknown action %q in transitions for %s %q", action.Action, kind, name)
		}
		for _, transition := range action.Transitions {
			if _, err := machine.destination(transition.Destination); err != nil {
				return fmt.Errorf("%w in %q transition from %s %q", err, action.Action, kind, name)
			}
			for _, external := range transition.External {
				if _, ok := machine.GroupsMap[external]; !ok {
					return fmt.Errorf("unknown external group %q in %q transition from %s %q", external, action.Action, kind, name)
				}
			}
		}
	}
	return nil
}

func (machine *Machine) collectActionMetadata(actions []config.TransitionAction) {
	for _, action := range actions {
		actionData := machine.ActionsMap[action.Action]
		for _, transition := range action.Transitions {
			if transition.Guard != "" {
				actionData.Guards.Add(transition.Guard)
			}
			if transition.Transition != "" {
				actionData.Transitions.Add(transition.Transition)
			}
		}
	}
}

func (machine *Machine) resolveTransitions() error {
	for _, state := range machine.States.Values {
		stateData := machine.StatesMap[state.Name]
		for _, action := range machine.Actions.Values {
			// Copy state rules before appending group rules: the configuration's
			// transition slices must not be changed during expansion.
			candidates := append([]config.Transition(nil), actionTransitions(state.TransitionActions, action.Name)...)
			for _, name := range slices.Backward(stateData.Groups.Values()) {
				group := machine.GroupsMap[name].Group
				candidates = append(candidates, actionTransitions(group.TransitionActions, action.Name)...)
			}
			for _, candidate := range candidates {
				transition, err := machine.resolveTransition(stateData, candidate)
				if err != nil {
					return err
				}
				stateData.Transitions[action.Name] = append(stateData.Transitions[action.Name], transition)
				if candidate.Guard == "" {
					break
				}
			}
		}
	}
	return nil
}

func actionTransitions(actions []config.TransitionAction, name string) []config.Transition {
	for _, action := range actions {
		if action.Action == name {
			return action.Transitions
		}
	}
	return nil
}

func (machine *Machine) resolveTransition(source *StateData, transition config.Transition) (TransitionData, error) {
	destination, err := machine.destination(transition.Destination)
	if err != nil {
		return TransitionData{}, err
	}
	result := TransitionData{
		Destination: destination,
		Guard:       transition.Guard,
		Transition:  transition.Transition,
	}
	external := orderedset.Create[string]()
	for _, name := range transition.External {
		external.Add(name)
		for _, descendant := range machine.GroupsMap[name].Descendants.Values() {
			external.Add(descendant)
		}
	}

	target := machine.StatesMap[destination.Name]
	for _, name := range slices.Backward(source.Groups.Values()) {
		if !target.Groups.Contains(name) || external.Contains(name) {
			result.ExitGroups = append(result.ExitGroups, machine.GroupsMap[name].Group)
		}
	}
	for _, name := range target.Groups.Values() {
		if !source.Groups.Contains(name) || external.Contains(name) {
			result.EnterGroups = append(result.EnterGroups, machine.GroupsMap[name].Group)
		}
	}
	return result, nil
}

func (machine *Machine) destination(name string) (*config.State, error) {
	if state, ok := machine.StatesMap[name]; ok {
		return state.State, nil
	}
	if group, ok := machine.GroupsMap[name]; ok {
		if group.Initial == nil {
			return nil, fmt.Errorf("group %q has no initial value", name)
		}
		return group.Initial, nil
	}
	return nil, fmt.Errorf("unknown destination %q", name)
}
