package model

func NewInlineFragment() *InlineFragment {
	return &InlineFragment{}
}

func (f *InlineFragment) AddSelections(list ...Selection) {
	f.selections.Add(list...)
}

func (f *InlineFragment) AddDirectives(list ...*Directive) {
	f.directives.Add(list...)
}

func (f *InlineFragment) SetTypeCondition(typ *NamedType) {
	f.typ = typ
}

func (f InlineFragment) Type() *NamedType {
	return f.typ
}

func (f InlineFragment) SelectionSet() chan Selection {
	return f.selections.Iterator()
}

func (f InlineFragment) Directives() chan *Directive {
	return f.directives.Iterator()
}

func NewFragmentDefinition(name string, typ *NamedType) *FragmentDefinition {
	return &FragmentDefinition{
		nameComponent: nameComponent(name),
		typ:           typ,
	}
}

func (f FragmentDefinition) SelectionSet() chan Selection {
	return f.selections.Iterator()
}

func (f *FragmentDefinition) AddSelections(selections ...Selection) {
	f.selections.Add(selections...)
}

func (def FragmentDefinition) Type() *NamedType {
	return def.typ
}

func (def *FragmentDefinition) AddDirectives(list ...*Directive) {
	def.directives.Add(list...)
}

func NewOperationDefinition(typ OperationType) *OperationDefinition {
	return &OperationDefinition{
		typ: typ,
	}
}

func (def OperationDefinition) Type() OperationType {
	return def.typ
}

func (def OperationDefinition) VariableDefinitions() chan *VariableDefinition {
	ch := make(chan *VariableDefinition, len(def.vardefs))
	for _, vdef := range def.vardefs {
		ch <- vdef
	}
	close(ch)
	return ch
}

func (def OperationDefinition) SelectionSet() chan Selection {
	return def.selections.Iterator()
}

func (def OperationDefinition) HasName() bool {
	return def.hasName
}

func (def OperationDefinition) Name() string {
	return def.name
}

func (def *OperationDefinition) SetName(s string) {
	def.hasName = true
	def.name = s
}

func (def *OperationDefinition) AddVariableDefinitions(list ...*VariableDefinition) {
	def.vardefs = append(def.vardefs, list...)
}

func (def *OperationDefinition) AddDirectives(list ...*Directive) {
	def.directives = append(def.directives, list...)
}

func (def *OperationDefinition) AddSelections(list ...Selection) {
	def.selections = append(def.selections, list...)
}

func NewVariableDefinition(name string, typ Type) *VariableDefinition {
	return &VariableDefinition{
		nameComponent: nameComponent(name),
		typeComponent: typeComponent{ typ: typ },
	}
}
