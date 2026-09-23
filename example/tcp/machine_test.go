package tcp_test

import (
	"testing"

	"github.com/arnavdugar/hsm/example/tcp"
	"github.com/arnavdugar/hsm/runtime"
	"github.com/stretchr/testify/assert"
)

func TestStateClosedHandleActionActiveOpen(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateClosed, tcp.ActionActiveOpen, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateSynSent, destination)
}

func TestStateSynSentHandleActionReceiveSynAck(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateSynSent, tcp.ActionReceiveSynAck, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateEstablished, destination)
}

func TestStateEstablishedHandleActionClose(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateEstablished, tcp.ActionClose, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateFinWait1, destination)
}

func TestStateFinWait1HandleActionReceiveAck(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateFinWait1, tcp.ActionReceiveAck, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateFinWait2, destination)
}

func TestStateFinWait2HandleActionReceiveFin(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateFinWait2, tcp.ActionReceiveFin, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateTimeWait, destination)
}

func TestStateTimeWaitHandleActionTimeout(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateTimeWait, tcp.ActionTimeout, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateClosed, destination)
}

func TestStateClosedHandleActionPassiveOpen(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateClosed, tcp.ActionPassiveOpen, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateListen, destination)
}

func TestStateListenHandleActionReceiveSyn(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateListen, tcp.ActionReceiveSyn, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateSynReceived, destination)
}

func TestStateSynReceivedHandleActionReceiveAck(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateSynReceived, tcp.ActionReceiveAck, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateEstablished, destination)
}

func TestStateEstablishedHandleActionReceiveFin(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateEstablished, tcp.ActionReceiveFin, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateCloseWait, destination)
}

func TestStateCloseWaitHandleActionClose(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateCloseWait, tcp.ActionClose, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateLastAck, destination)
}

func TestStateLastAckHandleActionReceiveAck(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateLastAck, tcp.ActionReceiveAck, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateClosed, destination)
}

func TestStateSynSentHandleActionReceiveSyn(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateSynSent, tcp.ActionReceiveSyn, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateSynReceived, destination)
}

func TestStateFinWait1HandleActionReceiveFin(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateFinWait1, tcp.ActionReceiveFin, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateClosing, destination)
}

func TestStateClosingHandleActionReceiveAck(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateClosing, tcp.ActionReceiveAck, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateTimeWait, destination)
}

func TestStateFinWait1HandleActionReceiveFinAck(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateFinWait1, tcp.ActionReceiveFinAck, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateTimeWait, destination)
}

func TestStateListenHandleActionClose(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateListen, tcp.ActionClose, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateClosed, destination)
}

func TestStateListenHandleActionSend(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateListen, tcp.ActionSend, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateSynSent, destination)
}

func TestStateSynSentHandleActionClose(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateSynSent, tcp.ActionClose, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateClosed, destination)
}

func TestStateSynReceivedHandleActionClose(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateSynReceived, tcp.ActionClose, struct{}{})

	assert.NoError(t, err)
	assert.Equal(t, tcp.StateFinWait1, destination)
}

func TestStateClosedHandleActionTimeout(t *testing.T) {
	destination, err := tcp.Handle(
		nil, tcp.StateClosed, tcp.ActionTimeout, struct{}{})

	assert.ErrorIs(t, err, runtime.ErrNoTransition)
	assert.Equal(t, tcp.StateClosed, destination)
}
