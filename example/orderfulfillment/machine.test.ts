import { expect, test } from "vitest";
import { ErrNoTransition } from "@arnavdugar/hsm/runtime";
import * as machine from "./machine.ts";

const transitions = [
  ["Paid / Ship", machine.StatePaid, machine.ActionShip, machine.StatePacking],
  [
    "Packing / Dispatch",
    machine.StatePacking,
    machine.ActionDispatch,
    machine.StateShipped,
  ],
  [
    "Shipped / Deliver",
    machine.StateShipped,
    machine.ActionDeliver,
    machine.StateCompleted,
  ],
  [
    "Paid / Pick Up",
    machine.StatePaid,
    machine.ActionPickUp,
    machine.StateReadyForPickup,
  ],
  [
    "Ready For Pickup / Collect",
    machine.StateReadyForPickup,
    machine.ActionCollect,
    machine.StateCompleted,
  ],
  [
    "Paid / Cancel",
    machine.StatePaid,
    machine.ActionCancel,
    machine.StateCancelled,
  ],
  [
    "Packing / Cancel",
    machine.StatePacking,
    machine.ActionCancel,
    machine.StateCancelled,
  ],
  [
    "Ready For Pickup / Cancel",
    machine.StateReadyForPickup,
    machine.ActionCancel,
    machine.StateCancelled,
  ],
] as const;

test.each(transitions)("%s", (_name, source, type, state) => {
  expect(machine.handle({}, source, { type })).toStrictEqual({
    ok: true,
    state,
  });
});

const rejectedTransitions = [
  ["Paid / Deliver", machine.StatePaid, machine.ActionDeliver],
  ["Completed / Cancel", machine.StateCompleted, machine.ActionCancel],
  ["Cancelled / Ship", machine.StateCancelled, machine.ActionShip],
] as const;

test.each(rejectedTransitions)("%s has no transition", (_name, state, type) => {
  const result = machine.handle({}, state, { type });
  expect(result).toStrictEqual({ ok: false, state, error: ErrNoTransition });
  if (!result.ok) expect(result.error).toBe(ErrNoTransition);
});
