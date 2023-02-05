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
provided under the [`examples`](examples) directory:
  - [boundary](example/boundary): A simple, three-state state machine that
    demonstrates the use of boundary (enter and exit) handlers.
  - [garagedoor](example/garagedoor): A state machine that simulates the
    functionality of a garage door.
  - [queue](example/queue): A state machine that encapsulates a queue object
    that can be closed only when the queue is empty, demonstrating the use of
    guard functions.
  - [simple](example/simple) A simple, three-state state machine containing a
    single action that cycles between the three states.

## Features

### Action hooks
The state machine configuration allows optionally specifying several three types
of hooks when handling an action, which are called in the following order:

1.  Exit: called when exiting a state or a group. The state exit handler is
    called first, followed by any applicable group exit handlers in the reverse
    order that the groups are specified.
1.  Transition: called for handling any transition logic. The handler is
    provided any action data, if specified.
1.  Enter: called when entering a state or a group. The group enter handlers are
    called first in the order that the groups are specified, followed by the
    state enter handler.

Each of the hooks are passed the context object, if specified. Also each hook
may return a non-nil error, which is propagated to the initial call to handle
the action, signaling that the transition should be aborted.

### Guard functions
Each transition can be preceeded by a guard function that determines whether or
not the transition should be taken given any context or action data, if
spedified. 

### Visualization
The renderer provides the option to output a
[mermaid.js](https://mermaid.js.org/) diagram in a markdown file. If checked
into GitHub, these diagrams will automatically be rendered.
