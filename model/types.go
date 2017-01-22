package model

func NewSchema() Schema {
	return &schema{}
}

func (s schema) Query() ObjectDefinition {
	return s.query
}

func (s schema) Types() chan ObjectDefinition {
	return s.types.Iterator()
}

func (s *schema) SetQuery(q ObjectDefinition) {
	s.query = q
}

func (s *schema) AddTypes(list ...ObjectDefinition) {
	s.types.Add(list...)
}

func NewNamedType(name string) NamedType {
	return &namedType{
		nameComponent: nameComponent(name),
		nullable:      true,
	}
}

func NewListType(t Type) ListType {
	return &listType{
		nullable:      true,
		typeComponent: typeComponent{typ: t},
	}
}

func NewObjectFieldArgumentDefinition(name string, typ Type) ObjectFieldArgumentDefinition {
	return &objectFieldArgumentDefinition{
		nameComponent: nameComponent(name),
		typeComponent: typeComponent{typ: typ},
	}
}

func NewObjectDefinition(name string) ObjectDefinition {
	return &objectDefinition{
		nameComponent: nameComponent(name),
		nullable: nullable(true),
	}
}

func (t objectDefinition) Fields() chan ObjectFieldDefinition {
	return t.fields.Iterator()
}

func (t *objectDefinition) AddFields(list ...ObjectFieldDefinition) {
	t.fields.Add(list...)
}

func NewObjectFieldDefinition(name string, typ Type) ObjectFieldDefinition {
	return &objectFieldDefinition{
		nameComponent: nameComponent(name),
		typeComponent: typeComponent{typ: typ},
	}
}

func (t *objectDefinition) SetImplements(typ NamedType) {
	t.hasImplements = true
	t.implements = typ
}

func (t objectDefinition) HasImplements() bool {
	return t.hasImplements
}

func (t objectDefinition) Implements() NamedType {
	return t.implements
}

func (t *objectFieldDefinition) AddArguments(list ...ObjectFieldArgumentDefinition) {
	t.arguments.Add(list...)
}

func (t objectFieldDefinition) Arguments() chan ObjectFieldArgumentDefinition {
	return t.arguments.Iterator()
}

func NewEnumDefinition(name string) EnumDefinition {
	return &enumDefinition{
		nameComponent: nameComponent(name),
		nullable:      nullable(true),
	}
}

func (t *enumDefinition) AddElements(list ...EnumElementDefinition) {
	t.elements.Add(list...)
}

func (t *enumDefinition) Elements() chan EnumElementDefinition {
	return t.elements.Iterator()
}

func NewEnumElementDefinition(name string, value Value) EnumElementDefinition {
	return &enumElementDefinition{
		nameComponent:  nameComponent(name),
		valueComponent: valueComponent{value: value},
	}
}

func NewInterfaceDefinition(name string) InterfaceDefinition {
	return &interfaceDefinition{
		nullable:      nullable(true),
		nameComponent: nameComponent(name),
	}
}

type Resolver interface {
	Resolve(interface{}) Type
}

func (iface *interfaceDefinition) SetTypeResolver(v Resolver) {}
func (iface *interfaceDefinition) TypeResolver() Resolver     { return nil }

func (iface interfaceDefinition) Fields() chan InterfaceFieldDefinition {
	return iface.fields.Iterator()
}

func (iface *interfaceDefinition) AddFields(list ...InterfaceFieldDefinition) {
	iface.fields.Add(list...)
}

func NewInterfaceFieldDefinition(name string, typ Type) InterfaceFieldDefinition {
	return &interfaceFieldDefinition{
		nameComponent: nameComponent(name),
		typeComponent: typeComponent{typ: typ},
	}
}

func (f *interfaceFieldDefinition) Type() Type {
	return f.typ
}

func NewUnionDefinition(name string) UnionDefinition {
	return &unionDefinition{
		nameComponent: nameComponent(name),
	}
}

func (def unionDefinition) Types() chan Type {
	return def.types.Iterator()
}
func (def *unionDefinition) AddTypes(list ...Type) {
	def.types.Add(list...)
}

func NewInputDefinition(name string) InputDefinition {
	return &inputDefinition{
		nameComponent: nameComponent(name),
	}
}

func NewInputFieldDefinition(name string) InputFieldDefinition {
	return &inputFieldDefinition{
		nameComponent: nameComponent(name),
	}
}

func (def *inputDefinition) AddFields(list ...InputFieldDefinition) {
	def.fields.Add(list...)
}

func (def inputDefinition) Fields() chan InputFieldDefinition {
	return def.fields.Iterator()
}
