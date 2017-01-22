package model

// Namer represents all those that have a name to share
type Namer interface {
	Name() string
}

// Typer represents all those that can get/set a type
type Typer interface {
	Type() Type
	SetType(Type)
}

// DefaultValuer represents all those that may or may not have
// a default value associated with it
type DefaultValuer interface {
	HasDefaultValue() bool
	DefaultValue() Value
	SetDefaultValue(Value)
}

// Nullable represents those types that can be specified that
// it could be null or not-null
type Nullable interface {
	IsNullable() bool
	SetNullable(bool)
}

type Type interface {}

type Definition interface{
}

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
	Variables() chan VariableDefinition
	Directives() chan Directive
	Selections() chan Selection
	AddVariableDefinitions(...VariableDefinition)
	AddDirectives(...Directive)
	AddSelections(...Selection)
}

type operationDefinition struct {
	typ        OperationType
	hasName    bool
	name       string
	variables  VariableDefinitionList
	directives DirectiveList
	selections SelectionList
}

type FragmentDefinition interface {
	Namer
	Typer

	Directives() chan Directive
	Selections() chan Selection
	AddDirectives(...Directive)
	AddSelections(...Selection)
}

type fragmentDefinition struct {
	nameComponent
	typeComponent
	directives DirectiveList
	selections SelectionList
}

type VariableDefinition interface {
	Namer
	Typer
	DefaultValuer
}

type variableDefinition struct {
	nameComponent
	typeComponent
	defaultValueComponent
}

// ObjectDefinition is a definition of a new object type
type ObjectDefinition interface {
	Namer
	Type
	Nullable
	AddFields(...ObjectFieldDefinition)
	Fields() chan ObjectFieldDefinition
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
	Namer
	Typer
	DefaultValuer
}

type objectFieldArgumentDefinition struct {
	nameComponent
	typeComponent
	defaultValueComponent
}

type ObjectFieldDefinition interface {
	Namer
	Typer
	Arguments() chan ObjectFieldArgumentDefinition
	AddArguments(...ObjectFieldArgumentDefinition)
}

type objectFieldDefinition struct {
	nameComponent
	typeComponent
	arguments ObjectFieldArgumentDefinitionList
}

type EnumDefinition interface {
	Namer
	Elements() chan EnumElementDefinition
	AddElements(...EnumElementDefinition)
}

type enumDefinition struct {
	nullable // is this kosher?
	nameComponent
	elements EnumElementDefinitionList
}

type EnumElementDefinition interface {
	Namer
	Value() Value
}
	
type enumElementDefinition struct {
	nameComponent
	valueComponent
}

type NamedType interface {
	Nullable
	Namer
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
	Namer
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

type Argument interface {
	Namer
	Value() Value
}

type argument struct {
	nameComponent
	valueComponent
}

type Directive interface {
	Namer
	Arguments() chan Argument
	AddArguments(...Argument)
}

type directive struct {
	name      string
	arguments ArgumentList
}

type Field struct {
	nameComponent
	hasAlias   bool
	alias      string
	arguments  ArgumentList
	directives DirectiveList
	selections SelectionList
}

type FragmentSpread struct {
	nameComponent
	directives DirectiveList
}

type InlineFragment struct {
	directives DirectiveList
	selections SelectionList
	typ        NamedType
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
