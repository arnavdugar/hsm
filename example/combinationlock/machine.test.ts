import { expect, test } from "vitest";
import { ErrNoTransition } from "@arnavdugar/hsm/runtime";
import * as machine from "./machine.ts";

const unexpectedHandler: machine.Handler = {
  IsFirstDigit() {
    throw new Error("Unexpected IsFirstDigit call");
  },
  IsSecondDigit() {
    throw new Error("Unexpected IsSecondDigit call");
  },
  IsThirdDigit() {
    throw new Error("Unexpected IsThirdDigit call");
  },
};

test("Locked / Digit Correct Digit", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    IsFirstDigit(value) {
      calls.push("IsFirstDigit");
      expect(value).toBe(4);
      return [true, undefined];
    },
  };
  const result = machine.handle(handler, machine.StateLocked, {
    type: machine.ActionDigit,
    data: 4,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateFirstDigit });
  expect(calls).toEqual(["IsFirstDigit"]);
});

test("Locked / Digit Incorrect Digit", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    IsFirstDigit(value) {
      calls.push("IsFirstDigit");
      expect(value).toBe(0);
      return [false, undefined];
    },
  };
  const result = machine.handle(handler, machine.StateLocked, {
    type: machine.ActionDigit,
    data: 0,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateLocked });
  expect(calls).toEqual(["IsFirstDigit"]);
});

test("Locked / Digit Guard Error", () => {
  const error = new Error("guard failed");
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    IsFirstDigit(value) {
      calls.push("IsFirstDigit");
      expect(value).toBe(4);
      return [false, error];
    },
  };
  const result = machine.handle(handler, machine.StateLocked, {
    type: machine.ActionDigit,
    data: 4,
  });
  expect(result).toStrictEqual({
    ok: false,
    state: machine.StateLocked,
    error,
  });
  if (!result.ok) expect(result.error).toBe(error);
  expect(calls).toEqual(["IsFirstDigit"]);
});

test("First Digit / Digit Correct Digit", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    IsSecondDigit(value) {
      calls.push("IsSecondDigit");
      expect(value).toBe(2);
      return [true, undefined];
    },
  };
  const result = machine.handle(handler, machine.StateFirstDigit, {
    type: machine.ActionDigit,
    data: 2,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateSecondDigit });
  expect(calls).toEqual(["IsSecondDigit"]);
});

test("First Digit / Digit Incorrect Digit", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    IsSecondDigit(value) {
      calls.push("IsSecondDigit");
      expect(value).toBe(0);
      return [false, undefined];
    },
  };
  const result = machine.handle(handler, machine.StateFirstDigit, {
    type: machine.ActionDigit,
    data: 0,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateLocked });
  expect(calls).toEqual(["IsSecondDigit"]);
});

test("Second Digit / Digit Correct Digit", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    IsThirdDigit(value) {
      calls.push("IsThirdDigit");
      expect(value).toBe(7);
      return [true, undefined];
    },
  };
  const result = machine.handle(handler, machine.StateSecondDigit, {
    type: machine.ActionDigit,
    data: 7,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateUnlocked });
  expect(calls).toEqual(["IsThirdDigit"]);
});

test("Second Digit / Digit Incorrect Digit", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    IsThirdDigit(value) {
      calls.push("IsThirdDigit");
      expect(value).toBe(0);
      return [false, undefined];
    },
  };
  const result = machine.handle(handler, machine.StateSecondDigit, {
    type: machine.ActionDigit,
    data: 0,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateLocked });
  expect(calls).toEqual(["IsThirdDigit"]);
});

test("Unlocked / Lock", () => {
  const result = machine.handle(unexpectedHandler, machine.StateUnlocked, {
    type: machine.ActionLock,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateLocked });
});

test("Unlocked / Digit", () => {
  const result = machine.handle(unexpectedHandler, machine.StateUnlocked, {
    type: machine.ActionDigit,
    data: 4,
  });
  expect(result).toStrictEqual({
    ok: false,
    state: machine.StateUnlocked,
    error: ErrNoTransition,
  });
  if (!result.ok) expect(result.error).toBe(ErrNoTransition);
});
