import { expect, test } from "vitest";
import { ErrNoTransition } from "@arnavdugar/hsm/runtime";
import * as machine from "./machine.ts";

const unexpectedHandler: machine.Handler = {
  StartSession() {
    throw new Error("Unexpected StartSession call");
  },
  ResetRetries() {
    throw new Error("Unexpected ResetRetries call");
  },
  CanRetry() {
    throw new Error("Unexpected CanRetry call");
  },
  RecordRetry() {
    throw new Error("Unexpected RecordRetry call");
  },
};

test("Disconnected / Connect", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    StartSession() {
      calls.push("StartSession");
    },
  };
  const result = machine.handle(handler, machine.StateDisconnected, {
    type: machine.ActionConnect,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateConnecting });
  expect(calls).toEqual(["StartSession"]);
});

test("Disconnected / Connect Handler Error", () => {
  const error = new Error("handler failed");
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    StartSession() {
      calls.push("StartSession");
      return error;
    },
  };
  const result = machine.handle(handler, machine.StateDisconnected, {
    type: machine.ActionConnect,
  });
  expect(result).toStrictEqual({
    ok: false,
    state: machine.StateDisconnected,
    error,
  });
  if (!result.ok) expect(result.error).toBe(error);
  expect(calls).toEqual(["StartSession"]);
});

test("Connecting / Failure", () => {
  const result = machine.handle(unexpectedHandler, machine.StateConnecting, {
    type: machine.ActionFailure,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateBackoff });
});

test("Connecting / Transport Ready", () => {
  const result = machine.handle(unexpectedHandler, machine.StateConnecting, {
    type: machine.ActionTransportReady,
  });
  expect(result).toStrictEqual({
    ok: true,
    state: machine.StateAuthenticating,
  });
});

test("Authenticating / Failure", () => {
  const result = machine.handle(
    unexpectedHandler,
    machine.StateAuthenticating,
    { type: machine.ActionFailure },
  );
  expect(result).toStrictEqual({ ok: true, state: machine.StateBackoff });
});

test("Authenticating / Authenticated", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    ResetRetries() {
      calls.push("ResetRetries");
    },
  };
  const result = machine.handle(handler, machine.StateAuthenticating, {
    type: machine.ActionAuthenticated,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateConnected });
  expect(calls).toEqual(["ResetRetries"]);
});

test("Connected / Connection Lost", () => {
  const result = machine.handle(unexpectedHandler, machine.StateConnected, {
    type: machine.ActionConnectionLost,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateBackoff });
});

test("Connected / Disconnect", () => {
  const result = machine.handle(unexpectedHandler, machine.StateConnected, {
    type: machine.ActionDisconnect,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateDisconnected });
});

test("Backoff / Disconnect", () => {
  const result = machine.handle(unexpectedHandler, machine.StateBackoff, {
    type: machine.ActionDisconnect,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateDisconnected });
});

test("Backoff / Timeout Can Retry", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    CanRetry() {
      calls.push("CanRetry");
      return [true, undefined];
    },
    RecordRetry() {
      calls.push("RecordRetry");
    },
  };
  const result = machine.handle(handler, machine.StateBackoff, {
    type: machine.ActionTimeout,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateConnecting });
  expect(calls).toEqual(["CanRetry", "RecordRetry"]);
});

test("Backoff / Timeout Budget Exhausted", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    CanRetry() {
      calls.push("CanRetry");
      return [false, undefined];
    },
  };
  const result = machine.handle(handler, machine.StateBackoff, {
    type: machine.ActionTimeout,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateFailed });
  expect(calls).toEqual(["CanRetry"]);
});

test("Backoff / Timeout Guard Error", () => {
  const error = new Error("guard failed");
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    CanRetry() {
      calls.push("CanRetry");
      return [false, error];
    },
  };
  const result = machine.handle(handler, machine.StateBackoff, {
    type: machine.ActionTimeout,
  });
  expect(result).toStrictEqual({
    ok: false,
    state: machine.StateBackoff,
    error,
  });
  if (!result.ok) expect(result.error).toBe(error);
  expect(calls).toEqual(["CanRetry"]);
});

test("Backoff / Timeout Handler Error", () => {
  const error = new Error("handler failed");
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    CanRetry() {
      calls.push("CanRetry");
      return [true, undefined];
    },
    RecordRetry() {
      calls.push("RecordRetry");
      return error;
    },
  };
  const result = machine.handle(handler, machine.StateBackoff, {
    type: machine.ActionTimeout,
  });
  expect(result).toStrictEqual({
    ok: false,
    state: machine.StateBackoff,
    error,
  });
  if (!result.ok) expect(result.error).toBe(error);
  expect(calls).toEqual(["CanRetry", "RecordRetry"]);
});

test("Failed / Connect", () => {
  const result = machine.handle(unexpectedHandler, machine.StateFailed, {
    type: machine.ActionConnect,
  });
  expect(result).toStrictEqual({
    ok: false,
    state: machine.StateFailed,
    error: ErrNoTransition,
  });
  if (!result.ok) expect(result.error).toBe(ErrNoTransition);
});
