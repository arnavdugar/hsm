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
