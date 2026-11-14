import { expect, test } from "vitest";
import * as machine from "./machine.ts";

const unexpectedHandler: machine.Handler = {
  OnEnterA() {
    throw new Error("Unexpected OnEnterA call");
  },
  OnExitA() {
    throw new Error("Unexpected OnExitA call");
  },
  OnEnterB() {
    throw new Error("Unexpected OnEnterB call");
  },
  OnExitB() {
    throw new Error("Unexpected OnExitB call");
  },
  OnEnterC() {
    throw new Error("Unexpected OnEnterC call");
  },
  OnExitC() {
    throw new Error("Unexpected OnExitC call");
  },
};

test("A / Next", () => {
  const context = {};
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    OnExitA(value) {
      calls.push("OnExitA");
      expect(value).toBe(context);
    },
    OnEnterB(value) {
      calls.push("OnEnterB");
      expect(value).toBe(context);
    },
  };
  const result = machine.handle(context, handler, machine.StateA, {
    type: machine.ActionNext,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateB });
  expect(calls).toEqual(["OnExitA", "OnEnterB"]);
});

test("A / Next Exit Error", () => {
  const context = {};
  const error = new Error("exit failed");
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    OnExitA(value) {
      calls.push("OnExitA");
      expect(value).toBe(context);
      return error;
    },
  };
  const result = machine.handle(context, handler, machine.StateA, {
    type: machine.ActionNext,
  });
  expect(result).toStrictEqual({ ok: false, state: machine.StateA, error });
  if (!result.ok) expect(result.error).toBe(error);
  expect(calls).toEqual(["OnExitA"]);
});

test("A / Next Enter Error", () => {
  const context = {};
  const error = new Error("enter failed");
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    OnExitA(value) {
      calls.push("OnExitA");
      expect(value).toBe(context);
    },
    OnEnterB(value) {
      calls.push("OnEnterB");
      expect(value).toBe(context);
      return error;
    },
  };
  const result = machine.handle(context, handler, machine.StateA, {
    type: machine.ActionNext,
  });
  expect(result).toStrictEqual({ ok: false, state: machine.StateA, error });
  if (!result.ok) expect(result.error).toBe(error);
  expect(calls).toEqual(["OnExitA", "OnEnterB"]);
});

test("B / Next", () => {
  const context = {};
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    OnExitB(value) {
      calls.push("OnExitB");
      expect(value).toBe(context);
    },
    OnEnterC(value) {
      calls.push("OnEnterC");
      expect(value).toBe(context);
    },
  };
  const result = machine.handle(context, handler, machine.StateB, {
    type: machine.ActionNext,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateC });
  expect(calls).toEqual(["OnExitB", "OnEnterC"]);
});

test("C / Next", () => {
  const context = {};
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    OnExitC(value) {
      calls.push("OnExitC");
      expect(value).toBe(context);
    },
    OnEnterA(value) {
      calls.push("OnEnterA");
      expect(value).toBe(context);
    },
  };
  const result = machine.handle(context, handler, machine.StateC, {
    type: machine.ActionNext,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateA });
  expect(calls).toEqual(["OnExitC", "OnEnterA"]);
});
