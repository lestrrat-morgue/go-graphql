package model

func NewSchema() *Schema {
	return &Schema{}
}

func (s Schema) Query() *ObjectDefinition {
	return s.query
}

func (s Schema) Types() chan *ObjectDefinition {
	return s.types.Iterator()
}

func (s *Schema) SetQuery(q *ObjectDefinition) {
	s.query = q
}

func (s *Schema) AddTypes(list ...*ObjectDefinition) {
	s.types.Add(list...)
}

func NewNamedType(name string) NamedType {
	return &namedType{
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

func (t *NamedTypeList) Add(list ...NamedType) {
	*t = append(*t, list...)
}

func (t NamedTypeList) Iterator() chan NamedType {
	ch := make(chan NamedType, len(t))
	for _, f := range t {
		ch <- f
	}
	close(ch)
	return ch
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

func (l *ObjectDefinitionList) Add(list ...*ObjectDefinition) {
	*l  = append(*l, list...)
}

func (l ObjectDefinitionList) Iterator() chan *ObjectDefinition {
	ch := make(chan *ObjectDefinition, len(l))
	for _, o := range l {
		ch <- o
	}
	close(ch)
	return ch
}

func NewObjectDefinition(name string) *ObjectDefinition {
	return &ObjectDefinition{
		nameComponent: nameComponent(name),
		nullable: nullable(true),
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

func (t *ObjectDefinition) SetImplements(typ NamedType) {
	t.hasImplements = true
	t.implements = typ
}

func (t ObjectDefinition) HasImplements() bool {
	return t.hasImplements
}

func (t ObjectDefinition) Implements() NamedType {
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
		nullable:      nullable(true),
	}
}

func (t *EnumDefinition) AddElements(list ...*EnumElementDefinition) {
	t.elements.Add(list...)
}

func (t *EnumDefinition) Elements() chan *EnumElementDefinition {
	return t.elements.Iterator()
}

func NewEnumElementDefinition(name string, value Value) *EnumElementDefinition {
	return &EnumElementDefinition{
		nameComponent:  nameComponent(name),
		valueComponent: valueComponent{value: value},
	}
}

func (e *EnumElementDefinitionList) Add(list ...*EnumElementDefinition) {
	*e = append(*e, list...)
}

func (e EnumElementDefinitionList) Iterator() chan *EnumElementDefinition {
	ch := make(chan *EnumElementDefinition, len(e))
	for _, el := range e {
		ch <- el
	}
	close(ch)
	return ch
}

func NewInterfaceDefinition(name string) *InterfaceDefinition {
	return &InterfaceDefinition{
		nullable:      nullable(true),
		nameComponent: nameComponent(name),
	}
}

type Resolver interface {
	Resolve(interface{}) Type
}

func (iface *InterfaceDefinition) SetTypeResolver(v Resolver) {}
func (iface *InterfaceDefinition) TypeResolver() Resolver     { return nil }

func (iface InterfaceDefinition) Fields() chan *InterfaceFieldDefinition {
	return iface.fields.Iterator()
}

func (iface *InterfaceDefinition) AddFields(list ...*InterfaceFieldDefinition) {
	iface.fields.Add(list...)
}

func NewInterfaceFieldDefinition(name string, typ Type) *InterfaceFieldDefinition {
	return &InterfaceFieldDefinition{
		nameComponent: nameComponent(name),
		typeComponent: typeComponent{typ: typ},
	}
}

func (f *InterfaceFieldDefinition) Type() Type {
	return f.typ
}

func (f *InterfaceFieldDefinitionList) Add(list ...*InterfaceFieldDefinition) {
	*f = append(*f, list...)
}

func (f InterfaceFieldDefinitionList) Iterator() chan *InterfaceFieldDefinition {
	ch := make(chan *InterfaceFieldDefinition, len(f))
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
