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

type DirectivesContainer interface {
	Directives() chan Directive
	AddDirectives(...Directive)
}

type SelectionsContainer interface {
	Selections() chan Selection
	AddSelections(...Selection)
}

type Type interface{}

type Definition interface {
	Namer
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
	Namer
	DirectivesContainer
	SelectionsContainer

	OperationType() OperationType
	HasName() bool
	SetName(string)
	Variables() chan VariableDefinition
	AddVariableDefinitions(...VariableDefinition)
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
	DirectivesContainer
	SelectionsContainer
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

type InterfaceDefinition interface {
	Nullable
	Namer
	Fields() chan InterfaceFieldDefinition
	AddFields(...InterfaceFieldDefinition)
}

type interfaceDefinition struct {
	nullable
	nameComponent
	fields InterfaceFieldDefinitionList
}

type InterfaceFieldDefinition interface {
	Namer
	Typer
}

type interfaceFieldDefinition struct {
	nameComponent
	typeComponent
}

type InputDefinition interface {
	Namer
	Fields() chan InputFieldDefinition
	AddFields(...InputFieldDefinition)
}

type inputDefinition struct {
	nameComponent
	fields InputFieldDefinitionList
}

type InputFieldDefinition interface {
	Namer
	Typer
}

type inputFieldDefinition struct {
	nameComponent
	typeComponent
}

type NamedType interface {
	Nullable
	Namer
}

type namedType struct {
	kindComponent
	nullable
	nameComponent
}

type ListType interface {
	Nullable
	Type() Type
}

type listType struct {
	nullable
	typeComponent
}

type Value interface {
	Kind() Kind
	Value() interface{}
}

type Variable interface {
	Namer
	Value
}

type variable struct {
	nameComponent
}

type intValue struct {
	value int
}

type floatValue struct {
	value float64
}

type stringValue struct {
	value string
}

type boolValue struct {
	value bool
}

type nullValue struct{}

type enumValue struct {
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

type SelectionField interface {
	Namer
	DirectivesContainer
	SelectionsContainer

	HasAlias() bool
	Alias() string
	SetAlias(string)
	Arguments() chan Argument
	AddArguments(...Argument)
}

type selectionField struct {
	nameComponent
	hasAlias   bool
	alias      string
	arguments  ArgumentList
	directives DirectiveList
	selections SelectionList
}

type FragmentSpread interface {
	Namer
	DirectivesContainer
}

type fragmentSpread struct {
	nameComponent
	directives DirectiveList
}

type InlineFragment interface {
	DirectivesContainer
	SelectionsContainer

	SetTypeCondition(NamedType)
	TypeCondition() NamedType
}

type inlineFragment struct {
	directives DirectiveList
	selections SelectionList
	typ        NamedType
}

type UnionDefinition interface {
	Namer
	Types() chan Type
	AddTypes(...Type)
}

type unionDefinition struct {
	nameComponent
	types TypeList
}

type Schema interface {
	Namer
	Query() NamedType
	SetQuery(NamedType)
	Mutation() NamedType
	SetMutation(NamedType)
	Subscription() NamedType
	SetSubscription(NamedType)
	Types() chan NamedType
	AddTypes(...NamedType)
	Directives() chan string
	AddDirectives(...string)
}

type schema struct {
	query        NamedType
	types        NamedTypeList
	mutation     NamedType
	subscription NamedType
}
