package atm_test

import (
	"testing"

	"github.com/arnavdugar/hsm/example/atm"
	"github.com/arnavdugar/hsm/runtime"
	"github.com/stretchr/testify/assert"
)

func TestStateMenuHandleActionCheckBalance(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateMenu, atm.ActionCheckBalance, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, atm.StateBalance, destination)
}

func TestStateBalanceHandleActionDone(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateBalance, atm.ActionDone, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, atm.StateMenu, destination)
}

func TestStateMenuHandleActionWithdraw(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateMenu, atm.ActionWithdraw, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, atm.StateEnterAmount, destination)
}

func TestStateEnterAmountHandleActionSubmit(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateEnterAmount, atm.ActionSubmit, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, atm.StateAuthorizeWithdrawal, destination)
}

func TestStateAuthorizeWithdrawalHandleActionApprove(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateAuthorizeWithdrawal, atm.ActionApprove, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, atm.StateDispenseCash, destination)
}

func TestStateDispenseCashHandleActionDone(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateDispenseCash, atm.ActionDone, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, atm.StateMenu, destination)
}

func TestStateMenuHandleActionDeposit(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateMenu, atm.ActionDeposit, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, atm.StateAcceptDeposit, destination)
}

func TestStateAcceptDepositHandleActionDone(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateAcceptDeposit, atm.ActionDone, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, atm.StateMenu, destination)
}

func TestStateMenuHandleActionExit(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateMenu, atm.ActionExit, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, atm.StateEnded, destination)
}

func TestStateAuthorizeWithdrawalHandleActionDecline(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateAuthorizeWithdrawal, atm.ActionDecline, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, atm.StateMenu, destination)
}

func TestStateEnterAmountHandleActionCancel(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateEnterAmount, atm.ActionCancel, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, atm.StateMenu, destination)
}

func TestStateAcceptDepositHandleActionCancel(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateAcceptDeposit, atm.ActionCancel, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, atm.StateMenu, destination)
}

func TestStateMenuHandleActionApprove(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateMenu, atm.ActionApprove, struct{}{})

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, atm.StateMenu, destination)
}

func TestStateEndedHandleActionWithdraw(t *testing.T) {
	destination, err := atm.Handle(
		nil, atm.StateEnded, atm.ActionWithdraw, struct{}{})

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, atm.StateEnded, destination)
}
