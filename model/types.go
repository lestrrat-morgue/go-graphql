package model

func (n nullable) IsNullable() bool {
	return bool(n)
}

func (n *nullable) SetNullable(b bool) {
	*n = nullable(b)
}

func NewNamedType(name string) *NamedType {
	return &NamedType{
		name: name,
		nullable: true,
	}
}

func (t *NamedType) Name() string {
	return t.name
}

func NewListType(t Type) *ListType {
	return &ListType{
		nullable: true,
		typ: t,
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

func (t *ObjectTypeFieldList) Add(list ...*ObjectTypeField) {
	*t = append(*t, list...)
}

func (t ObjectTypeFieldList) Iterator() chan *ObjectTypeField {
	ch := make(chan *ObjectTypeField, len(t))
	for _, f := range t {
		ch <- f
	}
	close(ch)
	return ch
}

func NewObjectTypeDefinition(name string) *ObjectTypeDefinition {
	return &ObjectTypeDefinition {
		name: name,
	}
}

func (t ObjectTypeDefinition) Name() string {
	return t.name
}

func (t ObjectTypeDefinition) Fields() chan *ObjectTypeField {
	return t.fields.Iterator()
}

func (t *ObjectTypeDefinition) AddFields(list ...*ObjectTypeField) {
	t.fields.Add(list...)
}

func NewObjectTypeField(name string, typ Type) *ObjectTypeField {
	return &ObjectTypeField{
		name: name,
		typ: typ,
	}
}

func (t ObjectTypeField) Name() string {
	return t.name
}

func (t ObjectTypeField) Type() Type {
	return t.typ
}

func (t *ObjectTypeField) AddArguments(list ...*Argument) {
	t.arguments.Add(list...)
}

func (t ObjectTypeField) Arguments() chan *Argument {
	return t.arguments.Iterator()
}
