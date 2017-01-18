package model

type Document struct {
	definitions []Definition
}

type Definition interface{}
type OperationType string

const (
	OperationTypeQuery    OperationType = "query"
	OperationTypeMutation OperationType = "mutation"
)

type OperationDefinition struct {
	typ        OperationType
	hasName    bool
	name       string
	vardefs    VariableDefinitionList
	directives DirectiveList
	selections SelectionSet
}

type FragmentDefinition struct {
	name       string
	typ        *NamedType
	directives DirectiveList
	selections SelectionSet
}

type Type interface {
	IsNullable() bool
	SetNullable(bool)
}

type nullable bool
type NamedType struct {
	nullable
	name string
}

type ListType struct {
	nullable
	typ Type
}

type VariableDefinition struct {
	name            string
	typ             Type
	hasDefaultValue bool
	defaultValue    Value
}
type VariableDefinitionList []*VariableDefinition

type Value interface {
	String() string
}

type Variable struct {
	name string
}

type IntValue struct {
	value int
}

type FloatValue struct {
	value float64
}

type StringValue struct {
	value string
}

type BoolValue struct {
	value bool
}

type NullValue struct{}

type EnumValue struct {
	name string
}
type Selection interface{}

type SelectionSet []Selection

type Argument struct {
	name  string
	value Value
}
type ArgumentList []*Argument

type Directive struct {
	name      string
	arguments ArgumentList
}
type DirectiveList []*Directive

type Field struct {
	hasAlias   bool
	alias      string
	arguments  ArgumentList
	directives DirectiveList
	name       string
	selections SelectionSet
}

type FragmentSpread struct {
	name       string
	directives DirectiveList
}

type InlineFragment struct {
	directives DirectiveList
	selections SelectionSet
	typ        *NamedType
}
