package combinationlock_test

import (
	"errors"
	"testing"

	"github.com/arnavdugar/hsm/example/combinationlock"
	"github.com/arnavdugar/hsm/runtime"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStateLockedHandleActionDigitCorrectDigit(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := combinationlock.NewMockHandler(ctrl)

	mockHandler.EXPECT().IsFirstDigit(4).Return(true, nil)
	destination, err := combinationlock.Handle(
		mockHandler, combinationlock.StateLocked, combinationlock.ActionDigit, 4)

	assert.NoError(t, err)
	assert.Equal(t, combinationlock.StateFirstDigit, destination)
}

func TestStateLockedHandleActionDigitIncorrectDigit(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := combinationlock.NewMockHandler(ctrl)

	mockHandler.EXPECT().IsFirstDigit(0).Return(false, nil)
	destination, err := combinationlock.Handle(
		mockHandler, combinationlock.StateLocked, combinationlock.ActionDigit, 0)

	assert.NoError(t, err)
	assert.Equal(t, combinationlock.StateLocked, destination)
}

func TestStateLockedHandleActionDigitGuardError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := combinationlock.NewMockHandler(ctrl)

	guardErr := errors.New("guard failed")

	mockHandler.EXPECT().IsFirstDigit(4).Return(false, guardErr)
	destination, err := combinationlock.Handle(
		mockHandler, combinationlock.StateLocked, combinationlock.ActionDigit, 4)

	assert.ErrorIs(t, err, guardErr)
	assert.Equal(t, combinationlock.StateLocked, destination)
}

func TestStateFirstDigitHandleActionDigitCorrectDigit(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := combinationlock.NewMockHandler(ctrl)

	mockHandler.EXPECT().IsSecondDigit(2).Return(true, nil)
	destination, err := combinationlock.Handle(
		mockHandler, combinationlock.StateFirstDigit, combinationlock.ActionDigit, 2)

	assert.NoError(t, err)
	assert.Equal(t, combinationlock.StateSecondDigit, destination)
}

func TestStateFirstDigitHandleActionDigitIncorrectDigit(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := combinationlock.NewMockHandler(ctrl)

	mockHandler.EXPECT().IsSecondDigit(0).Return(false, nil)
	destination, err := combinationlock.Handle(
		mockHandler, combinationlock.StateFirstDigit, combinationlock.ActionDigit, 0)

	assert.NoError(t, err)
	assert.Equal(t, combinationlock.StateLocked, destination)
}

func TestStateSecondDigitHandleActionDigitCorrectDigit(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := combinationlock.NewMockHandler(ctrl)

	mockHandler.EXPECT().IsThirdDigit(7).Return(true, nil)
	destination, err := combinationlock.Handle(
		mockHandler, combinationlock.StateSecondDigit, combinationlock.ActionDigit, 7)

	assert.NoError(t, err)
	assert.Equal(t, combinationlock.StateUnlocked, destination)
}

func TestStateSecondDigitHandleActionDigitIncorrectDigit(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := combinationlock.NewMockHandler(ctrl)

	mockHandler.EXPECT().IsThirdDigit(0).Return(false, nil)
	destination, err := combinationlock.Handle(
		mockHandler, combinationlock.StateSecondDigit, combinationlock.ActionDigit, 0)

	assert.NoError(t, err)
	assert.Equal(t, combinationlock.StateLocked, destination)
}

func TestStateUnlockedHandleActionLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := combinationlock.NewMockHandler(ctrl)

	destination, err := combinationlock.Handle(
		mockHandler, combinationlock.StateUnlocked, combinationlock.ActionLock, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, combinationlock.StateLocked, destination)
}

func TestStateUnlockedHandleActionDigit(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := combinationlock.NewMockHandler(ctrl)

	destination, err := combinationlock.Handle(
		mockHandler, combinationlock.StateUnlocked, combinationlock.ActionDigit, 4)

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, combinationlock.StateUnlocked, destination)
}

func TestStateFirstDigitHandleActionDigitInvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := combinationlock.NewMockHandler(ctrl)

	destination, err := combinationlock.Handle(
		mockHandler, combinationlock.StateFirstDigit, combinationlock.ActionDigit, "2")

	assert.ErrorIs(t, err, runtime.ErrInvalidActionDataType)
	assert.Equal(t, combinationlock.StateFirstDigit, destination)
}
