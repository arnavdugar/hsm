package parser

import (
	"fmt"
	"io"

	"github.com/arnavdugar/hsm/codegen/config"
	"github.com/arnavdugar/hsm/codegen/util/orderedset"
	"gopkg.in/yaml.v3"
)

type Machine struct {
	config.Machine
	ActionsMap map[string]*ActionData
	StatesMap  map[string]*config.State
}

type ActionData struct {
	Action      *config.Action
	Guards      orderedset.OrderedSet[string]
	Transitions orderedset.OrderedSet[string]
}

func Parse(reader io.Reader) (*Machine, error) {
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	machine := Machine{
		ActionsMap: map[string]*ActionData{},
		StatesMap:  map[string]*config.State{},
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
		machine.ActionsMap[action.Name] = &ActionData{
			Action:      action,
			Guards:      orderedset.Create[string](),
			Transitions: orderedset.Create[string](),
		}

		if action.Symbol == "" {
			action.Symbol = fmt.Sprintf("Action%s", action.Name)
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
		machine.StatesMap[state.Name] = state

		if state.Symbol == "" {
			state.Symbol = fmt.Sprintf("State%s", state.Name)
		}
	}

	for _, state := range machine.Machine.States.Values {
		for _, action := range state.TransitionActions {
			actionsMap, ok := machine.ActionsMap[action.Action]
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
					actionsMap.Guards.Add(transition.Guard)
				}

				if transition.Transition != "" {
					actionsMap.Transitions.Add(transition.Transition)
				}
			}
		}
	}

	return &machine, nil
}
