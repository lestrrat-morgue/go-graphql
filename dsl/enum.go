package dsl

import "github.com/lestrrat/go-graphql/model"

func EnumValue(name string, value model.Value, attrs ...Attribute) model.EnumElementDefinition {
	return model.NewEnumElementDefinition(name, value)
}

func Enum(attrs ...Attribute) model.EnumDefinition {
	var name string
	var elements model.EnumElementDefinitionList
	for _, attr := range attrs {
		switch attr.(type) {
		case nameAttr:
			name = attr.(nameAttr).Value().(string)
		case model.EnumElementDefinition:
			elements.Add(attr.(model.EnumElementDefinition))
		}
	}

	def := model.NewEnumDefinition(name)
	def.AddElements(elements...)
	return def
}
