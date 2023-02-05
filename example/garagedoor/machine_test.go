package garagedoor_test

import (
	"testing"

	"github.com/arnavdugar/hsm/example/garagedoor"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
