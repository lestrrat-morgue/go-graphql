package model

func NewNamedType(name string) *NamedType {
	return &NamedType{
		nameComponent: nameComponent(name),
		nullable:      true,
	}
}

func NewListType(t Type) *ListType {
	return &ListType{
		nullable:      true,
		typeComponent: typeComponent{typ: t},
	}
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

func (args *ObjectFieldArgumentDefinitionList) Add(list ...*ObjectFieldArgumentDefinition) {
	*args = append(*args, list...)
}

func (args ObjectFieldArgumentDefinitionList) Iterator() chan *ObjectFieldArgumentDefinition {
	ch := make(chan *ObjectFieldArgumentDefinition, len(args))
	for _, arg := range args {
		ch <- arg
	}
	close(ch)
	return ch
}

func NewObjectFieldArgumentDefinition(name string, typ Type) *ObjectFieldArgumentDefinition {
	return &ObjectFieldArgumentDefinition{
		nameComponent: nameComponent(name),
		typeComponent: typeComponent{typ: typ},
	}
}

func (t *ObjectFieldDefinitionList) Add(list ...*ObjectFieldDefinition) {
	*t = append(*t, list...)
}

func (t ObjectFieldDefinitionList) Iterator() chan *ObjectFieldDefinition {
	ch := make(chan *ObjectFieldDefinition, len(t))
	for _, f := range t {
		ch <- f
	}
	close(ch)
	return ch
}

func NewObjectDefinition(name string) *ObjectDefinition {
	return &ObjectDefinition{
		nameComponent: nameComponent(name),
	}
}

func (t ObjectDefinition) Fields() chan *ObjectFieldDefinition {
	return t.fields.Iterator()
}

func (t *ObjectDefinition) AddFields(list ...*ObjectFieldDefinition) {
	t.fields.Add(list...)
}

func NewObjectFieldDefinition(name string, typ Type) *ObjectFieldDefinition {
	return &ObjectFieldDefinition{
		nameComponent: nameComponent(name),
		typeComponent: typeComponent{typ: typ},
	}
}

func (t *ObjectDefinition) SetImplements(n string) {
	t.hasImplements = true
	t.implements = n
}

func (t ObjectDefinition) HasImplements() bool {
	return t.hasImplements
}

func (t ObjectDefinition) Implements() string {
	return t.implements
}

func (t *ObjectFieldDefinition) AddArguments(list ...*ObjectFieldArgumentDefinition) {
	t.arguments.Add(list...)
}

func (t ObjectFieldDefinition) Arguments() chan *ObjectFieldArgumentDefinition {
	return t.arguments.Iterator()
}

func NewEnumDefinition(name string) *EnumDefinition {
	return &EnumDefinition{
		nameComponent: nameComponent(name),
	}
}

func (t *EnumDefinition) AddElements(list ...*EnumElement) {
	t.elements.Add(list...)
}

func (t *EnumDefinition) Elements() chan *EnumElement {
	return t.elements.Iterator()
}

func NewEnumElement(name string) *EnumElement {
	return &EnumElement{
		nameComponent: nameComponent(name),
	}
}

func (e *EnumElementList) Add(list ...*EnumElement) {
	*e = append(*e, list...)
}

func (e EnumElementList) Iterator() chan *EnumElement {
	ch := make(chan *EnumElement, len(e))
	for _, el := range e {
		ch <- el
	}
	close(ch)
	return ch
}

func NewInterfaceDefinition(name string) *InterfaceDefinition {
	return &InterfaceDefinition{
		nameComponent: nameComponent(name),
	}
}

func (iface InterfaceDefinition) Fields() chan *InterfaceField {
	return iface.fields.Iterator()
}

func (iface *InterfaceDefinition) AddFields(list ...*InterfaceField) {
	iface.fields.Add(list...)
}

func NewInterfaceField(name string, typ Type) *InterfaceField {
	return &InterfaceField{
		nameComponent: nameComponent(name),
		typeComponent: typeComponent{typ: typ},
	}
}

func (f *InterfaceField) Type() Type {
	return f.typ
}

func (f *InterfaceFieldList) Add(list ...*InterfaceField) {
	*f = append(*f, list...)
}

func (f InterfaceFieldList) Iterator() chan *InterfaceField {
	ch := make(chan *InterfaceField, len(f))
	for _, field := range f {
		ch <- field
	}
	close(ch)
	return ch
}

func NewUnionDefinition(name string) *UnionDefinition {
	return &UnionDefinition{
		nameComponent: nameComponent(name),
	}
}

func (def UnionDefinition) Types() chan Type {
	return def.types.Iterator()
}
func (def *UnionDefinition) AddTypes(list ...Type) {
	def.types.Add(list...)
}

func NewInputDefinition(name string) *InputDefinition {
	return &InputDefinition{
		nameComponent: nameComponent(name),
	}
}

func NewInputFieldDefinition(name string) *InputFieldDefinition {
	return &InputFieldDefinition{
		nameComponent: nameComponent(name),
	}
}

func (def *InputDefinition) AddFields(list ...*InputFieldDefinition) {
	def.fields.Add(list...)
}

func (def InputDefinition) Fields() chan *InputFieldDefinition {
	return def.fields.Iterator()
}

func (l *InputFieldDefinitionList) Add(list ...*InputFieldDefinition) {
	*l = append(*l, list...)
}

func (l InputFieldDefinitionList) Iterator() chan *InputFieldDefinition {
	ch := make(chan *InputFieldDefinition, len(l))
	for _, e := range l {
		ch <- e
	}
	close(ch)
	return ch
}


