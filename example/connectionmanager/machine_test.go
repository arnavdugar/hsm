package connectionmanager_test

import (
	"errors"
	"testing"

	"github.com/arnavdugar/hsm/example/connectionmanager"
	"github.com/arnavdugar/hsm/runtime"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStateDisconnectedHandleActionConnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	mockHandler.EXPECT().StartSession().Return(nil)
	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateDisconnected, connectionmanager.ActionConnect, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, connectionmanager.StateConnecting, destination)
}

func TestStateDisconnectedHandleActionConnectHandlerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	handlerErr := errors.New("handler failed")

	mockHandler.EXPECT().StartSession().Return(handlerErr)
	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateDisconnected, connectionmanager.ActionConnect, struct{}{})

	assert.ErrorIs(t, err, handlerErr)
	assert.Equal(t, connectionmanager.StateDisconnected, destination)
}

func TestStateConnectingHandleActionFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateConnecting, connectionmanager.ActionFailure, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, connectionmanager.StateBackoff, destination)
}

func TestStateConnectingHandleActionTransportReady(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateConnecting, connectionmanager.ActionTransportReady, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, connectionmanager.StateAuthenticating, destination)
}

func TestStateAuthenticatingHandleActionFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateAuthenticating, connectionmanager.ActionFailure, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, connectionmanager.StateBackoff, destination)
}

func TestStateAuthenticatingHandleActionAuthenticated(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	mockHandler.EXPECT().ResetRetries().Return(nil)
	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateAuthenticating, connectionmanager.ActionAuthenticated, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, connectionmanager.StateConnected, destination)
}

func TestStateConnectedHandleActionConnectionLost(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateConnected, connectionmanager.ActionConnectionLost, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, connectionmanager.StateBackoff, destination)
}

func TestStateConnectedHandleActionDisconnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateConnected, connectionmanager.ActionDisconnect, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, connectionmanager.StateDisconnected, destination)
}

func TestStateBackoffHandleActionDisconnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateBackoff, connectionmanager.ActionDisconnect, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, connectionmanager.StateDisconnected, destination)
}

func TestStateBackoffHandleActionTimeoutCanRetry(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	gomock.InOrder(
		mockHandler.EXPECT().CanRetry().Return(true, nil),
		mockHandler.EXPECT().RecordRetry().Return(nil),
	)
	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateBackoff, connectionmanager.ActionTimeout, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, connectionmanager.StateConnecting, destination)
}

func TestStateBackoffHandleActionTimeoutBudgetExhausted(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	mockHandler.EXPECT().CanRetry().Return(false, nil)
	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateBackoff, connectionmanager.ActionTimeout, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, connectionmanager.StateFailed, destination)
}

func TestStateBackoffHandleActionTimeoutGuardError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	guardErr := errors.New("guard failed")

	mockHandler.EXPECT().CanRetry().Return(false, guardErr)
	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateBackoff, connectionmanager.ActionTimeout, struct{}{})

	assert.ErrorIs(t, err, guardErr)
	assert.Equal(t, connectionmanager.StateBackoff, destination)
}

func TestStateBackoffHandleActionTimeoutHandlerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	handlerErr := errors.New("handler failed")

	gomock.InOrder(
		mockHandler.EXPECT().CanRetry().Return(true, nil),
		mockHandler.EXPECT().RecordRetry().Return(handlerErr),
	)
	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateBackoff, connectionmanager.ActionTimeout, struct{}{})

	assert.ErrorIs(t, err, handlerErr)
	assert.Equal(t, connectionmanager.StateBackoff, destination)
}

func TestStateFailedHandleActionConnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := connectionmanager.NewMockHandler(ctrl)

	destination, err := connectionmanager.Handle(
		mockHandler, connectionmanager.StateFailed, connectionmanager.ActionConnect, struct{}{})

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, connectionmanager.StateFailed, destination)
}
