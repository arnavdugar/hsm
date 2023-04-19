package parser

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/arnavdugar/hsm/codegen/config"
	"github.com/arnavdugar/hsm/codegen/util/orderedset"
	"gopkg.in/yaml.v3"
)

type Machine struct {
	config.Machine
	ActionsMap       map[string]*ActionData
	BoundaryHandlers orderedset.OrderedSet[string]
	GroupsMap        map[string]*GroupData
	StatesMap        map[string]*StateData
}

type ActionData struct {
	Action         *config.Action
	Guards         orderedset.OrderedSet[string]
	Transitions    orderedset.OrderedSet[string]
	TransitionData map[string][]*TransitionData
}

type TransitionData struct {
	EnterHandlers []string
	ExitHandlers  []string
	Guard         string
	Transition    string
}

type GroupData struct {
	*config.Group
	ancestors  map[string]struct{}
	Children   []*GroupData
	Index      int
	parents    []*GroupData
	statesData []*StateData
}

type StateData struct {
	*config.State
	ancestors map[string]struct{}
	Groups    []*GroupData
}

func Parse(reader io.Reader) (*Machine, error) {
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	machine := Machine{
		ActionsMap:       map[string]*ActionData{},
		BoundaryHandlers: orderedset.Create[string](),
		GroupsMap:        map[string]*GroupData{},
		StatesMap:        map[string]*StateData{},
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

		_, ok := machine.ActionsMap[action.Name]
		if ok {
			return nil, fmt.Errorf("duplicate action name: %s", action.Name)
		}

		if action.Symbol == "" {
			action.Symbol = fmt.Sprintf("Action%s", action.Name)
		}

		machine.ActionsMap[action.Name] = &ActionData{
			Action:         action,
			Guards:         orderedset.Create[string](),
			Transitions:    orderedset.Create[string](),
			TransitionData: map[string]*TransitionData{},
		}
	}

	states := &machine.Machine.States
	if states.Type == "" {
		states.Type = "StateType"
	}

	for index := range states.Values {
		state := &states.Values[index]
		_, ok := machine.StatesMap[state.Name]
		if ok {
			return nil, fmt.Errorf("duplicate state name: %s", state.Name)
		}

		if state.Symbol == "" {
			state.Symbol = fmt.Sprintf("State%s", state.Name)
		}

		if state.Enter != "" {
			machine.BoundaryHandlers.Add(state.Enter)
		}
		if state.Exit != "" {
			machine.BoundaryHandlers.Add(state.Exit)
		}

		machine.StatesMap[state.Name] = &StateData{
			State:     state,
			ancestors: map[string]struct{}{},
			Groups:    []*GroupData{},
		}
	}

	for _, state := range machine.Machine.States.Values {
		for _, action := range state.TransitionActions {
			actionData, ok := machine.ActionsMap[action.Action]
			if !ok {
				return nil, fmt.Errorf(
					`unknown action "%s" in transitions for state "%s"`,
					action.Action, state.Name)
			}

			for _, transition := range action.Transitions {
				_, ok := machine.StatesMap[transition.Destination]
				if !ok {
					return nil, fmt.Errorf(
						`unknown destination state "%s" in "%s" transition from state "%s"`,
						transition.Destination, action.Action, state.Name)
				}

				if transition.Guard != "" {
					actionData.Guards.Add(transition.Guard)
				}
				if transition.Transition != "" {
					actionData.Transitions.Add(transition.Transition)
				}
			}
		}
	}

	for index := range machine.Groups {
		group := &machine.Groups[index]

		machine.GroupsMap[group.Name] = &GroupData{
			Group:      group,
			ancestors:  map[string]struct{}{},
			Children:   []*GroupData{},
			Index:      index,
			parents:    []*GroupData{},
			statesData: []*StateData{},
		}
	}

	for _, groupData := range machine.GroupsMap {
		for _, child := range groupData.Groups {
			childData, ok := machine.GroupsMap[child]
			if !ok {
				return nil, fmt.Errorf(
					`unknown subgroup "%s" in group "%s"`, child, groupData.Name)
			}

			groupData.Children = append(groupData.Children, childData)
			childData.parents = append(childData.parents, groupData)
		}

		for _, state := range groupData.States {
			stateData, ok := machine.StatesMap[state]
			if !ok {
				return nil, fmt.Errorf(
					`unknown substate "%s" in group "%s"`, state, groupData.Name)
			}
			groupData.statesData = append(groupData.statesData, stateData)
		}
	}

	pendingGroups := map[string]int{}
	currentGroups := []*GroupData{}

	for _, group := range machine.Groups {
		groupData := machine.GroupsMap[group.Name]

		if len(groupData.parents) == 0 {
			currentGroups = append(currentGroups, groupData)
		} else {
			pendingGroups[group.Name] = len(groupData.parents)
		}
	}

	for len(currentGroups) > 0 {
		current := currentGroups[0]
		currentGroups = currentGroups[1:]

		for _, child := range current.Children {
			child.ancestors[current.Name] = struct{}{}
			for ancestor := range current.ancestors {
				child.ancestors[ancestor] = struct{}{}
			}

			if pendingGroups[child.Name] == 1 {
				delete(pendingGroups, child.Name)
				currentGroups = append(currentGroups, child)
			} else {
				pendingGroups[child.Name] -= 1
			}
		}

		for _, state := range current.statesData {
			state.ancestors[current.Name] = struct{}{}
			for ancestor := range current.ancestors {
				state.ancestors[ancestor] = struct{}{}
			}
		}
	}

	if len(pendingGroups) > 0 {
		return nil, errors.New("groups declaration contains a cycle")
	}

	for _, stateData := range machine.StatesMap {
		for group := range stateData.ancestors {
			stateData.Groups = append(stateData.Groups, machine.GroupsMap[group])
		}
		sort.Slice(stateData.Groups, func(i, j int) bool {
			return stateData.Groups[i].Index < stateData.Groups[j].Index
		})
	}

	for _, stateData := range machine.StatesMap {
		for _, transitionAction := range stateData.TransitionActions {
			actionData := machine.ActionsMap[transitionAction.Action]
			transitionData := actionData.TransitionData[stateData.Name]
			for _, transition := range transitionAction.Transitions {
				exitGroups := stateData.Groups
				enterGroups := machine.StatesMap[transition.Destination].Groups

				exitIndex, enterIndex := 0, 0
				for exitIndex < len(exitGroups) || enterIndex < len(enterGroups) {
					exitGroupIndex := exitGroups[exitIndex].Index
					enterGroupIndex := enterGroups[enterIndex].Index
					switch {
					case exitGroupIndex < enterGroupIndex:

					case exitGroupIndex > enterGroupIndex:

					default:
						// TODO: Handle
					}
				}

				transitionData = append(transitionData, &TransitionData{
					EnterHandlers: []string{},
					ExitHandlers:  []string{},
					Guard:         transition.Guard,
					Transition:    transition.Transition,
				})
			}
		}
	}

	return &machine, nil
}
