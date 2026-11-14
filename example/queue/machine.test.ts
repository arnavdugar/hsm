import { expect, test } from "vitest";
import { ErrNoTransition } from "@arnavdugar/hsm/runtime";
import * as machine from "./machine.ts";

const unexpectedHandler: machine.Handler = {
  HandlePushElement() {
    throw new Error("Unexpected HandlePushElement call");
  },
  HasSingleElement() {
    throw new Error("Unexpected HasSingleElement call");
  },
  HandleConsumeElement() {
    throw new Error("Unexpected HandleConsumeElement call");
  },
  HandleClose() {
    throw new Error("Unexpected HandleClose call");
  },
};

test("Queue Empty / Push Element", () => {
  const element = {};
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HandlePushElement(value) {
      calls.push("HandlePushElement");
      expect(value).toBe(element);
    },
  };
  const result = machine.handle(handler, machine.StateQueueEmpty, {
    type: machine.ActionPushElement,
    data: element,
  });
  expect(result).toStrictEqual({
    ok: true,
    state: machine.StateQueueHasElements,
  });
  expect(calls).toEqual(["HandlePushElement"]);
});

test("Queue Empty / Push Element Handler Error", () => {
  const element = {};
  const error = new Error("push failed");
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HandlePushElement(value) {
      calls.push("HandlePushElement");
      expect(value).toBe(element);
      return error;
    },
  };
  const result = machine.handle(handler, machine.StateQueueEmpty, {
    type: machine.ActionPushElement,
    data: element,
  });
  expect(result).toStrictEqual({
    ok: false,
    state: machine.StateQueueEmpty,
    error,
  });
  if (!result.ok) expect(result.error).toBe(error);
  expect(calls).toEqual(["HandlePushElement"]);
});

test("Queue Has Elements / Push Element", () => {
  const element = {};
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HandlePushElement(value) {
      calls.push("HandlePushElement");
      expect(value).toBe(element);
    },
  };
  const result = machine.handle(handler, machine.StateQueueHasElements, {
    type: machine.ActionPushElement,
    data: element,
  });
  expect(result).toStrictEqual({
    ok: true,
    state: machine.StateQueueHasElements,
  });
  expect(calls).toEqual(["HandlePushElement"]);
});

test("Queue Has Elements / Consume Element Last Element", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HasSingleElement() {
      calls.push("HasSingleElement");
      return [true, undefined];
    },
    HandleConsumeElement() {
      calls.push("HandleConsumeElement");
    },
  };
  const result = machine.handle(handler, machine.StateQueueHasElements, {
    type: machine.ActionConsumeElement,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateQueueEmpty });
  expect(calls).toEqual(["HasSingleElement", "HandleConsumeElement"]);
});

test("Queue Has Elements / Consume Element Last Element Handler Error", () => {
  const error = new Error("consume failed");
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HasSingleElement() {
      calls.push("HasSingleElement");
      return [true, undefined];
    },
    HandleConsumeElement() {
      calls.push("HandleConsumeElement");
      return error;
    },
  };
  const result = machine.handle(handler, machine.StateQueueHasElements, {
    type: machine.ActionConsumeElement,
  });
  expect(result).toStrictEqual({
    ok: false,
    state: machine.StateQueueHasElements,
    error,
  });
  if (!result.ok) expect(result.error).toBe(error);
  expect(calls).toEqual(["HasSingleElement", "HandleConsumeElement"]);
});

test("Queue Has Elements / Consume Element Multiple Elements", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HasSingleElement() {
      calls.push("HasSingleElement");
      return [false, undefined];
    },
    HandleConsumeElement() {
      calls.push("HandleConsumeElement");
    },
  };
  const result = machine.handle(handler, machine.StateQueueHasElements, {
    type: machine.ActionConsumeElement,
  });
  expect(result).toStrictEqual({
    ok: true,
    state: machine.StateQueueHasElements,
  });
  expect(calls).toEqual(["HasSingleElement", "HandleConsumeElement"]);
});

test("Queue Has Elements / Consume Element Guard Error", () => {
  const error = new Error("guard failed");
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HasSingleElement() {
      calls.push("HasSingleElement");
      return [false, error];
    },
  };
  const result = machine.handle(handler, machine.StateQueueHasElements, {
    type: machine.ActionConsumeElement,
  });
  expect(result).toStrictEqual({
    ok: false,
    state: machine.StateQueueHasElements,
    error,
  });
  if (!result.ok) expect(result.error).toBe(error);
  expect(calls).toEqual(["HasSingleElement"]);
});

test("Queue Empty / Close", () => {
  const calls: string[] = [];
  const handler: machine.Handler = {
    ...unexpectedHandler,
    HandleClose() {
      calls.push("HandleClose");
    },
  };
  const result = machine.handle(handler, machine.StateQueueEmpty, {
    type: machine.ActionClose,
  });
  expect(result).toStrictEqual({ ok: true, state: machine.StateClosed });
  expect(calls).toEqual(["HandleClose"]);
});

test("Queue Empty / Consume Element", () => {
  const result = machine.handle(unexpectedHandler, machine.StateQueueEmpty, {
    type: machine.ActionConsumeElement,
  });
  expect(result).toStrictEqual({
    ok: false,
    state: machine.StateQueueEmpty,
    error: ErrNoTransition,
  });
  if (!result.ok) expect(result.error).toBe(ErrNoTransition);
});

test("Queue Has Elements / Close", () => {
  const result = machine.handle(
    unexpectedHandler,
    machine.StateQueueHasElements,
    { type: machine.ActionClose },
  );
  expect(result).toStrictEqual({
    ok: false,
    state: machine.StateQueueHasElements,
    error: ErrNoTransition,
  });
  if (!result.ok) expect(result.error).toBe(ErrNoTransition);
});

test("Closed / Push Element", () => {
  const result = machine.handle(unexpectedHandler, machine.StateClosed, {
    type: machine.ActionPushElement,
    data: {},
  });
  expect(result).toStrictEqual({
    ok: false,
    state: machine.StateClosed,
    error: ErrNoTransition,
  });
  if (!result.ok) expect(result.error).toBe(ErrNoTransition);
});
