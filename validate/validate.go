package validate

import (
	"errors"

	"github.com/lestrrat/go-graphql/model"
	"github.com/lestrrat/go-graphql/visitor"
	"golang.org/x/net/context"
)

var h = &visitor.Handler{
	EnterFragmentSpread:     enterFragmentSpread,
	EnterFragmentDefinition: enterFragmentDefinition,
	LeaveFragmentDefinition: leaveFragmentDefinition,
}

type validationCtx struct {
	context.Context

	container interface{} // Current top level container
	schema    model.Document
}

func Validate(c context.Context, schema, doc model.Document) error {
	var ctx validationCtx
	ctx.Context = c
	ctx.schema = schema
	return visitor.Visit(&ctx, h, doc)
}

func enterFragmentDefinition(c context.Context, v model.FragmentDefinition) error {
	ctx := c.(*validationCtx)
	ctx.container = v
	return nil
}

func leaveFragmentDefinition(c context.Context, v model.FragmentDefinition) error {
	ctx := c.(*validationCtx)
	if x, ok := ctx.container.(model.FragmentDefinition); ok {
		if x.Name() == v.Name() {
			ctx.container = nil
		}
	}
	return nil
}

func enterFragmentSpread(c context.Context, v model.FragmentSpread) error {
	ctx := c.(*validationCtx)
	if ctx.container == nil {
		return errors.New(`fragment spread in top level`)
	}

	f, ok := ctx.container.(model.FragmentDefinition)
	if ok {
		if f.Name() == v.Name() {
			return errors.New(`can not spread fragment inside the same named fragment`)
		}
	}
	return nil
}
