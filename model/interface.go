package model

type TypeList []Type
type Document struct {
	definitions []Definition
	types       TypeList
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
	nameComponent
	typ        *NamedType
	directives DirectiveList
	selections SelectionSet
}

type Type interface {
	IsNullable() bool
	SetNullable(bool)
}

type NamedType struct {
	nullable
	nameComponent
}

type ListType struct {
	nullable
	typeComponent
}

type VariableDefinition struct {
	nameComponent
	typeComponent
	defaultValueComponent
}
type VariableDefinitionList []*VariableDefinition

type Value interface {
	Value() interface{}
}

type Variable struct {
	nameComponent
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
	nameComponent
}

// ObjectField represents a literal object's field (NOT a type)
type ObjectField struct {
	nameComponent
	value Value
}
type ObjectFieldList []*ObjectField

type ObjectValue struct {
	fields ObjectFieldList
}

type Selection interface{}

type SelectionSet []Selection

type Argument struct {
	nameComponent
	value Value
}
type ArgumentList []*Argument

type Directive struct {
	name      string
	arguments ArgumentList
}
type DirectiveList []*Directive

type Field struct {
	nameComponent
	hasAlias   bool
	alias      string
	arguments  ArgumentList
	directives DirectiveList
	selections SelectionSet
}

type FragmentSpread struct {
	nameComponent
	directives DirectiveList
}

type InlineFragment struct {
	directives DirectiveList
	selections SelectionSet
	typ        *NamedType
}


// ObjectDefinition is a definition of a new object type
type ObjectDefinition struct {
	nameComponent
	fields        ObjectFieldDefinitionList
	hasImplements bool
	implements    string
}

type ObjectFieldArgumentDefinition struct {
	nameComponent
	typeComponent
	defaultValueComponent
}

type ObjectFieldArgumentDefinitionList []*ObjectFieldArgumentDefinition

type ObjectFieldDefinition struct {
	nameComponent
	typeComponent
	arguments ObjectFieldArgumentDefinitionList
}

type ObjectFieldDefinitionList []*ObjectFieldDefinition

type EnumDefinition struct {
	nameComponent
	elements EnumElementList
}

type EnumElement struct {
	nameComponent
}
type EnumElementList []*EnumElement

type InterfaceDefinition struct {
	nameComponent
	fields InterfaceFieldList
}

type InterfaceField struct {
	nameComponent
	typeComponent
}

type InterfaceFieldList []*InterfaceField

type UnionDefinition struct {
	nameComponent
	types TypeList
}

type InputDefinition struct {
	nameComponent
	fields InputFieldDefinitionList
}

type InputFieldDefinition struct {
	nameComponent
	typeComponent
}

type InputFieldDefinitionList []*InputFieldDefinition
