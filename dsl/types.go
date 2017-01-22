package dsl

import "github.com/lestrrat/go-graphql/model"

func NotNull(v model.Type) model.Type {
	if n, ok := v.(model.Nullable); ok {
		n.SetNullable(false)
	}
	return v
}

func String() model.Type {
	// This is wrong. just a place holder
	return model.NewNamedType(`String`)
}

func List(t model.Type) model.ListType {
	return model.NewListType(t)
}

func Interface(name string, attrs ...Attribute) InterfaceDefinition {
	var def InterfaceDefinition
	def.typ = model.NewInterfaceDefinition(name)
	return def.Configure(attrs...)
}

func InterfaceField(name string, typ model.Type, attrs ...Attribute) model.InterfaceFieldDefinition {
	return model.NewInterfaceFieldDefinition(name, typ)
}

func Object(name string, attrs ...Attribute) ObjectDefinition {
	var def ObjectDefinition
	def.typ = model.NewObjectDefinition(name)
	return def.Configure(attrs...)
}

func ObjectField(name string, typ model.Type, attrs ...Attribute) model.ObjectFieldDefinition {
	var def ObjectFieldDefinition
	def.field = model.NewObjectFieldDefinition(name, typ)
	return def.Configure(attrs...).Field()
}

func ObjectFieldArgument(name string, typ model.Type, attrs ...Attribute) model.ObjectFieldArgumentDefinition {
	return model.NewObjectFieldArgumentDefinition(name, typ)
}
