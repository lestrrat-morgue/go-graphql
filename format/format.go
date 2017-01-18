package format

import (
	"bytes"
	"io"
	"strconv"

	"github.com/lestrrat/go-graphql/model"
	"github.com/pkg/errors"
)

type fmtCtx struct {
	dst       io.Writer
	indentbuf []byte
}

func GraphQL(dst io.Writer, v interface{}) error {
	var ctx fmtCtx
	switch v.(type) {
	case *model.Document:
		return ctx.fmtDocument(dst, v.(*model.Document))
	case *model.VariableDefinition:
		return ctx.fmtVariableDefinition(dst, v.(*model.VariableDefinition))
	default:
		return errors.New(`unknown grahql type`)
	}
}

func (ctx *fmtCtx) indent() []byte {
	return ctx.indentbuf
}

func (ctx *fmtCtx) enterleave(f func() error) error {
	ctx.enter()
	defer ctx.leave()
	return f()
}

func (ctx *fmtCtx) enter() {
	ctx.indentbuf = append(ctx.indentbuf, "  "...)
}

func (ctx *fmtCtx) leave() {
	if len(ctx.indentbuf) >= 2 {
		ctx.indentbuf = ctx.indentbuf[:len(ctx.indentbuf)-2]
	}
}

func (ctx *fmtCtx) fmtDocument(dst io.Writer, v *model.Document) error {
	var buf bytes.Buffer
	for def := range v.Definitions() {
		if buf.Len() > 0 {
			buf.WriteString("\n\n")
		}

		switch def.(type) {
		case *model.OperationDefinition:
			if err := ctx.fmtOperationDefinition(&buf, def.(*model.OperationDefinition)); err != nil {
				return errors.Wrap(err, `failed to format operation definition`)
			}
		case *model.FragmentDefinition:
			if err := ctx.fmtFragmentDefinition(&buf, def.(*model.FragmentDefinition)); err != nil {
				return errors.Wrap(err, `failed to format fragment definition`)
			}
		}
	}

	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtFragmentDefinition(dst io.Writer, v *model.FragmentDefinition) error {
	var buf bytes.Buffer

	buf.WriteString("fragment ")
	buf.WriteString(v.Name())
	buf.WriteByte(' ')

	if err := ctx.fmtTypeCondition(&buf, v.Type()); err != nil {
		return errors.Wrap(err, `failed to format type condition`)
	}

	// Directives

	selch := v.SelectionSet()
	if len(selch) > 0 {
		buf.WriteByte(' ')
		if err := ctx.fmtSelectionSet(&buf, selch); err != nil {
			return errors.Wrap(err, `failed to format selection set`)
		}
	}

	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtTypeCondition(dst io.Writer, typ *model.NamedType) error {
	var buf bytes.Buffer
	buf.WriteString("on ")
	buf.WriteString(typ.Name())

	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtVariableDefinitionList(dst io.Writer, vdefch chan *model.VariableDefinition) error {
	l := len(vdefch)
	if l == 0 {
		return nil
	}
	var buf bytes.Buffer
	buf.WriteByte('(')

	i := 0
	for vdef := range vdefch {
		if err := ctx.fmtVariableDefinition(&buf, vdef); err != nil {
			return errors.Wrap(err, `failed to format variable defintiion`)
		}
		if l-1 > i {
			buf.WriteString(", ")
		}
		i++
	}
	buf.WriteByte(')')
	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtOperationDefinition(dst io.Writer, v *model.OperationDefinition) error {
	var buf bytes.Buffer
	buf.WriteString(string(v.Type()))
	if v.HasName() {
		buf.WriteByte(' ')
		buf.WriteString(v.Name())
	}

	if err := ctx.fmtVariableDefinitionList(&buf, v.VariableDefinitions()); err != nil {
		return errors.Wrap(err, `failed to format variable definitions`)
	}

	buf.WriteByte(' ')
	if err := ctx.fmtSelectionSet(&buf, v.SelectionSet()); err != nil {
		return errors.Wrap(err, `failed to format selection set`)
	}

	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtSelectionSet(dst io.Writer, ch chan model.Selection) error {
	if len(ch) == 0 {
		return nil
	}

	var buf bytes.Buffer
	buf.WriteString("{")
	err := ctx.enterleave(func() error {
		for sel := range ch {
			buf.WriteByte('\n')
			buf.Write(ctx.indent())
			if err := ctx.fmtSelection(&buf, sel); err != nil {
				return errors.Wrap(err, `failed to format selection`)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	buf.WriteByte('\n')
	buf.Write(ctx.indent())
	buf.WriteString("}")

	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtSelection(dst io.Writer, v model.Selection) error {
	switch v.(type) {
	case *model.Field:
		return ctx.fmtField(dst, v.(*model.Field))
	case *model.FragmentSpread:
		return ctx.fmtFragmentSpread(dst, v.(*model.FragmentSpread))
	case *model.InlineFragment:
		return ctx.fmtInlineFragment(dst, v.(*model.InlineFragment))
	default:
		return errors.New(`unknown selection`)
	}
}

func (ctx *fmtCtx) fmtInlineFragment(dst io.Writer, v *model.InlineFragment) error {
	var buf bytes.Buffer
	buf.WriteString("... ")

	if typ := v.Type(); typ != nil {
		if err := ctx.fmtTypeCondition(&buf, typ); err != nil {
			return errors.Wrap(err, `failed to format type condition`)
		}
	}

	dirch := v.Directives()
	if len(dirch) > 0 {
		buf.WriteByte(' ')
		if err := ctx.fmtDirectives(&buf, dirch); err != nil {
			return errors.Wrap(err, `failed to format directives`)
		}
	}

	buf.WriteByte(' ')
	if err := ctx.fmtSelectionSet(&buf, v.SelectionSet()); err != nil {
		return errors.Wrap(err, `failed to format selection set`)
	}

	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtArgumentList(dst io.Writer, argch chan *model.Argument) error {
	l := len(argch)
	if l == 0 {
		return nil
	}

	var buf bytes.Buffer
	buf.WriteByte('(')

	argc := 0
	for arg := range argch {
		if err := ctx.fmtArgument(&buf, arg); err != nil {
			return errors.Wrap(err, `failed to format argument`)
		}
		if l-1 > argc {
			buf.WriteString(", ")
		}
		argc++
	}
	buf.WriteByte(')')
	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
	return nil
}

func (ctx *fmtCtx) fmtField(dst io.Writer, v *model.Field) error {
	var buf bytes.Buffer

	if v.HasAlias() {
		buf.WriteString(v.Alias())
		buf.WriteString(": ")
	}
	buf.WriteString(v.Name())

	if err := ctx.fmtArgumentList(&buf, v.Arguments()); err != nil {
		return errors.Wrap(err, `failed to format arguments`)
	}

	dirch := v.Directives()
	if len(dirch) > 0 {
		buf.WriteByte(' ')
		if err := ctx.fmtDirectives(&buf, dirch); err != nil {
			return errors.Wrap(err, `failed to format directives`)
		}
	}

	selch := v.SelectionSet()
	if len(selch) > 0 {
		buf.WriteByte(' ')
		if err := ctx.fmtSelectionSet(&buf, selch); err != nil {
			return errors.Wrap(err, `failed to format selection set`)
		}
	}

	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtFragmentSpread(dst io.Writer, v *model.FragmentSpread) error {
	var buf bytes.Buffer
	buf.WriteString("...")
	buf.WriteString(v.Name())
	// Directives...

	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtArgument(dst io.Writer, v *model.Argument) error {
	var buf bytes.Buffer
	buf.WriteString(v.Name())
	buf.WriteString(": ")
	if err := ctx.fmtValue(&buf, v.Value()); err != nil {
		return errors.Wrap(err, `failed to format value`)
	}

	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtVariableDefinition(dst io.Writer, v *model.VariableDefinition) error {
	var buf bytes.Buffer
	buf.WriteByte('$')
	buf.WriteString(v.Name())
	buf.WriteString(": ")
	if err := ctx.fmtType(&buf, v.Type()); err != nil {
		return errors.Wrap(err, `failed to format type`)
	}

	if v.HasDefaultValue() {
		buf.WriteString(" = ")
		if err := ctx.fmtValue(&buf, v.DefaultValue()); err != nil {
			return errors.Wrap(err, `failed to format default value`)
		}
	}

	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtValue(dst io.Writer, v model.Value) error {
	var buf bytes.Buffer

	switch v.(type) {
	case *model.Variable, model.Variable:
		buf.WriteByte('$')
		buf.WriteString(v.Value().(string))
	case *model.IntValue, model.IntValue:
		buf.WriteString(strconv.Itoa(v.Value().(int)))
	case *model.FloatValue, model.FloatValue:
		buf.WriteString(strconv.FormatFloat(v.Value().(float64), 'g', -1, 64))
	case *model.StringValue, model.StringValue:
		buf.WriteString(v.Value().(string))
	case *model.BoolValue, model.BoolValue:
		buf.WriteString(strconv.FormatBool(v.Value().(bool)))
	case *model.NullValue, model.NullValue:
		buf.WriteString("null")
	case *model.EnumValue, model.EnumValue:
		buf.WriteString(v.Value().(string))
	case *model.ObjectValue:
		buf.WriteByte('{')
		err := ctx.enterleave(func() error {
			for field := range v.(*model.ObjectValue).Fields() {
				buf.WriteByte('\n')
				buf.Write(ctx.indent())
				buf.WriteString(field.Name())
				buf.WriteString(": ")
				if err := ctx.fmtValue(&buf,field.Value()); err != nil {
					return errors.Wrap(err, `failed to format object field value`)
				}
			}
			return nil
		})
		if err != nil {
			return errors.Wrap(err, `failed to format object`)
		}
		buf.WriteByte('\n')
		buf.Write(ctx.indent())
		buf.WriteByte('}')
	default:
		return errors.New(`unsupported value`)
	}
	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtType(dst io.Writer, v model.Type) error {
	switch v.(type) {
	case *model.NamedType:
		return ctx.fmtNamedType(dst, v.(*model.NamedType))
	case *model.ListType:
		return ctx.fmtListType(dst, v.(*model.ListType))
	default:
		return errors.Errorf(`invalid type %s`, v)
	}
}

func (ctx *fmtCtx) fmtNamedType(dst io.Writer, v *model.NamedType) error {
	var buf bytes.Buffer
	buf.WriteString(v.Name())
	if !v.IsNullable() {
		buf.WriteByte('!')
	}
	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtListType(dst io.Writer, v *model.ListType) error {
	var buf bytes.Buffer
	buf.WriteString("[ ")
	if err := ctx.fmtType(dst, v.Type()); err != nil {
		return errors.Wrap(err, `failed to format type`)
	}

	buf.WriteString(" ]")
	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}

func (ctx *fmtCtx) fmtDirectives(dst io.Writer, dirch chan *model.Directive) error {
	l := len(dirch)
	if l == 0 {
		return nil
	}

	var buf bytes.Buffer
	i := 0
	for dir := range dirch {
		buf.WriteByte('@')
		buf.WriteString(dir.Name())
		if err := ctx.fmtArgumentList(&buf, dir.Arguments()); err != nil {
			return errors.Wrap(err, `failed to format arguments`)
		}

		if l-1 > i {
			buf.WriteByte(' ')
		}
		i++
	}

	if _, err := buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
}
