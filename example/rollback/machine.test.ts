import { expect, test } from "vitest";
import { ErrNoTransition } from "@arnavdugar/hsm/runtime";
import * as machine from "./machine.ts";

const transitions = [
  [
    "Reserve / Reserved",
    machine.StateReserve,
    machine.ActionReserved,
    machine.StateCharge,
  ],
  [
    "Charge / Charged",
    machine.StateCharge,
    machine.ActionCharged,
    machine.StateFulfill,
  ],
  [
    "Fulfill / Fulfilled",
    machine.StateFulfill,
    machine.ActionFulfilled,
    machine.StateCompleted,
  ],
  [
    "Reserve / Failure",
    machine.StateReserve,
    machine.ActionFailure,
    machine.StateFailed,
  ],
  [
    "Charge / Failure",
    machine.StateCharge,
    machine.ActionFailure,
    machine.StateReleaseReservation,
  ],
  [
    "Release Reservation / Released",
    machine.StateReleaseReservation,
    machine.ActionReleased,
    machine.StateFailed,
  ],
  [
    "Fulfill / Failure",
    machine.StateFulfill,
    machine.ActionFailure,
    machine.StateRefund,
  ],
  [
    "Refund / Refunded",
    machine.StateRefund,
    machine.ActionRefunded,
    machine.StateReleaseReservation,
  ],
] as const;

test.each(transitions)("%s", (_name, source, type, state) => {
  expect(machine.handle({}, source, { type })).toStrictEqual({
    ok: true,
    state,
  });
});

const rejectedTransitions = [
  ["Reserve / Fulfilled", machine.StateReserve, machine.ActionFulfilled],
  ["Completed / Failure", machine.StateCompleted, machine.ActionFailure],
  ["Failed / Reserved", machine.StateFailed, machine.ActionReserved],
] as const;

test.each(rejectedTransitions)("%s has no transition", (_name, state, type) => {
  const result = machine.handle({}, state, { type });
  expect(result).toStrictEqual({ ok: false, state, error: ErrNoTransition });
  if (!result.ok) expect(result.error).toBe(ErrNoTransition);
});
