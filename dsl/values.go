package dsl

import "github.com/lestrrat/go-graphql/model"

func IntValue(v int) model.Value {
	return model.NewIntValue(v)
}
