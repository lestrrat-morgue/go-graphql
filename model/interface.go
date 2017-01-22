package model

type Definition interface{}
type Document interface {
	Definitions() chan Definition
	AddDefinitions(...Definition)
}
type document struct {
	definitions DefinitionList
	types       TypeList
}

type OperationType string

const (
	OperationTypeQuery    OperationType = "query"
	OperationTypeMutation OperationType = "mutation"
)

type OperationDefinition interface {
	OperationType() OperationType
	HasName() bool
	Name() string
	SetName(string)
	Variables() chan *VariableDefinition
	Directives() chan *Directive
	Selections() chan Selection
	AddVariableDefinitions(...*VariableDefinition)
	AddDirectives(...*Directive)
	AddSelections(...Selection)
}

type operationDefinition struct {
	typ        OperationType
	hasName    bool
	name       string
	variables  VariableDefinitionList
	directives DirectiveList
	selections SelectionSet
}

type FragmentDefinition struct {
	nameComponent
	typeComponent
	directives DirectiveList
	selections SelectionSet
}

type Type interface {
	IsNullable() bool
	SetNullable(bool)
}

type NamedType interface {
	Name() string
	Type
}

type namedType struct {
	kindComponent
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
type ObjectField interface {
	Name() string
	Value() Value
	SetValue(Value)
}

type objectField struct {
	nameComponent
	valueComponent
}

type ObjectValue interface {
	Value

	Fields() chan ObjectField
	AddFields(...ObjectField)
}
type objectValue struct {
	fields ObjectFieldList
}

type Selection interface{}

type SelectionSet []Selection

type Argument struct {
	nameComponent
	valueComponent
}

type Directive struct {
	name      string
	arguments ArgumentList
}

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
	typ        NamedType
}

// ObjectDefinition is a definition of a new object type
type ObjectDefinition interface {
	Type
	AddFields(...ObjectFieldDefinition)
	Fields() chan ObjectFieldDefinition
	Name() string
	HasImplements() bool
	Implements() NamedType
	SetImplements(NamedType)
}

type objectDefinition struct {
	nullable
	nameComponent
	fields        ObjectFieldDefinitionList
	hasImplements bool
	implements    NamedType
}

type ObjectFieldArgumentDefinition interface {
	Name() string
	Type() Type
	HasDefaultValue() bool
	DefaultValue() Value
	SetDefaultValue(Value)
}

type objectFieldArgumentDefinition struct {
	nameComponent
	typeComponent
	defaultValueComponent
}

type ObjectFieldDefinition interface {
	Name() string
	Type() Type
	Arguments() chan ObjectFieldArgumentDefinition
	AddArguments(...ObjectFieldArgumentDefinition)
}

type objectFieldDefinition struct {
	nameComponent
	typeComponent
	arguments ObjectFieldArgumentDefinitionList
}

type EnumDefinition struct {
	nullable // is this kosher?
	nameComponent
	elements EnumElementDefinitionList
}

type EnumElementDefinition struct {
	nameComponent
	valueComponent
}

type InterfaceDefinition struct {
	nullable
	nameComponent
	fields InterfaceFieldDefinitionList
}

type InterfaceFieldDefinition struct {
	nameComponent
	typeComponent
}

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

type Schema struct {
	query ObjectDefinition // But must be a query
	types ObjectDefinitionList
}
