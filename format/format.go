package format

import (
	"bytes"
	"io"
	"strconv"

	"github.com/lestrrat/go-graphql/model"
	"github.com/lestrrat/go-graphql/visitor"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type fmtCtx struct {
	context.Context
	buf       *bytes.Buffer
	indentbuf []byte

	elements         []int
	element          int
	padSelectionList bool
}

var fmtHandler = &visitor.Handler{
	EnterDefinition:               enterDefinition,
	EnterDefinitionList:           enterList,
	LeaveDefinitionList:           leaveList,
	EnterOperationDefinition:      enterOperationDefinition,
	EnterInputDefinition:          enterInputDefinition,
	LeaveInputDefinition:          leaveInputDefinition,
	EnterInputFieldDefinition:     enterInputFieldDefinition,
	EnterSelectionList:            enterSelectionList,
	EnterSelection:                enterSelection,
	LeaveSelectionList:            leaveSelectionList,
	EnterSelectionField:           enterSelectionField,
	EnterInlineFragment:           enterInlineFragment,
	EnterFragmentSpread:           enterFragmentSpread,
	EnterFragmentDefinition:       enterFragmentDefinition,
	EnterDirective:                enterDirective,
	EnterDirectiveList:            enterDirectiveList,
	EnterUnionDefinition:          enterUnionDefinition,
	EnterInterfaceDefinition:      enterInterfaceDefinition,
	LeaveInterfaceDefinition:      leaveInterfaceDefinition,
	EnterInterfaceFieldDefinition: enterInterfaceFieldDefinition,
	EnterObjectDefinition:         enterObjectDefinition,
	LeaveObjectDefinition:         leaveObjectDefinition,
	EnterObjectFieldDefinition:    enterObjectFieldDefinition,
	EnterEnumDefinition:           enterEnumDefinition,
	EnterSchema:                   enterSchema,
}

const singleindent = "  "

func enterList(c context.Context) error {
	// this pushes a new counter stack to keep track of list elements
	ctx := c.(*fmtCtx)
	ctx.elements = append(ctx.elements, 0)
	ctx.element++
	return nil
}

func moreIndent(c context.Context) {
	ctx := c.(*fmtCtx)
	ctx.indentbuf = append(ctx.indentbuf, singleindent...)
}

func leaveList(c context.Context) error {
	// pops the counter stack
	ctx := c.(*fmtCtx)
	ctx.elements = ctx.elements[:ctx.element]
	ctx.element--
	return nil
}

func lessIndent(c context.Context) {
	ctx := c.(*fmtCtx)
	ctx.indentbuf = ctx.indentbuf[:len(ctx.indentbuf)-len(singleindent)]
}

func enterDefinition(c context.Context, v model.Definition) error {
	ctx := c.(*fmtCtx)
	if ctx.elements[ctx.element] > 0 {
		ctx.buf.WriteString("\n\n")
	}
	ctx.elements[ctx.element]++
	return nil
}

func enterOperationDefinition(c context.Context, v model.OperationDefinition) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf

	buf.WriteString(string(v.OperationType()))
	if v.HasName() {
		buf.WriteByte(' ')
		buf.WriteString(v.Name())
	}

	ch := v.Variables()
	if l := len(ch); l > 0 {
		buf.WriteByte('(')
		i := 0
		for vardef := range ch {
			buf.WriteByte('$')
			buf.WriteString(vardef.Name())
			buf.WriteString(": ")
			if err := fmtType(ctx, vardef.Type()); err != nil {
				return errors.Wrap(err, `failed to format type`)
			}

			if vardef.HasDefaultValue() {
				buf.WriteString(" = ")
				if err := fmtValue(ctx, vardef.DefaultValue()); err != nil {
					return errors.Wrap(err, `failed to format default value`)
				}
			}
			if l-1 > i {
				buf.WriteString(", ")
			}
			i++
		}
		buf.WriteByte(')')
	}
	ctx.padSelectionList = true
	return nil
}

func enterSelectionList(c context.Context) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf
	if ctx.padSelectionList {
		buf.WriteByte(' ')
		ctx.padSelectionList = false
	}
	buf.WriteString("{")

	moreIndent(c)
	return enterList(c)
}

func leaveSelectionList(c context.Context) error {
	if err := leaveList(c); err != nil {
		return err
	}
	lessIndent(c)

	ctx := c.(*fmtCtx)
	buf := ctx.buf
	buf.WriteByte('\n')
	buf.Write(ctx.indentbuf)
	buf.WriteByte('}')
	return nil
}

func enterSelection(c context.Context, v model.Selection) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf
	buf.WriteString("\n")
	for i := 0; i < ctx.element; i++ {
		buf.WriteString("  ")
	}
	return nil
}

func enterSelectionField(c context.Context, v model.SelectionField) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf

	if v.HasAlias() {
		buf.WriteString(v.Alias())
		buf.WriteString(": ")
	}
	buf.WriteString(v.Name())

	if err := fmtArgumentList(ctx, v.Arguments()); err != nil {
		return errors.Wrap(err, `failed to format arguments`)
	}

	ctx.padSelectionList = true
	return nil
}

func enterInlineFragment(c context.Context, v model.InlineFragment) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf
	buf.WriteString("... ")

	if typ := v.TypeCondition(); typ != nil {
		if err := fmtTypeCondition(ctx, typ); err != nil {
			return errors.Wrap(err, `failed to format type condition`)
		}
	}
	ctx.padSelectionList = true
	return nil
}

func enterFragmentSpread(c context.Context, v model.FragmentSpread) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf
	buf.WriteString("...")
	buf.WriteString(v.Name())
	return nil
}

func enterFragmentDefinition(c context.Context, v model.FragmentDefinition) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf
	buf.WriteString("fragment ")
	buf.WriteString(v.Name())
	buf.WriteByte(' ')

	if err := fmtTypeCondition(ctx, v.Type().(model.NamedType)); err != nil {
		return errors.Wrap(err, `failed to format type condition`)
	}
	ctx.padSelectionList = true
	return nil
}

func enterDirectiveList(c context.Context) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf
	buf.WriteByte(' ')
	return nil
}

func enterSchema(c context.Context, v model.Schema) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf

	buf.WriteString("schema {")
	moreIndent(c)

	buf.WriteByte('\n')
	buf.Write(ctx.indentbuf)
	buf.WriteString("query: ")
	buf.WriteString(v.Query().Name())

	if ch := v.Types(); len(ch) > 0 {
		buf.WriteByte('\n')
		buf.Write(ctx.indentbuf)
		buf.WriteString("types: [")
		l := len(ch)
		i := 0
		for typ := range ch {
			buf.WriteString(typ.Name())
			if l-1 > i {
				buf.WriteString(", ")
			}
			i++
		}
		buf.WriteByte(']')
	}
	buf.WriteByte('\n')
	buf.WriteByte('}')
	return nil
}

func enterDirective(c context.Context, v model.Directive) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf
	if ctx.elements[ctx.element] > 0 {
		ctx.buf.WriteByte('0')
	}

	buf.WriteByte('@')
	buf.WriteString(v.Name())
	if err := fmtArgumentList(ctx, v.Arguments()); err != nil {
		return errors.Wrap(err, `failed to format arguments`)
	}
	return nil
}

func enterInputDefinition(c context.Context, v model.InputDefinition) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf
	buf.WriteString("input ")
	buf.WriteString(v.Name())
	buf.WriteString(" {")
	moreIndent(c)
	return enterList(c)
}

func leaveInputDefinition(c context.Context, v model.InputDefinition) error {
	if err := leaveList(c); err != nil {
		return err
	}
	lessIndent(c)

	ctx := c.(*fmtCtx)
	buf := ctx.buf
	buf.WriteByte('\n')
	buf.Write(ctx.indentbuf)
	buf.WriteByte('}')
	return nil
}

func enterInputFieldDefinition(c context.Context, v model.InputFieldDefinition) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf

	buf.WriteByte('\n')
	buf.Write(ctx.indentbuf)
	buf.WriteString(v.Name())
	buf.WriteString(": ")
	if err := fmtType(ctx, v.Type()); err != nil {
		return errors.Wrap(err, `failed to format field type`)
	}
	return nil
}

func enterUnionDefinition(c context.Context, v model.UnionDefinition) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf

	buf.WriteString("union ")
	buf.WriteString(v.Name())
	buf.WriteString(" = ")

	ch := v.Types()
	if len(ch) == 0 {
		return errors.New(`union without any types to compose is meaningless`)
	}

	// write the first one
	t := <-ch
	if err := fmtType(ctx, t); err != nil {
		return errors.Wrap(err, `failed to format type`)
	}

	// write the rest, if any. these will be preceded by
	// a '|' (pipe)
	for len(ch) > 0 {
		t = <-ch
		buf.WriteString(" | ")
		if err := fmtType(ctx, t); err != nil {
			return errors.Wrap(err, `failed to format type`)
		}
	}

	return nil
}

func enterInterfaceDefinition(c context.Context, v model.InterfaceDefinition) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf

	buf.WriteString("interface ")
	buf.WriteString(v.Name())
	buf.WriteString(" {")
	moreIndent(c)
	return nil
}

func leaveInterfaceDefinition(c context.Context, v model.InterfaceDefinition) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf

	lessIndent(c)

	buf.WriteByte('\n')
	buf.Write(ctx.indentbuf)
	buf.WriteByte('}')
	return nil
}

func enterInterfaceFieldDefinition(c context.Context, v model.InterfaceFieldDefinition) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf

	buf.WriteByte('\n')
	buf.Write(ctx.indentbuf)
	buf.WriteString(v.Name())
	buf.WriteString(": ")
	if err := fmtType(ctx, v.Type()); err != nil {
		return errors.Wrap(err, `failed to format field type`)
	}
	return nil
}

func enterObjectDefinition(c context.Context, v model.ObjectDefinition) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf

	buf.WriteString("type ")
	buf.WriteString(v.Name())
	if v.HasImplements() {
		buf.WriteString(" implements ")
		buf.WriteString(v.Implements().(model.NamedType).Name())
	}
	buf.WriteString(" {")
	moreIndent(c)
	return nil
}

func leaveObjectDefinition(c context.Context, v model.ObjectDefinition) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf

	lessIndent(c)

	buf.WriteByte('\n')
	buf.Write(ctx.indentbuf)
	buf.WriteByte('}')
	return nil
}

func enterObjectFieldDefinition(c context.Context, v model.ObjectFieldDefinition) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf

	buf.WriteByte('\n')
	buf.Write(ctx.indentbuf)
	buf.WriteString(v.Name())
	if err := fmtObjectFieldArgumentDefinitionList(ctx, v.Arguments()); err != nil {
		return errors.Wrap(err, `failed to format object field argumets`)
	}
	buf.WriteString(": ")
	if err := fmtType(ctx, v.Type()); err != nil {
		return errors.Wrap(err, `failed to format object field type`)
	}
	return nil
}

func enterEnumDefinition(c context.Context, v model.EnumDefinition) error {
	ctx := c.(*fmtCtx)
	buf := ctx.buf

	buf.WriteString("enum ")
	buf.WriteString(v.Name())
	buf.WriteString(" {")
	ch := v.Elements()
	if len(ch) == 0 {
		return errors.New(`enum without any elements to compose is meaningless`)
	}

	moreIndent(c)
	for e := range ch {
		buf.WriteByte('\n')
		buf.Write(ctx.indentbuf)
		buf.WriteString(e.Name())
	}
	lessIndent(c)
	buf.WriteString("\n}")
	return nil
}

func GraphQL(c context.Context, dst io.Writer, v interface{}) error {
	var b = make([]byte, 0, 4096)
	var ctx fmtCtx
	ctx.Context = c
	ctx.buf = bytes.NewBuffer(b)
	ctx.element = -1

	if err := visitor.Visit(&ctx, fmtHandler, v); err != nil {
		return err
	}
	if _, err := ctx.buf.WriteTo(dst); err != nil {
		return errors.Wrap(err, `failed to write to destination`)
	}
	return nil
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

func fmtTypeCondition(ctx *fmtCtx, typ model.NamedType) error {
	buf := ctx.buf
	buf.WriteString("on ")
	buf.WriteString(typ.Name())
	return nil
}

func fmtArgumentList(ctx *fmtCtx, argch chan model.Argument) error {
	l := len(argch)
	if l == 0 {
		return nil
	}

	buf := ctx.buf
	buf.WriteByte('(')

	argc := 0
	for arg := range argch {
		if err := fmtArgument(ctx, arg); err != nil {
			return errors.Wrap(err, `failed to format argument`)
		}
		if l-1 > argc {
			buf.WriteString(", ")
		}
		argc++
	}
	buf.WriteByte(')')
	return nil
}

func fmtArgument(ctx *fmtCtx, v model.Argument) error {
	buf := ctx.buf

	buf.WriteString(v.Name())
	buf.WriteString(": ")
	if err := fmtValue(ctx, v.Value()); err != nil {
		return errors.Wrap(err, `failed to format value`)
	}
	return nil
}

func fmtValue(ctx *fmtCtx, v model.Value) error {
	buf := ctx.buf

	switch v.Kind() {
	case model.VariableKind:
		buf.WriteByte('$')
		buf.WriteString(v.Value().(string))
	case model.IntKind:
		buf.WriteString(strconv.Itoa(v.Value().(int)))
	case model.FloatKind:
		buf.WriteString(strconv.FormatFloat(v.Value().(float64), 'g', -1, 64))
	case model.StringKind, model.EnumKind:
		buf.WriteString(v.Value().(string))
	case model.BooleanKind:
		buf.WriteString(strconv.FormatBool(v.Value().(bool)))
	case model.NullKind:
		buf.WriteString("null")
	case model.ObjectKind:
		buf.WriteByte('{')
		moreIndent(ctx)

		for field := range v.(model.ObjectValue).Fields() {
			buf.WriteByte('\n')
			buf.Write(ctx.indentbuf)
			buf.WriteString(field.Name())
			buf.WriteString(": ")
			if err := fmtValue(ctx, field.Value()); err != nil {
				return errors.Wrap(err, `failed to format object field value`)
			}
		}
		lessIndent(ctx)

		buf.WriteByte('\n')
		buf.Write(ctx.indent())
		buf.WriteByte('}')
	default:
		return errors.New(`unsupported value`)
	}
	return nil
}

func fmtType(ctx *fmtCtx, v model.Type) error {
	switch v.(type) {
	case model.NamedType:
		return fmtNamedType(ctx, v.(model.NamedType))
	case model.ListType:
		return fmtListType(ctx, v.(model.ListType))
	default:
		return errors.Errorf(`invalid type %s`, v)
	}
}

func fmtNamedType(ctx *fmtCtx, v model.NamedType) error {
	buf := ctx.buf

	buf.WriteString(v.Name())
	if !v.IsNullable() {
		buf.WriteByte('!')
	}
	return nil
}

func fmtListType(ctx *fmtCtx, v model.ListType) error {
	buf := ctx.buf
	buf.WriteByte('[')
	if err := fmtType(ctx, v.Type()); err != nil {
		return errors.Wrap(err, `failed to format type`)
	}

	buf.WriteByte(']')
	if !v.IsNullable() {
		buf.WriteByte('!')
	}
	return nil
}

func fmtObjectFieldArgumentDefinition(ctx *fmtCtx, v model.ObjectFieldArgumentDefinition) error {
	buf := ctx.buf

	buf.WriteString(v.Name())
	buf.WriteString(": ")
	if err := fmtType(ctx, v.Type()); err != nil {
		return errors.Wrap(err, `failed to format type`)
	}

	if v.HasDefaultValue() {
		buf.WriteString(" = ")
		if err := fmtValue(ctx, v.DefaultValue()); err != nil {
			return errors.Wrap(err, `failed to format default value`)
		}
	}
	return nil
}

func fmtObjectFieldArgumentDefinitionList(ctx *fmtCtx, argch chan model.ObjectFieldArgumentDefinition) error {
	l := len(argch)
	if l == 0 {
		return nil
	}

	buf := ctx.buf
	buf.WriteByte('(')

	argc := 0
	for arg := range argch {
		if err := fmtObjectFieldArgumentDefinition(ctx, arg); err != nil {
			return errors.Wrap(err, `failed to format argument`)
		}
		if l-1 > argc {
			buf.WriteString(", ")
		}
		argc++
	}
	buf.WriteByte(')')

	return nil
}
