package model

func (s *SelectionSet) Add(list ...Selection) {
	*s = append(*s, list...)
}

func (s SelectionSet) Iterator() chan Selection {
	ch := make(chan Selection, len(s))
	for _, sel := range s {
		ch <- sel
	}
	close(ch)
	return ch
}

func NewField(n string) *Field {
	return &Field{
		nameComponent: nameComponent(n),
	}
}

func (f Field) HasAlias() bool {
	return f.hasAlias
}

func (f Field) Alias() string {
	return f.alias
}

func (f *Field) SetAlias(s string) {
	f.hasAlias = true
	f.alias = s
}

func (f Field) Arguments() chan *Argument {
	return f.arguments.Iterator()
}

func (f Field) Directives() chan *Directive {
	return f.directives.Iterator()
}

func (f Field) SelectionSet() chan Selection {
	return f.selections.Iterator()
}

func (f *Field) AddArguments(args ...*Argument) {
	f.arguments.Add(args...)
}

func (f *Field) AddDirectives(directives ...*Directive) {
	f.directives.Add(directives...)
}

func (f *Field) AddSelections(selections ...Selection) {
	f.selections.Add(selections...)
}

func NewFragmentSpread(name string) *FragmentSpread {
	return &FragmentSpread{
		nameComponent: nameComponent(name),
	}
}

func (f *FragmentSpread) AddDirectives(directives ...*Directive) {
	f.directives.Add(directives...)
}
