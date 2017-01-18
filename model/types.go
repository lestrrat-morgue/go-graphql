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
