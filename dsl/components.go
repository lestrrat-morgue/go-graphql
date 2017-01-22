package dsl

import "github.com/lestrrat/go-graphql/model"

func Schema(attrs ...model.Definition) *model.Schema {
	var query model.ObjectDefinition
	var types model.ObjectDefinitionList
	for _, attr := range attrs {
		if q, ok := attr.(model.ObjectDefinition); ok && q.Name() == `Query` {
			query = q
		} else {
			switch attr.(type) {
			case model.ObjectDefinition:
				types.Add(attr.(model.ObjectDefinition))
			}
		}
	}

	s := model.NewSchema()
	s.SetQuery(query)
	s.AddTypes(types...)
	return s

}

type Attribute interface {
}

type stringAttr string

func (s stringAttr) Value() interface{} {
	return string(s)
}

type nameAttr struct {
	stringAttr
}

func Name(s string) Attribute {
	return nameAttr{stringAttr(s)}
}

func Description(s string) Attribute {
	return stringAttr(s)
}

type ObjectDefinition struct {
	typ model.ObjectDefinition
}

func (v ObjectDefinition) Type() model.ObjectDefinition {
	return v.typ
}

func (v ObjectDefinition) Configure(attrs ...Attribute) ObjectDefinition {
	var fields model.ObjectFieldDefinitionList
	for _, attr := range attrs {
		switch attr.(type) {
		case ObjectBlock:
			attr.(ObjectBlock).Call(v)
		case ImplementsDefinition:
			v.typ.SetImplements(attr.(ImplementsDefinition).typ)
		case model.ObjectFieldDefinition:
			fields.Add(attr.(model.ObjectFieldDefinition))
		}
	}

	v.typ.AddFields(fields...)
	return v
}

type ObjectBlock func(ObjectDefinition)

func (f ObjectBlock) Call(v ObjectDefinition) {
	f(v)
}

type InterfaceDefinition struct {
	typ *model.InterfaceDefinition
}

func (v InterfaceDefinition) Type() *model.InterfaceDefinition {
	return v.typ
}

type ImplementsDefinition struct {
	typ model.NamedType
}

func Implements(t model.NamedType) ImplementsDefinition {
	return ImplementsDefinition{typ: t}
}

func (def ImplementsDefinition) Type() model.Type {
	return def.typ
}

func (v InterfaceDefinition) Configure(attrs ...Attribute) InterfaceDefinition {
	var fields model.InterfaceFieldDefinitionList
	for _, attr := range attrs {
		switch attr.(type) {
		case InterfaceBlock:
			attr.(InterfaceBlock).Call(v)
		case *model.InterfaceFieldDefinition:
			fields.Add(attr.(*model.InterfaceFieldDefinition))
		}
	}

	v.typ.AddFields(fields...)
	return v
}

type InterfaceBlock func(InterfaceDefinition)

func (f InterfaceBlock) Call(v InterfaceDefinition) {
	f(v)
}

type ObjectFieldDefinition struct {
	field model.ObjectFieldDefinition
}
func (def ObjectFieldDefinition) Field() model.ObjectFieldDefinition {
	return def.field
}

func (v *ObjectFieldDefinition) Configure(attrs ...Attribute) *ObjectFieldDefinition {
	var arguments model.ObjectFieldArgumentDefinitionList
	for _, attr := range attrs {
		switch attr.(type) {
		case model.ObjectFieldArgumentDefinition:
			arguments.Add(attr.(model.ObjectFieldArgumentDefinition))
		}
	}
	v.field.AddArguments(arguments...)
	return v
}
