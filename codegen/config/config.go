package config

type Machine struct {
	// Required. Code generation features.
	Codegen Codegen `yaml:"codegen"`
	// Required. State machine actions.
	Actions Actions `yaml:"actions"`
	// Optional. State machine groups.
	Groups Groups `yaml:"groups"`
	// Required. State machine states.
	States States `yaml:"states"`
}

type Codegen struct {
	// Required. Golang code generation features.
	Golang Golang `yaml:"golang"`
	// Optional. Mermaid code generation features.
	Mermaid Mermaid `yaml:"mermaid"`
}

type Golang struct {
	// Required. The golang package name for the generated file.
	Package string `yaml:"package"`
	// Optional. Additional imports.
	Imports []Import `yaml:"imports"`
	// Optional. The type for a context parameter in all of the handler methods.
	Context string `yaml:"context"`
	// Optional. Create a struct declaring the action types.
	DeclareActions bool `yaml:"declare_actions"`
	// Optional. Create a struct declaring the state types.
	DeclareStates bool `yaml:"declare_states"`
}

type Import struct {
	// Optional. The import name.
	Name string `yaml:"name"`
	// Required. The import path.
	Path string `yaml:"path"`
}

type Mermaid struct {
	// Optional. A flag to enable generating a mermaid diagram. All other fields
	// are ignored if this is set to false.
	Enabled bool `yaml:"enabled"`
	// Required. The name of the output markdown file.
	Filename string `yaml:"filename"`
}

type Actions struct {
	// Optional. The type of the action enum. If omitted, the type "ActionType" is
	// used.
	Type string `yaml:"type"`
	// Required. A list containing all of the actions.
	Values []Action `yaml:"values"`
}

type Action struct {
	// Required. The action name. This is used for generating handler methods for
	// the action.
	Name string `yaml:"name"`
	// Optional. The action symbol. If omitted, the symbol "Action" appended with
	// the action name is used.
	Symbol string `yaml:"symbol"`
	// Optional. The action data type. If omitted, no action data parameter is
	// created.
	DataType string `yaml:"data_type"`
}

type Groups struct {
	// Required. A list of all of the groups. Containing groups must appear before
	// their subgroups. Entry hooks run in this order, exit hooks in reverse order,
	// and action rules are considered in reverse order after state rules.
	Values []Group `yaml:"values"`
}

type Group struct {
	// Required. The group name. State and group names share a namespace.
	Name string `yaml:"name"`
	// Optional unless the group is a destination or another group's initial
	// target. The name of a member state or subgroup, directly or transitively.
	// Group targets are resolved recursively to their initial state.
	Initial string `yaml:"initial"`
	// Optional. The name of the function called when the group is entered. If
	// omitted, no function is created.
	Enter string `yaml:"enter"`
	// Optional. The name of the function called when the group is exited. If
	// omitted, no function is created.
	Exit string `yaml:"exit"`
	// Optional. Direct subgroups. Their transitive member states are also members
	// of this group. Containment must be acyclic; membership may overlap.
	Groups []string `yaml:"groups"`
	// Optional. Direct member states.
	States []string `yaml:"states"`
	// Optional. Action rules available to every member state. These are considered
	// after state rules, in reverse group declaration order, until a guard passes
	// or an unguarded transition is reached. Guard errors abort the transition.
	TransitionActions []TransitionAction `yaml:"actions"`
}

type States struct {
	// Optional. The type of the state enum. If omitted, the type "StateType" is
	// used.
	Type string `yaml:"type"`
	// Required. A list containing all of the states.
	Values []State `yaml:"values"`
}

type State struct {
	// Required. The state name. State and group names share a namespace.
	Name string `yaml:"name"`
	// Optional. The state symbol. If omitted, the symbol "State" appended with
	// the state name is used.
	Symbol string `yaml:"symbol"`
	// Optional. The name of the function called before the state is entered. If
	// omitted, no function is created.
	Enter string `yaml:"enter"`
	// Optional. The name of the function called when the state is exited. If
	// omitted, no function is created.
	Exit string `yaml:"exit"`
	// Optional. A list of actions that will cause transitions from the state, and
	// the corresponding transitions.
	TransitionActions []TransitionAction `yaml:"actions"`
}

type TransitionAction struct {
	// Required. The name of the action that triggers the transition.
	Action string `yaml:"action"`
	// Required. Transitions in evaluation order. The first passing guard or
	// unguarded transition wins. If every guard is false, evaluation continues
	// with the next applicable group's rules. A guard error aborts evaluation.
	Transitions []Transition `yaml:"transitions"`
}

type Transition struct {
	// Optional. Name of the transition function.
	Transition string `yaml:"transition"`
	// Required. Name of the destination state or group. A group destination
	// resolves through its initial values to a concrete state.
	Destination string `yaml:"destination"`
	// Optional. Groups whose boundaries are crossed even when both endpoints are
	// members. Includes their declared subgroups transitively, but not groups
	// that merely overlap. Each group is exited/entered at most once.
	External []string `yaml:"external"`
	// Optional. Name of the guard function.
	Guard string `yaml:"guard"`
}
