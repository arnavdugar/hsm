import { expect, test } from "vitest";
import * as machine from "./machine.ts";

const unexpectedHandler: machine.Handler = {
  HandleButtonWhenStoppedClosing() {
    throw new Error("Unexpected HandleButtonWhenStoppedClosing call");
  },
  HandleButtonWhenOpening() {
    throw new Error("Unexpected HandleButtonWhenOpening call");
  },
  HandleButtonWhenStoppedOpening() {
    throw new Error("Unexpected HandleButtonWhenStoppedOpening call");
  },
  HandleButtonWhenClosing() {
    throw new Error("Unexpected HandleButtonWhenClosing call");
  },
  HandleOpened() {
    throw new Error("Unexpected HandleOpened call");
  },
  HandleClosed() {
    throw new Error("Unexpected HandleClosed call");
  },
  HandleSensor() {
    throw new Error("Unexpected HandleSensor call");
  },
};

test("Opening / Button", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HandleButtonWhenOpening() {
      calls.push("HandleButtonWhenOpening");
    },
  };
  const result = machine.handle(handler, machine.StateOpening, {
    type: machine.ActionButton,
  });
  expect(result).toStrictEqual({
    ok: true,
    state: machine.StateStoppedOpening,
  });
  expect(calls).toEqual(["HandleButtonWhenOpening"]);
});

test("Stopped Closing / Button", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HandleButtonWhenStoppedClosing() {
      calls.push("HandleButtonWhenStoppedClosing");
    },
  };
  const result = machine.handle(handler, machine.StateStoppedClosing, {
    type: machine.ActionButton,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateOpening });
  expect(calls).toEqual(["HandleButtonWhenStoppedClosing"]);
});

test("Opening / Opened", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HandleOpened() {
      calls.push("HandleOpened");
    },
  };
  const result = machine.handle(handler, machine.StateOpening, {
    type: machine.ActionOpened,
  });
  expect(result).toStrictEqual({
    ok: true,
    state: machine.StateStoppedOpening,
  });
  expect(calls).toEqual(["HandleOpened"]);
});

test("Stopped Opening / Button", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HandleButtonWhenStoppedOpening() {
      calls.push("HandleButtonWhenStoppedOpening");
    },
  };
  const result = machine.handle(handler, machine.StateStoppedOpening, {
    type: machine.ActionButton,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateClosing });
  expect(calls).toEqual(["HandleButtonWhenStoppedOpening"]);
});

test("Closing / Button", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HandleButtonWhenClosing() {
      calls.push("HandleButtonWhenClosing");
    },
  };
  const result = machine.handle(handler, machine.StateClosing, {
    type: machine.ActionButton,
  });
  expect(result).toStrictEqual({
    ok: true,
    state: machine.StateStoppedClosing,
  });
  expect(calls).toEqual(["HandleButtonWhenClosing"]);
});

test("Closing / Sensor", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HandleSensor() {
      calls.push("HandleSensor");
    },
  };
  const result = machine.handle(handler, machine.StateClosing, {
    type: machine.ActionSensor,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateOpening });
  expect(calls).toEqual(["HandleSensor"]);
});

test("Closing / Closed", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HandleClosed() {
      calls.push("HandleClosed");
    },
  };
  const result = machine.handle(handler, machine.StateClosing, {
    type: machine.ActionClosed,
  });
  expect(result).toStrictEqual({
    ok: true,
    state: machine.StateStoppedClosing,
  });
  expect(calls).toEqual(["HandleClosed"]);
});
