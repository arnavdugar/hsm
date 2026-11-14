import { expect, test } from 'vitest';
import { ErrNoTransition } from '@arnavdugar/hsm/runtime';
import * as machine from './machine.ts';

const transitions = [
  ['Closed / Active Open', machine.StateClosed, machine.ActionActiveOpen, machine.StateSynSent],
  ['Syn Sent / Receive Syn Ack', machine.StateSynSent, machine.ActionReceiveSynAck, machine.StateEstablished],
  ['Established / Close', machine.StateEstablished, machine.ActionClose, machine.StateFinWait1],
  ['Fin Wait1 / Receive Ack', machine.StateFinWait1, machine.ActionReceiveAck, machine.StateFinWait2],
  ['Fin Wait2 / Receive Fin', machine.StateFinWait2, machine.ActionReceiveFin, machine.StateTimeWait],
  ['Time Wait / Timeout', machine.StateTimeWait, machine.ActionTimeout, machine.StateClosed],
  ['Closed / Passive Open', machine.StateClosed, machine.ActionPassiveOpen, machine.StateListen],
  ['Listen / Receive Syn', machine.StateListen, machine.ActionReceiveSyn, machine.StateSynReceived],
  ['Syn Received / Receive Ack', machine.StateSynReceived, machine.ActionReceiveAck, machine.StateEstablished],
  ['Established / Receive Fin', machine.StateEstablished, machine.ActionReceiveFin, machine.StateCloseWait],
  ['Close Wait / Close', machine.StateCloseWait, machine.ActionClose, machine.StateLastAck],
  ['Last Ack / Receive Ack', machine.StateLastAck, machine.ActionReceiveAck, machine.StateClosed],
  ['Syn Sent / Receive Syn', machine.StateSynSent, machine.ActionReceiveSyn, machine.StateSynReceived],
  ['Fin Wait1 / Receive Fin', machine.StateFinWait1, machine.ActionReceiveFin, machine.StateClosing],
  ['Closing / Receive Ack', machine.StateClosing, machine.ActionReceiveAck, machine.StateTimeWait],
  ['Fin Wait1 / Receive Fin Ack', machine.StateFinWait1, machine.ActionReceiveFinAck, machine.StateTimeWait],
  ['Listen / Close', machine.StateListen, machine.ActionClose, machine.StateClosed],
  ['Listen / Send', machine.StateListen, machine.ActionSend, machine.StateSynSent],
  ['Syn Sent / Close', machine.StateSynSent, machine.ActionClose, machine.StateClosed],
  ['Syn Received / Close', machine.StateSynReceived, machine.ActionClose, machine.StateFinWait1],
] as const;

test.each(transitions)('%s', (_name, source, type, state) => {
  expect(machine.handle({}, source, { type })).toStrictEqual({ ok: true, state });
});

const rejectedTransitions = [
  ['Closed / Timeout', machine.StateClosed, machine.ActionTimeout],
] as const;

test.each(rejectedTransitions)('%s has no transition', (_name, state, type) => {
  const result = machine.handle({}, state, { type });
  expect(result).toStrictEqual({ ok: false, state, error: ErrNoTransition });
  if (!result.ok) expect(result.error).toBe(ErrNoTransition);
});
