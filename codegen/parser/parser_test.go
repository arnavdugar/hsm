package parser_test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/arnavdugar/hsm/codegen/parser"
	"github.com/stretchr/testify/assert"
)

//go:embed test/cyclic_groups.yaml
var CyclicGroupsConfig string

func TestParseCyclicGroups(t *testing.T) {
	reader := strings.NewReader(CyclicGroupsConfig)

	_, err := parser.Parse(reader)

	assert.Error(t, err)
	assert.Equal(t, "groups declaration contains a cycle", err.Error())
}

//go:embed test/duplicate_action.yaml
var DuplicateActionConfig string

func TestParseDuplicateAction(t *testing.T) {
	reader := strings.NewReader(DuplicateActionConfig)

	_, err := parser.Parse(reader)

	assert.Error(t, err)
	assert.Equal(t, "duplicate action name: ActionName", err.Error())
}

//go:embed test/duplicate_state.yaml
var DuplicateStateConfig string

func TestParseDuplicateState(t *testing.T) {
	reader := strings.NewReader(DuplicateStateConfig)

	_, err := parser.Parse(reader)

	assert.Error(t, err)
	assert.Equal(t, "duplicate state name: StateName", err.Error())
}

//go:embed test/groups.yaml
var GroupsConfig string

func TestParseGroups(t *testing.T) {
	reader := strings.NewReader(GroupsConfig)

	machine, err := parser.Parse(reader)

	assert.NoError(t, err)
	if assert.Len(t, machine.StatesMap["StateA"].Groups, 4) {
		assert.Equal(t, "GroupC", machine.StatesMap["StateA"].Groups[0].Name)
		assert.Equal(t, "GroupD", machine.StatesMap["StateA"].Groups[1].Name)
		assert.Equal(t, "GroupA", machine.StatesMap["StateA"].Groups[2].Name)
		assert.Equal(t, "GroupB", machine.StatesMap["StateA"].Groups[3].Name)
	}
	if assert.Len(t, machine.StatesMap["StateB"].Groups, 4) {
		assert.Equal(t, "GroupC", machine.StatesMap["StateB"].Groups[0].Name)
		assert.Equal(t, "GroupE", machine.StatesMap["StateB"].Groups[1].Name)
		assert.Equal(t, "GroupA", machine.StatesMap["StateB"].Groups[2].Name)
		assert.Equal(t, "GroupB", machine.StatesMap["StateB"].Groups[3].Name)
	}
}

//go:embed test/unknown_destination.yaml
var UnknownDestinationConfig string

func TestParseUnknownDestination(t *testing.T) {
	reader := strings.NewReader(UnknownDestinationConfig)

	_, err := parser.Parse(reader)

	assert.Error(t, err)
	assert.Equal(t, `unknown destination state "UnknownState" in "ActionName" transition from state "StateName"`, err.Error())
}

//go:embed test/unknown_subgroup.yaml
var UnknownSubgroupConfig string

func TestParseUnknownSubgroup(t *testing.T) {
	reader := strings.NewReader(UnknownSubgroupConfig)

	_, err := parser.Parse(reader)

	assert.Error(t, err)
	assert.Equal(t, `unknown subgroup "UnknownGroup" in group "GroupA"`, err.Error())
}

//go:embed test/unknown_substate.yaml
var UnknownSubstateConfig string

func TestParseUnknownSubstate(t *testing.T) {
	reader := strings.NewReader(UnknownSubstateConfig)

	_, err := parser.Parse(reader)

	assert.Error(t, err)
	assert.Equal(t, `unknown substate "UnknownState" in group "GroupA"`, err.Error())
}

//go:embed test/unknown_transition_action.yaml
var UnknownTransitionActionConfig string

func TestParseUnknownTransitionAction(t *testing.T) {
	reader := strings.NewReader(UnknownTransitionActionConfig)

	_, err := parser.Parse(reader)

	assert.Error(t, err)
	assert.Equal(t, `unknown action "UnknownActionName" in transitions for state "StateName"`, err.Error())
}
