package orderfulfillment_test

import (
	"testing"

	"github.com/arnavdugar/hsm/example/orderfulfillment"
	"github.com/arnavdugar/hsm/runtime"
	"github.com/stretchr/testify/assert"
)

func TestStatePaidHandleActionShip(t *testing.T) {
	destination, err := orderfulfillment.Handle(
		nil, orderfulfillment.StatePaid, orderfulfillment.ActionShip, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, orderfulfillment.StatePacking, destination)
}

func TestStatePackingHandleActionDispatch(t *testing.T) {
	destination, err := orderfulfillment.Handle(
		nil, orderfulfillment.StatePacking, orderfulfillment.ActionDispatch, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, orderfulfillment.StateShipped, destination)
}

func TestStateShippedHandleActionDeliver(t *testing.T) {
	destination, err := orderfulfillment.Handle(
		nil, orderfulfillment.StateShipped, orderfulfillment.ActionDeliver, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, orderfulfillment.StateCompleted, destination)
}

func TestStatePaidHandleActionPickUp(t *testing.T) {
	destination, err := orderfulfillment.Handle(
		nil, orderfulfillment.StatePaid, orderfulfillment.ActionPickUp, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, orderfulfillment.StateReadyForPickup, destination)
}

func TestStateReadyForPickupHandleActionCollect(t *testing.T) {
	destination, err := orderfulfillment.Handle(
		nil, orderfulfillment.StateReadyForPickup, orderfulfillment.ActionCollect, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, orderfulfillment.StateCompleted, destination)
}

func TestStatePaidHandleActionCancel(t *testing.T) {
	destination, err := orderfulfillment.Handle(
		nil, orderfulfillment.StatePaid, orderfulfillment.ActionCancel, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, orderfulfillment.StateCancelled, destination)
}

func TestStatePackingHandleActionCancel(t *testing.T) {
	destination, err := orderfulfillment.Handle(
		nil, orderfulfillment.StatePacking, orderfulfillment.ActionCancel, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, orderfulfillment.StateCancelled, destination)
}

func TestStateReadyForPickupHandleActionCancel(t *testing.T) {
	destination, err := orderfulfillment.Handle(
		nil, orderfulfillment.StateReadyForPickup, orderfulfillment.ActionCancel, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, orderfulfillment.StateCancelled, destination)
}

func TestStatePaidHandleActionDeliver(t *testing.T) {
	destination, err := orderfulfillment.Handle(
		nil, orderfulfillment.StatePaid, orderfulfillment.ActionDeliver, struct{}{})

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, orderfulfillment.StatePaid, destination)
}

func TestStateCompletedHandleActionCancel(t *testing.T) {
	destination, err := orderfulfillment.Handle(
		nil, orderfulfillment.StateCompleted, orderfulfillment.ActionCancel, struct{}{})

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, orderfulfillment.StateCompleted, destination)
}

func TestStateCancelledHandleActionShip(t *testing.T) {
	destination, err := orderfulfillment.Handle(
		nil, orderfulfillment.StateCancelled, orderfulfillment.ActionShip, struct{}{})

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, orderfulfillment.StateCancelled, destination)
}
