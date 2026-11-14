import { expect, test } from "vitest";
import { ErrNoTransition } from "@arnavdugar/hsm/runtime";
import * as machine from "./machine.ts";

const transitions = [
  [
    "Menu / Check Balance",
    machine.StateMenu,
    machine.ActionCheckBalance,
    machine.StateBalance,
  ],
  [
    "Balance / Done",
    machine.StateBalance,
    machine.ActionDone,
    machine.StateMenu,
  ],
  [
    "Menu / Withdraw",
    machine.StateMenu,
    machine.ActionWithdraw,
    machine.StateEnterAmount,
  ],
  [
    "Enter Amount / Submit",
    machine.StateEnterAmount,
    machine.ActionSubmit,
    machine.StateAuthorizeWithdrawal,
  ],
  [
    "Authorize Withdrawal / Approve",
    machine.StateAuthorizeWithdrawal,
    machine.ActionApprove,
    machine.StateDispenseCash,
  ],
  [
    "Dispense Cash / Done",
    machine.StateDispenseCash,
    machine.ActionDone,
    machine.StateMenu,
  ],
  [
    "Menu / Deposit",
    machine.StateMenu,
    machine.ActionDeposit,
    machine.StateAcceptDeposit,
  ],
  [
    "Accept Deposit / Done",
    machine.StateAcceptDeposit,
    machine.ActionDone,
    machine.StateMenu,
  ],
  ["Menu / Exit", machine.StateMenu, machine.ActionExit, machine.StateEnded],
  [
    "Authorize Withdrawal / Decline",
    machine.StateAuthorizeWithdrawal,
    machine.ActionDecline,
    machine.StateMenu,
  ],
  [
    "Enter Amount / Cancel",
    machine.StateEnterAmount,
    machine.ActionCancel,
    machine.StateMenu,
  ],
  [
    "Accept Deposit / Cancel",
    machine.StateAcceptDeposit,
    machine.ActionCancel,
    machine.StateMenu,
  ],
] as const;

test.each(transitions)("%s", (_name, source, type, state) => {
  expect(machine.handle({}, source, { type })).toStrictEqual({
    ok: true,
    state,
  });
});

const rejectedTransitions = [
  ["Menu / Approve", machine.StateMenu, machine.ActionApprove],
  ["Ended / Withdraw", machine.StateEnded, machine.ActionWithdraw],
] as const;

test.each(rejectedTransitions)("%s has no transition", (_name, state, type) => {
  const result = machine.handle({}, state, { type });
  expect(result).toStrictEqual({ ok: false, state, error: ErrNoTransition });
  if (!result.ok) expect(result.error).toBe(ErrNoTransition);
});
