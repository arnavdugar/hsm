import { expect, test } from "vitest";
import * as machine from "./machine.ts";

const transitions = [
  ["A / Next", machine.StateA, machine.ActionNext, machine.StateB],
  ["B / Next", machine.StateB, machine.ActionNext, machine.StateC],
  ["C / Next", machine.StateC, machine.ActionNext, machine.StateA],
] as const;

test.each(transitions)("%s", (_name, source, type, state) => {
  expect(machine.handle({}, source, { type })).toStrictEqual({
    ok: true,
    state,
  });
});
