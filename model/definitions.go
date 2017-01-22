package model

func NewInlineFragment() InlineFragment {
	return &inlineFragment{}
}

func (f *inlineFragment) AddSelections(list ...Selection) {
	f.selections.Add(list...)
}

func (f *inlineFragment) AddDirectives(list ...Directive) {
	f.directives.Add(list...)
}

func (f *inlineFragment) SetTypeCondition(typ NamedType) {
	f.typ = typ
}

func (f inlineFragment) TypeCondition() NamedType {
	return f.typ
}

func (f inlineFragment) Selections() chan Selection {
	return f.selections.Iterator()
}

func (f inlineFragment) Directives() chan Directive {
	return f.directives.Iterator()
}

func NewFragmentDefinition(name string, typ NamedType) FragmentDefinition {
	return &fragmentDefinition{
		nameComponent: nameComponent(name),
		typeComponent: typeComponent{typ: typ},
	}
}

func (f fragmentDefinition) Selections() chan Selection {
	return f.selections.Iterator()
}

func (f *fragmentDefinition) AddSelections(selections ...Selection) {
	f.selections.Add(selections...)
}

func (f fragmentDefinition) Directives() chan Directive {
	return f.directives.Iterator()
}

func (f *fragmentDefinition) AddDirectives(list ...Directive) {
	f.directives.Add(list...)
}

func NewOperationDefinition(typ OperationType) OperationDefinition {
	return &operationDefinition{
		typ: typ,
	}
}

func (def operationDefinition) OperationType() OperationType {
	return def.typ
}

func (def operationDefinition) Variables() chan VariableDefinition {
	return def.variables.Iterator()
}

func (def operationDefinition) Selections() chan Selection {
	return def.selections.Iterator()
}

func (def operationDefinition) Directives() chan Directive {
	return def.directives.Iterator()
}

func (def operationDefinition) HasName() bool {
	return def.hasName
}

func (def operationDefinition) Name() string {
	return def.name
}

func (def *operationDefinition) SetName(s string) {
	def.hasName = true
	def.name = s
}

func (def *operationDefinition) AddVariableDefinitions(list ...VariableDefinition) {
	def.variables.Add(list...)
}

func (def *operationDefinition) AddDirectives(list ...Directive) {
	def.directives.Add(list...)
}

func (def *operationDefinition) AddSelections(list ...Selection) {
	def.selections.Add(list...)
}

func NewVariableDefinition(name string, typ Type) VariableDefinition {
	return &variableDefinition{
		nameComponent: nameComponent(name),
		typeComponent: typeComponent{ typ: typ },
	}
}
