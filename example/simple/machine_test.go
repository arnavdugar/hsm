package simple_test

import (
	"testing"

	"github.com/arnavdugar/hsm/example/simple"
	"github.com/stretchr/testify/assert"
)

func TestStateAHandleActionNext(t *testing.T) {
	destination, err := simple.Handle(
		nil, simple.StateA, simple.ActionNext, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, simple.StateB, destination)
}

func TestStateBHandleActionNext(t *testing.T) {
	destination, err := simple.Handle(
		nil, simple.StateB, simple.ActionNext, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, simple.StateC, destination)
}

func TestStateCHandleActionNext(t *testing.T) {
	destination, err := simple.Handle(
		nil, simple.StateC, simple.ActionNext, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, simple.StateA, destination)
}
