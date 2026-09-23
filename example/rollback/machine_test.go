package rollback_test

import (
	"testing"

	"github.com/arnavdugar/hsm/example/rollback"
	"github.com/arnavdugar/hsm/runtime"
	"github.com/stretchr/testify/assert"
)

func TestStateReserveHandleActionReserved(t *testing.T) {
	destination, err := rollback.Handle(
		nil, rollback.StateReserve, rollback.ActionReserved, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, rollback.StateCharge, destination)
}

func TestStateChargeHandleActionCharged(t *testing.T) {
	destination, err := rollback.Handle(
		nil, rollback.StateCharge, rollback.ActionCharged, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, rollback.StateFulfill, destination)
}

func TestStateFulfillHandleActionFulfilled(t *testing.T) {
	destination, err := rollback.Handle(
		nil, rollback.StateFulfill, rollback.ActionFulfilled, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, rollback.StateCompleted, destination)
}

func TestStateReserveHandleActionFailure(t *testing.T) {
	destination, err := rollback.Handle(
		nil, rollback.StateReserve, rollback.ActionFailure, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, rollback.StateFailed, destination)
}

func TestStateChargeHandleActionFailure(t *testing.T) {
	destination, err := rollback.Handle(
		nil, rollback.StateCharge, rollback.ActionFailure, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, rollback.StateReleaseReservation, destination)
}

func TestStateReleaseReservationHandleActionReleased(t *testing.T) {
	destination, err := rollback.Handle(
		nil, rollback.StateReleaseReservation, rollback.ActionReleased, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, rollback.StateFailed, destination)
}

func TestStateFulfillHandleActionFailure(t *testing.T) {
	destination, err := rollback.Handle(
		nil, rollback.StateFulfill, rollback.ActionFailure, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, rollback.StateRefund, destination)
}

func TestStateRefundHandleActionRefunded(t *testing.T) {
	destination, err := rollback.Handle(
		nil, rollback.StateRefund, rollback.ActionRefunded, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, rollback.StateReleaseReservation, destination)
}

func TestStateReserveHandleActionFulfilled(t *testing.T) {
	destination, err := rollback.Handle(
		nil, rollback.StateReserve, rollback.ActionFulfilled, struct{}{})

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, rollback.StateReserve, destination)
}

func TestStateCompletedHandleActionFailure(t *testing.T) {
	destination, err := rollback.Handle(
		nil, rollback.StateCompleted, rollback.ActionFailure, struct{}{})

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, rollback.StateCompleted, destination)
}

func TestStateFailedHandleActionReserved(t *testing.T) {
	destination, err := rollback.Handle(
		nil, rollback.StateFailed, rollback.ActionReserved, struct{}{})

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, rollback.StateFailed, destination)
}
