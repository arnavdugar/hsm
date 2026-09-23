package queue_test

import (
	"errors"
	"testing"

	"github.com/arnavdugar/hsm/example/queue"
	"github.com/arnavdugar/hsm/runtime"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStateQueueEmptyHandleActionPushElement(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := queue.NewMockHandler(ctrl)
	element := queue.QueueElement{}

	mockHandler.EXPECT().HandlePushElement(element).Return(nil)
	destination, err := queue.Handle(
		mockHandler, queue.StateQueueEmpty, queue.ActionPushElement, element)

	assert.NoError(t, err)
	assert.Equal(t, queue.StateQueueHasElements, destination)
}

func TestStateQueueEmptyHandleActionPushElementHandlerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := queue.NewMockHandler(ctrl)
	element := queue.QueueElement{}
	handlerErr := errors.New("push failed")

	mockHandler.EXPECT().HandlePushElement(element).Return(handlerErr)
	destination, err := queue.Handle(
		mockHandler, queue.StateQueueEmpty, queue.ActionPushElement, element)

	assert.ErrorIs(t, err, handlerErr)
	assert.Equal(t, queue.StateQueueEmpty, destination)
}

func TestStateQueueHasElementsHandleActionPushElement(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := queue.NewMockHandler(ctrl)
	element := queue.QueueElement{}

	mockHandler.EXPECT().HandlePushElement(element).Return(nil)
	destination, err := queue.Handle(
		mockHandler, queue.StateQueueHasElements, queue.ActionPushElement, element)

	assert.NoError(t, err)
	assert.Equal(t, queue.StateQueueHasElements, destination)
}

func TestStateQueueHasElementsHandleActionConsumeElementLastElement(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := queue.NewMockHandler(ctrl)

	gomock.InOrder(
		mockHandler.EXPECT().HasSingleElement().Return(true, nil),
		mockHandler.EXPECT().HandleConsumeElement().Return(nil),
	)
	destination, err := queue.Handle(
		mockHandler, queue.StateQueueHasElements, queue.ActionConsumeElement, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, queue.StateQueueEmpty, destination)
}

func TestStateQueueHasElementsHandleActionConsumeElementLastElementHandlerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := queue.NewMockHandler(ctrl)
	handlerErr := errors.New("consume failed")

	gomock.InOrder(
		mockHandler.EXPECT().HasSingleElement().Return(true, nil),
		mockHandler.EXPECT().HandleConsumeElement().Return(handlerErr),
	)
	destination, err := queue.Handle(
		mockHandler, queue.StateQueueHasElements, queue.ActionConsumeElement, struct{}{})

	assert.ErrorIs(t, err, handlerErr)
	assert.Equal(t, queue.StateQueueHasElements, destination)
}

func TestStateQueueHasElementsHandleActionConsumeElementMultipleElements(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := queue.NewMockHandler(ctrl)

	gomock.InOrder(
		mockHandler.EXPECT().HasSingleElement().Return(false, nil),
		mockHandler.EXPECT().HandleConsumeElement().Return(nil),
	)
	destination, err := queue.Handle(
		mockHandler, queue.StateQueueHasElements, queue.ActionConsumeElement, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, queue.StateQueueHasElements, destination)
}

func TestStateQueueHasElementsHandleActionConsumeElementGuardError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := queue.NewMockHandler(ctrl)
	guardErr := errors.New("guard failed")

	mockHandler.EXPECT().HasSingleElement().Return(false, guardErr)
	destination, err := queue.Handle(
		mockHandler, queue.StateQueueHasElements, queue.ActionConsumeElement, struct{}{})

	assert.ErrorIs(t, err, guardErr)
	assert.Equal(t, queue.StateQueueHasElements, destination)
}

func TestStateQueueEmptyHandleActionClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := queue.NewMockHandler(ctrl)

	mockHandler.EXPECT().HandleClose().Return(nil)
	destination, err := queue.Handle(
		mockHandler, queue.StateQueueEmpty, queue.ActionClose, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, queue.StateClosed, destination)
}

func TestStateQueueEmptyHandleActionConsumeElement(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := queue.NewMockHandler(ctrl)

	destination, err := queue.Handle(
		mockHandler, queue.StateQueueEmpty, queue.ActionConsumeElement, struct{}{})

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, queue.StateQueueEmpty, destination)
}

func TestStateQueueHasElementsHandleActionClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := queue.NewMockHandler(ctrl)

	destination, err := queue.Handle(
		mockHandler, queue.StateQueueHasElements, queue.ActionClose, struct{}{})

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, queue.StateQueueHasElements, destination)
}

func TestStateClosedHandleActionPushElement(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := queue.NewMockHandler(ctrl)

	destination, err := queue.Handle(
		mockHandler, queue.StateClosed, queue.ActionPushElement, queue.QueueElement{})

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, queue.StateClosed, destination)
}

func TestStateQueueEmptyHandleActionPushElementInvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := queue.NewMockHandler(ctrl)

	destination, err := queue.Handle(
		mockHandler, queue.StateQueueEmpty, queue.ActionPushElement, struct{}{})

	assert.ErrorIs(t, err, runtime.ErrInvalidActionDataType)
	assert.Equal(t, queue.StateQueueEmpty, destination)
}
