package model

func NewSelectionField(n string) SelectionField {
	return &selectionField{
		nameComponent: nameComponent(n),
	}
}

func (f selectionField) HasAlias() bool {
	return f.hasAlias
}

func (f selectionField) Alias() string {
	return f.alias
}

func (f *selectionField) SetAlias(s string) {
	f.hasAlias = true
	f.alias = s
}

func (f selectionField) Arguments() chan Argument {
	return f.arguments.Iterator()
}

func (f selectionField) Directives() chan Directive {
	return f.directives.Iterator()
}

func (f selectionField) Selections() chan Selection {
	return f.selections.Iterator()
}

func (f *selectionField) AddArguments(args ...Argument) {
	f.arguments.Add(args...)
}

func (f *selectionField) AddDirectives(directives ...Directive) {
	f.directives.Add(directives...)
}

func (f *selectionField) AddSelections(selections ...Selection) {
	f.selections.Add(selections...)
}

func NewFragmentSpread(name string) FragmentSpread {
	return &fragmentSpread{
		nameComponent: nameComponent(name),
	}
}

func (f fragmentSpread) Directives() chan Directive {
	return f.directives.Iterator()
}

func (f *fragmentSpread) AddDirectives(directives ...Directive) {
	f.directives.Add(directives...)
}
