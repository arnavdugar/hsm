import { expect, test } from "vitest";
import {
  ErrNoTransition,
  ErrUnknownAction,
  ErrUnknownState,
} from "@arnavdugar/hsm/runtime";
import * as machine from "./machine.ts";

function recordingHandler(failAt?: string) {
  const context = { requestId: "test" };
  const calls: string[] = [];
  const error = new Error("handler failed");

  function record(ctx: typeof context, name: string) {
    expect(ctx).toBe(context);
    calls.push(name);
    return name === failAt ? error : undefined;
  }

  const handler: machine.Handler = {
    EnterOffline: (ctx) => record(ctx, "EnterOffline"),
    ExitOffline: (ctx) => record(ctx, "ExitOffline"),
    EnterDialing: (ctx) => record(ctx, "EnterDialing"),
    ExitDialing: (ctx) => record(ctx, "ExitDialing"),
    EnterAuthenticating: (ctx) => record(ctx, "EnterAuthenticating"),
    ExitAuthenticating: (ctx) => record(ctx, "ExitAuthenticating"),
    EnterReady: (ctx) => record(ctx, "EnterReady"),
    ExitReady: (ctx) => record(ctx, "ExitReady"),
    EnterDraining: (ctx) => record(ctx, "EnterDraining"),
    ExitDraining: (ctx) => record(ctx, "ExitDraining"),
    EnterSession: (ctx) => record(ctx, "EnterSession"),
    ExitSession: (ctx) => record(ctx, "ExitSession"),
    EnterConnecting: (ctx) => record(ctx, "EnterConnecting"),
    ExitConnecting: (ctx) => record(ctx, "ExitConnecting"),
    EnterTransport: (ctx) => record(ctx, "EnterTransport"),
    ExitTransport: (ctx) => record(ctx, "ExitTransport"),
    ResetSession: (ctx) => record(ctx, "ResetSession"),
    CanRefresh: (ctx) => [false, record(ctx, "CanRefresh")],
    CanRetryState: (ctx, data) => [
      data === 1,
      record(ctx, `CanRetryState(${data})`),
    ],
    CanRetryTransport: (ctx, data) => [
      data === 2,
      record(ctx, `CanRetryTransport(${data})`),
    ],
    CanRetryConnecting: (ctx, data) => [
      data === 3,
      record(ctx, `CanRetryConnecting(${data})`),
    ],
    RetryState: (ctx, data) => record(ctx, `RetryState(${data})`),
    RetryTransport: (ctx, data) => record(ctx, `RetryTransport(${data})`),
    RetryConnecting: (ctx, data) => record(ctx, `RetryConnecting(${data})`),
    RetrySession: (ctx, data) => record(ctx, `RetrySession(${data})`),
  };
  return { context, calls, error, handler };
}

const transitions = [
  [
    "transitive initial group",
    machine.StateOffline,
    { type: machine.ActionStart },
    machine.StateDialing,
    ["ExitOffline", "EnterSession", "EnterConnecting", "EnterDialing"],
  ],
  [
    "direct state bypasses group initials and activates overlap",
    machine.StateOffline,
    { type: machine.ActionDirect },
    machine.StateAuthenticating,
    [
      "ExitOffline",
      "EnterSession",
      "EnterConnecting",
      "EnterTransport",
      "EnterAuthenticating",
    ],
  ],
  [
    "enter overlapping group",
    machine.StateDialing,
    { type: machine.ActionNext },
    machine.StateAuthenticating,
    ["ExitDialing", "EnterTransport", "EnterAuthenticating"],
  ],
  [
    "leave subgroup while retaining overlapping group",
    machine.StateAuthenticating,
    { type: machine.ActionNext },
    machine.StateReady,
    ["ExitAuthenticating", "ExitConnecting", "EnterReady"],
  ],
  [
    "remain inside shared groups",
    machine.StateReady,
    { type: machine.ActionNext },
    machine.StateDraining,
    ["ExitReady", "EnterDraining"],
  ],
  [
    "leave all groups",
    machine.StateDraining,
    { type: machine.ActionNext },
    machine.StateOffline,
    ["ExitDraining", "ExitTransport", "ExitSession", "EnterOffline"],
  ],
  [
    "external parent reenters subgroup",
    machine.StateDialing,
    { type: machine.ActionReset },
    machine.StateDialing,
    [
      "ExitDialing",
      "ExitConnecting",
      "ExitSession",
      "ResetSession",
      "EnterSession",
      "EnterConnecting",
      "EnterDialing",
    ],
  ],
  [
    "external parent and natural overlap exit",
    machine.StateAuthenticating,
    { type: machine.ActionReset },
    machine.StateDialing,
    [
      "ExitAuthenticating",
      "ExitTransport",
      "ExitConnecting",
      "ExitSession",
      "ResetSession",
      "EnterSession",
      "EnterConnecting",
      "EnterDialing",
    ],
  ],
  [
    "external parent retains unrelated overlap and state rule overrides group guard",
    machine.StateReady,
    { type: machine.ActionRefresh },
    machine.StateReady,
    ["ExitReady", "ExitSession", "EnterSession", "EnterReady"],
  ],
  [
    "group destination retains already active groups",
    machine.StateAuthenticating,
    { type: machine.ActionRestart },
    machine.StateDialing,
    ["ExitAuthenticating", "ExitTransport", "EnterDialing"],
  ],
  [
    "self transition still calls state boundaries",
    machine.StateDialing,
    { type: machine.ActionRestart },
    machine.StateDialing,
    ["ExitDialing", "EnterDialing"],
  ],
  [
    "inherited group action exits in reverse declaration order",
    machine.StateAuthenticating,
    { type: machine.ActionClose },
    machine.StateOffline,
    [
      "ExitAuthenticating",
      "ExitTransport",
      "ExitConnecting",
      "ExitSession",
      "EnterOffline",
    ],
  ],
  [
    "all group guards false",
    machine.StateDraining,
    { type: machine.ActionRefresh },
    machine.StateDraining,
    ["CanRefresh"],
    ErrNoTransition,
  ],
  [
    "group action unavailable outside group",
    machine.StateOffline,
    { type: machine.ActionClose },
    machine.StateOffline,
    [],
    ErrNoTransition,
  ],
  [
    "unknown state",
    -1 as machine.StateType,
    { type: machine.ActionNext },
    -1,
    [],
    ErrUnknownState,
  ],
  [
    "unknown action",
    machine.StateDialing,
    { type: -1 } as unknown as machine.Action,
    machine.StateDialing,
    [],
    ErrUnknownAction,
  ],
] as const;

test.each(transitions)("%s", (_name, source, action, state, calls, error?) => {
  const actual = recordingHandler();
  const result = machine.handle(actual.context, actual.handler, source, action);
  if (error) {
    expect(result).toStrictEqual({ ok: false, state, error });
    if (!result.ok) expect(result.error).toBe(error);
  } else {
    expect(result).toStrictEqual({ ok: true, state });
  }
  expect(actual.calls).toEqual(calls);
});

const retries = [
  [
    "state rule first",
    1,
    machine.StateAuthenticating,
    [
      "CanRetryState(1)",
      "ExitAuthenticating",
      "RetryState(1)",
      "EnterAuthenticating",
    ],
  ],
  [
    "later overlapping group first",
    2,
    machine.StateReady,
    [
      "CanRetryState(2)",
      "CanRetryTransport(2)",
      "ExitAuthenticating",
      "ExitConnecting",
      "RetryTransport(2)",
      "EnterReady",
    ],
  ],
  [
    "earlier group fallback",
    3,
    machine.StateDialing,
    [
      "CanRetryState(3)",
      "CanRetryTransport(3)",
      "CanRetryConnecting(3)",
      "ExitAuthenticating",
      "ExitTransport",
      "RetryConnecting(3)",
      "EnterDialing",
    ],
  ],
  [
    "outer group fallback",
    0,
    machine.StateDialing,
    [
      "CanRetryState(0)",
      "CanRetryTransport(0)",
      "CanRetryConnecting(0)",
      "ExitAuthenticating",
      "ExitTransport",
      "RetrySession(0)",
      "EnterDialing",
    ],
  ],
] as const;

test.each(retries)("%s", (_name, data, state, calls) => {
  const actual = recordingHandler();
  const result = machine.handle(
    actual.context,
    actual.handler,
    machine.StateAuthenticating,
    {
      type: machine.ActionRetry,
      data,
    },
  );
  expect(result).toStrictEqual({ ok: true, state });
  expect(actual.calls).toEqual(calls);
});

const resetCalls = [
  "ExitAuthenticating",
  "ExitTransport",
  "ExitConnecting",
  "ExitSession",
  "ResetSession",
  "EnterSession",
  "EnterConnecting",
  "EnterDialing",
];

test.each(resetCalls.map((call, index) => [call, index] as const))(
  "reset stops after %s fails",
  (call, index) => {
    const { context, handler, error, calls } = recordingHandler(call);
    const state = machine.StateAuthenticating;
    const result = machine.handleReset(context, handler, state);
    expect(result).toStrictEqual({ ok: false, state, error });
    if (!result.ok) expect(result.error).toBe(error);
    expect(calls).toEqual(resetCalls.slice(0, index + 1));
  },
);

const retryGuards = [
  "CanRetryState(3)",
  "CanRetryTransport(3)",
  "CanRetryConnecting(3)",
];

test.each(retryGuards.map((guard, index) => [guard, index] as const))(
  "retry stops after %s fails",
  (guard, index) => {
    const { context, handler, error, calls } = recordingHandler(guard);
    const state = machine.StateAuthenticating;
    const result = machine.handleRetry(context, handler, state, 3);
    expect(result).toStrictEqual({ ok: false, state, error });
    if (!result.ok) expect(result.error).toBe(error);
    expect(calls).toEqual(retryGuards.slice(0, index + 1));
  },
);

test("retry stops after the transition handler fails", () => {
  const { context, handler, error, calls } =
    recordingHandler("RetryTransport(2)");
  const state = machine.StateAuthenticating;
  const result = machine.handleRetry(context, handler, state, 2);
  expect(result).toStrictEqual({ ok: false, state, error });
  if (!result.ok) expect(result.error).toBe(error);
  expect(calls).toEqual([
    "CanRetryState(2)",
    "CanRetryTransport(2)",
    "ExitAuthenticating",
    "ExitConnecting",
    "RetryTransport(2)",
  ]);
});

test("a missing action returns the shared unknown action error", () => {
  const { context, handler, calls } = recordingHandler();
  const state = machine.StateAuthenticating;
  // @ts-expect-error Untyped callers can still supply a null action.
  const result = machine.handle(context, handler, state, null);
  expect(result).toStrictEqual({ ok: false, state, error: ErrUnknownAction });
  if (!result.ok) expect(result.error).toBe(ErrUnknownAction);
  expect(calls).toEqual([]);
});
