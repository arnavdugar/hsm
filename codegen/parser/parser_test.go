package parser_test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/arnavdugar/hsm/codegen/parser"
	"github.com/stretchr/testify/assert"
)

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

//go:embed test/unknown_destination.yaml
var UnknownDestinationConfig string

func TestParseUnknownDestination(t *testing.T) {
	reader := strings.NewReader(UnknownDestinationConfig)

	_, err := parser.Parse(reader)

	assert.Error(t, err)
	assert.Equal(t, `unknown destination state "UnknownState" in "ActionName" transition from state "StateName"`, err.Error())
}

//go:embed test/unknown_transition_action.yaml
var UnknownTransitionActionConfig string

func TestParseUnknownTransitionAction(t *testing.T) {
	reader := strings.NewReader(UnknownTransitionActionConfig)

	_, err := parser.Parse(reader)

	assert.Error(t, err)
	assert.Equal(t, `unknown action "UnknownActionName" in transitions for state "StateName"`, err.Error())
}
