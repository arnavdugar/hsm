package boundary_test

import (
	"context"
	"errors"
	"testing"

	"github.com/arnavdugar/hsm/example/boundary"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStateAHandleActionNext(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := boundary.NewMockHandler(ctrl)
	ctx := context.Background()

	gomock.InOrder(
		mockHandler.EXPECT().OnExitA(ctx).Return(nil),
		mockHandler.EXPECT().OnEnterB(ctx).Return(nil),
	)
	destination, err := boundary.Handle(
		ctx, mockHandler, boundary.StateA, boundary.ActionNext, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, boundary.StateB, destination)
}

func TestStateAHandleActionNextExitError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := boundary.NewMockHandler(ctrl)
	ctx := context.Background()
	handlerErr := errors.New("exit failed")

	mockHandler.EXPECT().OnExitA(ctx).Return(handlerErr)
	destination, err := boundary.Handle(
		ctx, mockHandler, boundary.StateA, boundary.ActionNext, struct{}{})

	assert.ErrorIs(t, err, handlerErr)
	assert.Equal(t, boundary.StateA, destination)
}

func TestStateAHandleActionNextEnterError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := boundary.NewMockHandler(ctrl)
	ctx := context.Background()
	handlerErr := errors.New("enter failed")

	gomock.InOrder(
		mockHandler.EXPECT().OnExitA(ctx).Return(nil),
		mockHandler.EXPECT().OnEnterB(ctx).Return(handlerErr),
	)
	destination, err := boundary.Handle(
		ctx, mockHandler, boundary.StateA, boundary.ActionNext, struct{}{})

	assert.ErrorIs(t, err, handlerErr)
	assert.Equal(t, boundary.StateA, destination)
}

func TestStateBHandleActionNext(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := boundary.NewMockHandler(ctrl)
	ctx := context.Background()

	gomock.InOrder(
		mockHandler.EXPECT().OnExitB(ctx).Return(nil),
		mockHandler.EXPECT().OnEnterC(ctx).Return(nil),
	)
	destination, err := boundary.Handle(
		ctx, mockHandler, boundary.StateB, boundary.ActionNext, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, boundary.StateC, destination)
}

func TestStateCHandleActionNext(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := boundary.NewMockHandler(ctrl)
	ctx := context.Background()

	gomock.InOrder(
		mockHandler.EXPECT().OnExitC(ctx).Return(nil),
		mockHandler.EXPECT().OnEnterA(ctx).Return(nil),
	)
	destination, err := boundary.Handle(
		ctx, mockHandler, boundary.StateC, boundary.ActionNext, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, boundary.StateA, destination)
}
