package config

type Machine struct {
	// Required. Code generation features.
	Codegen Codegen `yaml:"codegen"`
	// Required. State machine actions.
	Actions Actions `yaml:"actions"`
	// Optional. A list of groups containing states or other groups.
	Groups []Group `yaml:"groups"`
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

type Group struct {
	// Required. The group name.
	Name string `yaml:"name"`
	// Optional. The name of the function called before the state is entered. If
	// omitted, no function is created.
	Enter string `yaml:"enter"`
	// Optional. The name of the function called when the state is exited. If
	// omitted, no function is created.
	Exit string `yaml:"exit"`
	// Optional. A list of group names of groups contained by this group.
	Groups []string `yaml:"groups"`
	// Optional. A list of state names of groups contained by this group.
	States []string `yaml:"states"`
}

type States struct {
	// Optional. The type of the state enum. If omitted, the type "StateType" is
	// used.
	Type string `yaml:"type"`
	// Required. A list containing all of the states.
	Values []State `yaml:"values"`
}

type State struct {
	// Required. The state name.
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
	// Required. A list of actions that will cause transitions from the state, and
	// the corresponding transitions.
	TransitionActions []TransitionAction `yaml:"actions"`
}

type TransitionAction struct {
	// Required. The name of the action that triggers the transition.
	Action string `yaml:"action"`
	// Required. The transitions from the state. If multiple transitions are
	// specified for a given action, the transitions should have guard functions
	// specified to select the appropriate destination.
	Transitions []Transition `yaml:"transitions"`
}

type Transition struct {
	// Required. Name of the destination.
	Destination string `yaml:"destination"`
	// Optional. Name of the guard function.
	Guard string `yaml:"guard"`
	// Optional. Name of the transition function.
	Transition string `yaml:"transition"`
	// Optional. A list of groups that should be re-entered during the transition.
	// The source and destination states must both be contained in the groups.
	External []string `yaml:"external"`
}
