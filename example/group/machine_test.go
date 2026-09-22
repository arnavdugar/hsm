package group_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/arnavdugar/hsm/example/group"
	"github.com/arnavdugar/hsm/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGroupTransitions(t *testing.T) {
	tests := []struct {
		name        string
		source      group.StateType
		action      group.ActionType
		destination group.StateType
		calls       []string
		err         error
	}{
		{
			name:   "transitive initial group",
			source: group.StateOffline, action: group.ActionStart, destination: group.StateDialing,
			calls: []string{"ExitOffline", "EnterSession", "EnterConnecting", "EnterDialing"},
		},
		{
			name:   "direct state bypasses group initials and activates overlap",
			source: group.StateOffline, action: group.ActionDirect, destination: group.StateAuthenticating,
			calls: []string{"ExitOffline", "EnterSession", "EnterConnecting", "EnterTransport", "EnterAuthenticating"},
		},
		{
			name:   "enter overlapping group",
			source: group.StateDialing, action: group.ActionNext, destination: group.StateAuthenticating,
			calls: []string{"ExitDialing", "EnterTransport", "EnterAuthenticating"},
		},
		{
			name:   "leave subgroup while retaining overlapping group",
			source: group.StateAuthenticating, action: group.ActionNext, destination: group.StateReady,
			calls: []string{"ExitAuthenticating", "ExitConnecting", "EnterReady"},
		},
		{
			name:   "remain inside shared groups",
			source: group.StateReady, action: group.ActionNext, destination: group.StateDraining,
			calls: []string{"ExitReady", "EnterDraining"},
		},
		{
			name:   "leave all groups",
			source: group.StateDraining, action: group.ActionNext, destination: group.StateOffline,
			calls: []string{"ExitDraining", "ExitTransport", "ExitSession", "EnterOffline"},
		},
		{
			name:   "external parent reenters subgroup",
			source: group.StateDialing, action: group.ActionReset, destination: group.StateDialing,
			calls: []string{"ExitDialing", "ExitConnecting", "ExitSession", "ResetSession", "EnterSession", "EnterConnecting", "EnterDialing"},
		},
		{
			name:   "external parent and natural overlap exit",
			source: group.StateAuthenticating, action: group.ActionReset, destination: group.StateDialing,
			calls: []string{"ExitAuthenticating", "ExitTransport", "ExitConnecting", "ExitSession", "ResetSession", "EnterSession", "EnterConnecting", "EnterDialing"},
		},
		{
			name:   "external parent retains unrelated overlap and state rule overrides group guard",
			source: group.StateReady, action: group.ActionRefresh, destination: group.StateReady,
			calls: []string{"ExitReady", "ExitSession", "EnterSession", "EnterReady"},
		},
		{
			name:   "group destination retains already active groups",
			source: group.StateAuthenticating, action: group.ActionRestart, destination: group.StateDialing,
			calls: []string{"ExitAuthenticating", "ExitTransport", "EnterDialing"},
		},
		{
			name:   "self transition still calls state boundaries",
			source: group.StateDialing, action: group.ActionRestart, destination: group.StateDialing,
			calls: []string{"ExitDialing", "EnterDialing"},
		},
		{
			name:   "inherited group action exits in reverse declaration order",
			source: group.StateAuthenticating, action: group.ActionClose, destination: group.StateOffline,
			calls: []string{"ExitAuthenticating", "ExitTransport", "ExitConnecting", "ExitSession", "EnterOffline"},
		},
		{
			name:   "all group guards false",
			source: group.StateDraining, action: group.ActionRefresh, destination: group.StateDraining,
			calls: []string{"CanRefresh"}, err: runtime.ErrNoTransition,
		},
		{
			name:   "group action unavailable outside group",
			source: group.StateOffline, action: group.ActionClose, destination: group.StateOffline,
			err: runtime.ErrNoTransition,
		},
		{
			name:   "unknown state",
			source: group.StateType(-1), action: group.ActionNext, destination: group.StateType(-1),
			err: runtime.ErrUnknownState,
		},
		{
			name:   "unknown action",
			source: group.StateDialing, action: group.ActionType(-1), destination: group.StateDialing,
			err: runtime.ErrUnknownAction,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := newRecordingHandler(t)
			destination, err := group.Handle(handler.ctx, handler, test.source, test.action, nil)
			if test.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, test.err)
			}
			assert.Equal(t, test.destination, destination)
			assert.Equal(t, test.calls, handler.calls)
		})
	}
}

func TestGroupActionPrecedence(t *testing.T) {
	tests := []struct {
		name        string
		data        int
		destination group.StateType
		calls       []string
	}{
		{
			name: "state rule first", data: 1, destination: group.StateAuthenticating,
			calls: []string{"CanRetryState(1)", "ExitAuthenticating", "RetryState(1)", "EnterAuthenticating"},
		},
		{
			name: "later overlapping group first", data: 2, destination: group.StateReady,
			calls: []string{"CanRetryState(2)", "CanRetryTransport(2)", "ExitAuthenticating", "ExitConnecting", "RetryTransport(2)", "EnterReady"},
		},
		{
			name: "earlier group fallback", data: 3, destination: group.StateDialing,
			calls: []string{"CanRetryState(3)", "CanRetryTransport(3)", "CanRetryConnecting(3)", "ExitAuthenticating", "ExitTransport", "RetryConnecting(3)", "EnterDialing"},
		},
		{
			name: "outer group fallback", data: 0, destination: group.StateDialing,
			calls: []string{"CanRetryState(0)", "CanRetryTransport(0)", "CanRetryConnecting(0)", "ExitAuthenticating", "ExitTransport", "RetrySession(0)", "EnterDialing"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := newRecordingHandler(t)
			destination, err := group.Handle(handler.ctx, handler, group.StateAuthenticating, group.ActionRetry, test.data)
			require.NoError(t, err)
			assert.Equal(t, test.destination, destination)
			assert.Equal(t, test.calls, handler.calls)
		})
	}
}

func TestGroupHookErrorsAbortTransition(t *testing.T) {
	calls := []string{"ExitAuthenticating", "ExitTransport", "ExitConnecting", "ExitSession", "ResetSession", "EnterSession", "EnterConnecting", "EnterDialing"}
	for index, call := range calls {
		t.Run(call, func(t *testing.T) {
			handler := newRecordingHandler(t)
			handler.failAt = call
			destination, err := group.HandleReset(handler.ctx, handler, group.StateAuthenticating)
			require.ErrorIs(t, err, handler.err)
			assert.Equal(t, group.StateAuthenticating, destination)
			assert.Equal(t, calls[:index+1], handler.calls)
		})
	}
}

func TestGroupGuardErrorsDoNotFallThrough(t *testing.T) {
	calls := []string{"CanRetryState(3)", "CanRetryTransport(3)", "CanRetryConnecting(3)"}
	for index, call := range calls {
		t.Run(call, func(t *testing.T) {
			handler := newRecordingHandler(t)
			handler.failAt = call
			destination, err := group.HandleRetry(handler.ctx, handler, group.StateAuthenticating, 3)
			require.ErrorIs(t, err, handler.err)
			assert.Equal(t, group.StateAuthenticating, destination)
			assert.Equal(t, calls[:index+1], handler.calls)
		})
	}
}

func TestGroupTransitionHandlerError(t *testing.T) {
	handler := newRecordingHandler(t)
	handler.failAt = "RetryTransport(2)"
	destination, err := group.HandleRetry(handler.ctx, handler, group.StateAuthenticating, 2)
	require.ErrorIs(t, err, handler.err)
	assert.Equal(t, group.StateAuthenticating, destination)
	assert.Equal(t, []string{"CanRetryState(2)", "CanRetryTransport(2)", "ExitAuthenticating", "ExitConnecting", "RetryTransport(2)"}, handler.calls)
}

func TestGroupActionInvalidData(t *testing.T) {
	handler := newRecordingHandler(t)
	destination, err := group.Handle(handler.ctx, handler, group.StateAuthenticating, group.ActionRetry, "invalid")
	require.ErrorIs(t, err, runtime.ErrInvalidActionDataType)
	assert.Equal(t, group.StateAuthenticating, destination)
	assert.Empty(t, handler.calls)
}

type recordingHandler struct {
	t      *testing.T
	ctx    context.Context
	calls  []string
	failAt string
	err    error
}

func newRecordingHandler(t *testing.T) *recordingHandler {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return &recordingHandler{t: t, ctx: ctx, err: errors.New("handler failed")}
}

func (handler *recordingHandler) record(ctx context.Context, name string) error {
	handler.t.Helper()
	assert.Same(handler.t, handler.ctx, ctx)
	handler.calls = append(handler.calls, name)
	if name == handler.failAt {
		return handler.err
	}
	return nil
}

func (handler *recordingHandler) EnterOffline(ctx context.Context) error {
	return handler.record(ctx, "EnterOffline")
}

func (handler *recordingHandler) ExitOffline(ctx context.Context) error {
	return handler.record(ctx, "ExitOffline")
}

func (handler *recordingHandler) EnterDialing(ctx context.Context) error {
	return handler.record(ctx, "EnterDialing")
}

func (handler *recordingHandler) ExitDialing(ctx context.Context) error {
	return handler.record(ctx, "ExitDialing")
}

func (handler *recordingHandler) EnterAuthenticating(ctx context.Context) error {
	return handler.record(ctx, "EnterAuthenticating")
}

func (handler *recordingHandler) ExitAuthenticating(ctx context.Context) error {
	return handler.record(ctx, "ExitAuthenticating")
}

func (handler *recordingHandler) EnterReady(ctx context.Context) error {
	return handler.record(ctx, "EnterReady")
}

func (handler *recordingHandler) ExitReady(ctx context.Context) error {
	return handler.record(ctx, "ExitReady")
}

func (handler *recordingHandler) EnterDraining(ctx context.Context) error {
	return handler.record(ctx, "EnterDraining")
}

func (handler *recordingHandler) ExitDraining(ctx context.Context) error {
	return handler.record(ctx, "ExitDraining")
}

func (handler *recordingHandler) EnterSession(ctx context.Context) error {
	return handler.record(ctx, "EnterSession")
}

func (handler *recordingHandler) ExitSession(ctx context.Context) error {
	return handler.record(ctx, "ExitSession")
}

func (handler *recordingHandler) EnterConnecting(ctx context.Context) error {
	return handler.record(ctx, "EnterConnecting")
}

func (handler *recordingHandler) ExitConnecting(ctx context.Context) error {
	return handler.record(ctx, "ExitConnecting")
}

func (handler *recordingHandler) EnterTransport(ctx context.Context) error {
	return handler.record(ctx, "EnterTransport")
}

func (handler *recordingHandler) ExitTransport(ctx context.Context) error {
	return handler.record(ctx, "ExitTransport")
}

func (handler *recordingHandler) ResetSession(ctx context.Context) error {
	return handler.record(ctx, "ResetSession")
}

func (handler *recordingHandler) CanRefresh(ctx context.Context) (bool, error) {
	return false, handler.record(ctx, "CanRefresh")
}

func (handler *recordingHandler) CanRetryState(ctx context.Context, data int) (bool, error) {
	return data == 1, handler.record(ctx, fmt.Sprintf("CanRetryState(%d)", data))
}

func (handler *recordingHandler) CanRetryTransport(ctx context.Context, data int) (bool, error) {
	return data == 2, handler.record(ctx, fmt.Sprintf("CanRetryTransport(%d)", data))
}

func (handler *recordingHandler) CanRetryConnecting(ctx context.Context, data int) (bool, error) {
	return data == 3, handler.record(ctx, fmt.Sprintf("CanRetryConnecting(%d)", data))
}

func (handler *recordingHandler) RetryState(ctx context.Context, data int) error {
	return handler.record(ctx, fmt.Sprintf("RetryState(%d)", data))
}

func (handler *recordingHandler) RetryTransport(ctx context.Context, data int) error {
	return handler.record(ctx, fmt.Sprintf("RetryTransport(%d)", data))
}

func (handler *recordingHandler) RetryConnecting(ctx context.Context, data int) error {
	return handler.record(ctx, fmt.Sprintf("RetryConnecting(%d)", data))
}

func (handler *recordingHandler) RetrySession(ctx context.Context, data int) error {
	return handler.record(ctx, fmt.Sprintf("RetrySession(%d)", data))
}
