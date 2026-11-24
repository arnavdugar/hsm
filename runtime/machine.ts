export const ErrNoTransition = new Error("no transition");
export const ErrUnknownAction = new Error("unknown action");
export const ErrUnknownState = new Error("unknown state");

export type HandleResult<StateType> = {
    ok: true;
    state: StateType;
} | {
    ok: false;
    state: StateType;
    error: unknown;
};
