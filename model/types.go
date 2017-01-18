package model

func (n nullable) IsNullable() bool {
	return bool(n)
}

func (n *nullable) SetNullable(b bool) {
	*n = nullable(b)
}

func NewNamedType(name string) *NamedType {
	return &NamedType{
		name:     name,
		nullable: true,
	}
}

func (t *NamedType) Name() string {
	return t.name
}

func NewListType(t Type) *ListType {
	return &ListType{
		nullable: true,
		typ:      t,
	}
}

func (t *ListType) Type() Type {
	return t.typ
}

func (t *TypeList) Add(list ...Type) {
	*t = append(*t, list...)
}

func (t TypeList) Iterator() chan Type {
	ch := make(chan Type, len(t))
	for _, f := range t {
		ch <- f
	}
	close(ch)
	return ch
}

func (args *ObjectTypeDefinitionFieldArgumentList) Add(list ...*ObjectTypeDefinitionFieldArgument) {
	*args = append(*args, list...)
}

func (args ObjectTypeDefinitionFieldArgumentList) Iterator() chan *ObjectTypeDefinitionFieldArgument {
	ch := make(chan *ObjectTypeDefinitionFieldArgument, len(args))
	for _, arg := range args {
		ch <- arg
	}
	close(ch)
	return ch
}

func NewObjectTypeDefinitionFieldArgument(name string, typ Type) *ObjectTypeDefinitionFieldArgument {
	return &ObjectTypeDefinitionFieldArgument{
		name: name,
		typ:  typ,
	}
}

func (arg ObjectTypeDefinitionFieldArgument) Name() string {
	return arg.name
}

func (arg ObjectTypeDefinitionFieldArgument) Type() Type {
	return arg.typ
}

func (arg *ObjectTypeDefinitionFieldArgument) SetDefaultValue(v Value) {
	arg.hasDefaultValue = true
	arg.defaultValue = v
}

func (arg ObjectTypeDefinitionFieldArgument) HasDefaultValue() bool {
	return arg.hasDefaultValue
}

func (arg ObjectTypeDefinitionFieldArgument) DefaultValue() Value {
	return arg.defaultValue
}

func (t *ObjectTypeDefinitionFieldList) Add(list ...*ObjectTypeDefinitionField) {
	*t = append(*t, list...)
}

func (t ObjectTypeDefinitionFieldList) Iterator() chan *ObjectTypeDefinitionField {
	ch := make(chan *ObjectTypeDefinitionField, len(t))
	for _, f := range t {
		ch <- f
	}
	close(ch)
	return ch
}

func NewObjectTypeDefinition(name string) *ObjectTypeDefinition {
	return &ObjectTypeDefinition{
		name: name,
	}
}

func (t ObjectTypeDefinition) Name() string {
	return t.name
}

func (t ObjectTypeDefinition) Fields() chan *ObjectTypeDefinitionField {
	return t.fields.Iterator()
}

func (t *ObjectTypeDefinition) AddFields(list ...*ObjectTypeDefinitionField) {
	t.fields.Add(list...)
}

func NewObjectTypeDefinitionField(name string, typ Type) *ObjectTypeDefinitionField {
	return &ObjectTypeDefinitionField{
		name: name,
		typ:  typ,
	}
}

func (t ObjectTypeDefinitionField) Name() string {
	return t.name
}

func (t ObjectTypeDefinitionField) Type() Type {
	return t.typ
}

func (t *ObjectTypeDefinitionField) AddArguments(list ...*ObjectTypeDefinitionFieldArgument) {
	t.arguments.Add(list...)
}

func (t ObjectTypeDefinitionField) Arguments() chan *ObjectTypeDefinitionFieldArgument {
	return t.arguments.Iterator()
}
