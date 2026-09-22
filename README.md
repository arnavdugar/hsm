# Hierarchical State Machine

A hierarchical state machine (HSM) is a representation of a
[finite-state machine](https://en.wikipedia.org/wiki/Finite-state_machine)
where states can be grouped together into higher-level states. In a traditional
implementation, the hierarchy of states is structured as a tree, where all
substates have exactly one parent state and parent states have disjoint sets of
substates. This implementation, however, is more general: states may be
arbitrarily grouped together, and groups of states may contain non-disjoint sets
of subgroups and substates.

Many hierarchical state machine provide a library to declare the structure of a
state machine at runtime, and evaluate the state machine as each action is
handled. This implementation uses a yaml configuration file to statically render
handler methods that perform the state machine logic.

## Status

This project is still under active development and is subject to breaking
changes without warning.

## Usage

### Configuration

The state machine structure and code generation features are configured using a
yaml file, with structure defined in
[codegen/config/config.go](codegen/config/config.go).

### Code generation

The logic for the state machine can be generated using the command

```
go run github.com/arnavdugar/hsm/codegen -i=${INPUT_FILE} -o=${OUTPUT_FILE}
```

where `INPUT_FILE` is the configuration file and `OUTPUT_FILE` is the generated
golang file. An additional file containing a mermaid.js diagram of the state
machine can be generated if specified in the configuration file.

### Examples

Several examples of state machine configurations and their rendered logic are
provided under the [`example`](example) directory:

- [boundary](example/boundary): A simple, three-state state machine that
  demonstrates the use of boundary (enter and exit) handlers.
- [group](example/group): Overlapping groups with nested initial targets, shared
  actions, guard fallback, and external transitions that reenter group boundaries.
- [garagedoor](example/garagedoor): A state machine that simulates the
  functionality of a garage door.
- [queue](example/queue): A state machine that encapsulates a queue object
  that can be closed only when the queue is empty, demonstrating the use of
  guard functions.
- [simple](example/simple) A simple, three-state state machine containing a
  single action that cycles between the three states.
- [orderfulfillment](example/orderfulfillment): Shipping and pickup branches
  reconverge at completion, with a separate cancellation terminal state.
- [atm](example/atm): A menu hub with balance, withdrawal, and deposit spokes.
- [combinationlock](example/combinationlock): A prefix chain with guards that
  advance or reset based on digit action data.
- [connectionmanager](example/connectionmanager): Shared recovery loops with
  a guarded retry budget and explicit timeout actions.
- [rollback](example/rollback): A forward job workflow with compensation paths
  determined by how far the job progressed.
- [tcp](example/tcp): A simplified TCP lifecycle with asymmetric opening and
  closing paths, including simultaneous open and close.

## Features

### Groups and initial targets

Groups are optional. Each group lists its direct member `states` and `groups`;
all states belonging to its subgroups are also members transitively. Membership
may overlap, and a subgroup may belong to multiple groups. Containment must be
acyclic, and containing groups must be declared before their subgroups in
`groups.values`. State and group names share one namespace and must be unique.

```yaml
groups:
  values:
    - name: Session
      groups: [Connecting]
      states: [Connected]
      initial: Connecting
      enter: OpenSession
      exit: CloseSession
    - name: Connecting
      states: [Dialing, Handshaking]
      initial: Dialing
```

A transition's `destination` may name a state or a group. A group destination
resolves recursively through `initial` values to a concrete state. An `initial`
value must name a member state or subgroup, directly or transitively. Every group
followed during resolution must have an initial value. Groups that are not used
as destinations or initial targets may omit it.

Initial values are used only to resolve group targets. A transition directly to
`Handshaking` enters that state, regardless of the initial values of its groups.
The resulting state's membership determines every active group, including groups
outside the initial-resolution path.

### Group boundaries and external transitions

Normally, a transition exits groups containing its source but not its destination,
and enters groups containing its destination but not its source. Groups containing
both remain active without exit or entry hooks.

The optional `external` list forces applicable boundaries of the named groups to
run even when both states belong to them:

```yaml
transitions:
  - destination: Session
    external: [Session]
```

Externality includes explicitly declared subgroups transitively. It does not
propagate to groups that merely overlap or happen to contain the same states.
Only groups containing the source are exited, and only groups containing the
resolved destination are entered. Each group is exited or entered at most once,
even when multiple membership paths or external declarations refer to it.
Omitting `external`, or specifying `external: []`, uses normal membership changes.

### Shared actions and guard fallback

Groups may declare `actions` using the same schema as states. When handling an
action, state rules are considered first, followed by rules from all containing
groups in reverse group declaration order. Within each rule, transitions are
evaluated in their listed order. Each containing group is considered once, even
if the state belongs to it through multiple paths.

The first transition with a passing guard, or with no guard, is selected. False
guards continue evaluation, including falling through from state rules to group
rules. A guard error aborts evaluation without running boundary hooks or trying
any fallback. If no transition is selected, the handler returns `ErrNoTransition`.

### Action hooks

The state machine configuration allows optionally specifying three types
of hooks when handling an action, which are called in the following order:

1.  Exit: called when exiting a state or a group. The state exit handler is
    called first, followed by any applicable group exit handlers in the reverse
    order that the groups are specified.
1.  Transition: called for handling any transition logic. The handler is
    provided any action data, if specified.
1.  Enter: called when entering a state or a group. The group enter handlers are
    called first in the order that the groups are specified, followed by the
    state enter handler.

State exit and entry hooks run for every selected transition, including
self-transitions. The `external` list controls group boundary hooks only.

Each of the hooks are passed the context object, if specified. Also each hook
may return a non-nil error, which is propagated to the initial call to handle
the action, signaling that the transition should be aborted.

### Guard functions

Each transition can be preceded by a guard function that determines whether or
not the transition should be taken given any context or action data, if
specified.

### Visualization

The renderer provides the option to output a
[mermaid.js](https://mermaid.js.org/) diagram in a markdown file. If checked
into GitHub, these diagrams will automatically be rendered.
