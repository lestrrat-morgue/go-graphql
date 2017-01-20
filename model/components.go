package model

type nullable bool

func (n nullable) IsNullable() bool {
	return bool(n)
}

func (n *nullable) SetNullable(b bool) {
	*n = nullable(b)
}

// nameComponent allows us to hide the name and to also
// provide a default Name() method for every component that requires
// a string name
type nameComponent string

func (n nameComponent) Name() string {
	return string(n)
}

// typeComponent allows us to hide the type and to also
// provide a default Type() as well as SetType() methods for
// every component that requires a Type
type typeComponent struct {
	typ Type
}

func (t typeComponent) Type() Type {
	return Type(t.typ)
}

func (t *typeComponent) SetType(newt Type) {
	t.typ = newt
}

// defaultValueComponent allows us to hide the defaultValue and to also
// provide a default HasDefaultValue(), DefaultValue(), SetDefaultValue(), and
// RemoveDefaultValue() for every component that requires a deault value
type defaultValueComponent struct {
	valid bool
	value Value
}

func (v defaultValueComponent) HasDefaultValue() bool {
	return v.valid
}

func (v defaultValueComponent) DefaultValue() Value {
	return v.value
}

func (v *defaultValueComponent) SetDefaultValue(newv Value) {
	v.valid = true
	v.value = newv
}

func (v *defaultValueComponent) RemoveDefaultValue() {
	v.valid = false
	v.value = nil
}
