package garagedoor_test

import (
	"testing"

	"github.com/arnavdugar/hsm/example/garagedoor"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestStateOpeningHandleActionButton(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := garagedoor.NewMockHandler(ctrl)

	mockHandler.EXPECT().HandleButtonWhenOpening()
	destination, err := garagedoor.Handle(
		mockHandler, garagedoor.StateOpening, garagedoor.ActionButton, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, garagedoor.StateStoppedOpening, destination)
}

func TestStateStoppedClosingHandleActionButton(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := garagedoor.NewMockHandler(ctrl)

	mockHandler.EXPECT().HandleButtonWhenStoppedClosing().Return(nil)
	destination, err := garagedoor.Handle(
		mockHandler, garagedoor.StateStoppedClosing, garagedoor.ActionButton, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, garagedoor.StateOpening, destination)
}

func TestStateOpeningHandleActionOpened(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := garagedoor.NewMockHandler(ctrl)

	mockHandler.EXPECT().HandleOpened().Return(nil)
	destination, err := garagedoor.Handle(
		mockHandler, garagedoor.StateOpening, garagedoor.ActionOpened, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, garagedoor.StateStoppedOpening, destination)
}

func TestStateStoppedOpeningHandleActionButton(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := garagedoor.NewMockHandler(ctrl)

	mockHandler.EXPECT().HandleButtonWhenStoppedOpening().Return(nil)
	destination, err := garagedoor.Handle(
		mockHandler, garagedoor.StateStoppedOpening, garagedoor.ActionButton, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, garagedoor.StateClosing, destination)
}

func TestStateClosingHandleActionButton(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := garagedoor.NewMockHandler(ctrl)

	mockHandler.EXPECT().HandleButtonWhenClosing().Return(nil)
	destination, err := garagedoor.Handle(
		mockHandler, garagedoor.StateClosing, garagedoor.ActionButton, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, garagedoor.StateStoppedClosing, destination)
}

func TestStateClosingHandleActionSensor(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := garagedoor.NewMockHandler(ctrl)

	mockHandler.EXPECT().HandleSensor().Return(nil)
	destination, err := garagedoor.Handle(
		mockHandler, garagedoor.StateClosing, garagedoor.ActionSensor, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, garagedoor.StateOpening, destination)
}

func TestStateClosingHandleActionClosed(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockHandler := garagedoor.NewMockHandler(ctrl)

	mockHandler.EXPECT().HandleClosed().Return(nil)
	destination, err := garagedoor.Handle(
		mockHandler, garagedoor.StateClosing, garagedoor.ActionClosed, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, garagedoor.StateStoppedClosing, destination)
}
