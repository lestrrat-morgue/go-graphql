package parser

import (
	"fmt"

	"github.com/lestrrat/go-graphql/model"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	queryKey      = "query"
	mutationKey   = "mutation"
	fragmentKey   = "fragment"
	onKey         = "on"
	typeKey       = "type"
	enumKey       = "enum"
	interfaceKey  = "interface"
	implementsKey = "implements"
)

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func syntaxErr(tok *Token, message string, args ...interface{}) error {
	return errors.Errorf(
		`%s at line %d, column %d`,
		fmt.Sprintf(message, args...),
		tok.Pos.Line,
		tok.Pos.Column,
	)
}

func unexpectedToken(tok *Token, message string, expected ...TokenType) error {
	var value string
	if len(tok.Value) > 0 {
		value = " (" + tok.Value + ")"
	}
	if len(expected) == 0 {
		return syntaxErr(tok, "%s: unexpected token %s%s", message, tok.Type, value)
	}
	return syntaxErr(tok, "%s: expected token %s, but got %s%s", message, expected, tok.Type, value)
}

func unexpectedName(tok *Token, message string, expected ...string) error {
	// XXX tok must be tok.Type == NAME
	if len(expected) == 0 {
		return syntaxErr(tok, "%s: unexpected name %s", message, tok.Value)
	}
	return syntaxErr(tok, "%s: expected name %v, but got %s", message, expected, tok.Value)
}

func consumeToken(pctx *parseCtx, typ TokenType) (*Token, error) {
	switch t := pctx.next(); t.Type {
	case typ:
		return t, nil
	default:
		return nil, syntaxErr(t, `expected token %s, got %s`, typ, t.Type)
	}
}

func consumeName(pctx *parseCtx, names ...string) (string, error) {
	t, err := consumeToken(pctx, NAME)
	if err != nil {
		return "", err
	}

	if len(names) == 0 { // any name is fine
		return t.Value, nil
	}

	for _, name := range names {
		if t.Value == name {
			return t.Value, nil
		}
	}
	return "", syntaxErr(t, `expected name %v, got %s`, names, t.Value)
}

func peekName(pctx *parseCtx, name string) bool {
	switch t := pctx.peek(); t.Type {
	case NAME:
		return t.Value == name
	default:
		return false
	}
}

func (p *Parser) Parse(ctx context.Context, src []byte) (*model.Document, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var pctx parseCtx
	pctx.Context = ctx
	pctx.lexsrc = make(chan *Token, 256)
	pctx.peekCount = -1
	pctx.peekTokens = [3]*Token{}
	pctx.types = make(map[string]*model.NamedType)

	go lex(ctx, src, pctx.lexsrc)

	doc, err := pctx.parseDocument()
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse document`)
	}
	return doc, nil
}

type parseCtx struct {
	context.Context

	lexsrc     chan *Token
	peekCount  int
	peekTokens [3]*Token
	types      map[string]*model.NamedType
}

var eofToken = Token{
	Type: EOF,
}

// peek the next token. this operation fills the peekTokens
// buffer. `next()` is a combination of peek+advance.
//
// note: we do NOT check for peekCout > 2 for efficiency.
// if you do that, you're f*cked.
func (pctx *parseCtx) peek() *Token {
	if pctx.peekCount < 0 {
		select {
		case <-pctx.Context.Done():
			return &eofToken
		case t, ok := <-pctx.lexsrc:
			if !ok {
				return &eofToken
			}
			pctx.peekCount++
			pctx.peekTokens[pctx.peekCount] = t
		}
	}
	return pctx.peekTokens[pctx.peekCount]
}

func (pctx *parseCtx) advance() {
	if pctx.peekCount >= 0 {
		pctx.peekCount--
	}
}

func (pctx *parseCtx) rewind() {
	if pctx.peekCount < 2 {
		pctx.peekCount++
	}
}

func (pctx *parseCtx) next() *Token {
	t := pctx.peek()
	pctx.advance()
	return t
}

func (pctx *parseCtx) registerType(t *model.NamedType) error {
	pctx.types[t.Name()] = t
	return nil
}

func (pctx *parseCtx) lookupType(n string) (*model.NamedType, error) {
	typ, ok := pctx.types[n]
	if !ok {
		return nil, errors.Errorf(`type %s not found`, n)
	}

	return typ, nil
}

func (pctx *parseCtx) parseDocument() (*model.Document, error) {
	doc := model.NewDocument()
	for {
		switch t := pctx.peek(); t.Type {
		case EOF:
			return doc, nil
		case BRACE_L:
			def, err := pctx.parseOperationDefinition(true)
			if err != nil {
				return nil, errors.Wrap(err, `failed to parse operation definition`)
			}
			doc.AddDefinitions(def)
		case NAME:
			switch t.Value {
			case queryKey, mutationKey:
				def, err := pctx.parseOperationDefinition(false)
				if err != nil {
					return nil, errors.Wrap(err, `failed to parse operation definition`)
				}
				doc.AddDefinitions(def)
			case fragmentKey:
				frag, err := pctx.parseFragmentDefinition()
				if err != nil {
					return nil, errors.Wrap(err, `failed to parse fragment definition`)
				}
				doc.AddDefinitions(frag)
			case typeKey:
				typ, err := pctx.parseObjectTypeDefinition()
				if err != nil {
					return nil, errors.Wrap(err, `failed to parse object type definition`)
				}
				doc.AddDefinitions(typ)
			case enumKey:
				enum, err := pctx.parseEnumDefinition()
				if err != nil {
					return nil, errors.Wrap(err, `failed to parse enum definition`)
				}
				doc.AddDefinitions(enum)
			case interfaceKey:
				iface, err := pctx.parseInterfaceDefinition()
				if err != nil {
					return nil, errors.Wrap(err, `failed to parse interface definition`)
				}
				doc.AddDefinitions(iface)
			default:
				return nil, unexpectedName(t, `document`, queryKey, mutationKey, fragmentKey, typeKey, enumKey, interfaceKey)
			}
		default:
			return nil, unexpectedToken(t, `document`)
		}
	}
	return nil, errors.New("error for now")
}

func (pctx *parseCtx) parseTypeCondition() (*model.NamedType, error) {
	switch t := pctx.next(); t.Type {
	case NAME:
		if t.Value != onKey {
			return nil, syntaxErr(t, `expected "on", but got %s`, t.Value)
		}
	default:
		return nil, unexpectedToken(t, `type condition`)
	}

	typ, err := pctx.parseNamedType()
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse named type`)
	}
	return typ.(*model.NamedType), nil
}

func (pctx *parseCtx) parseFragmentName() (string, error) {
	var name string
	switch t := pctx.next(); t.Type {
	case NAME:
		if t.Value == onKey {
			return "", syntaxErr(t, `illegal fragment name "on"`)
		}
		name = t.Value
	}
	return name, nil
}

// FragmentDefinition:
//   fragment FragmentName TypeCondition Directives? SelectionSet
// FragmentName:
//   Name but not on
func (pctx *parseCtx) parseFragmentDefinition() (*model.FragmentDefinition, error) {
	switch t := pctx.next(); t.Type {
	case NAME:
		switch t.Value {
		case fragmentKey:
		default:
			return nil, syntaxErr(t, `expected "fragment", but got %s`, t.Value)
		}
	default:
		return nil, unexpectedToken(t, `fragment`, NAME)
	}

	name, err := pctx.parseFragmentName()
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse fragment name`)
	}

	typ, err := pctx.parseTypeCondition()
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse type condition`)
	}

	fdef := model.NewFragmentDefinition(name, typ)

	switch t := pctx.peek(); t.Type {
	case AT:
		directives, err := pctx.parseDirectives()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse directives`)
		}
		fdef.AddDirectives(directives...)
	}

	set, err := pctx.parseSelectionSet()
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse selection set`)
	}
	fdef.AddSelections(set...)

	return fdef, nil
}

// OperationDefinition:
//   OperationType Name? VariableDefinitions? Directives? SelectionSet
//   SelectionSet
// OperationType: one of
//	 query	mutation
func (pctx *parseCtx) parseOperationDefinition(implicitType bool) (*model.OperationDefinition, error) {
	var optyp model.OperationType
	if implicitType {
		optyp = model.OperationTypeQuery
	} else {
		switch t := pctx.next(); t.Type {
		case NAME:
			switch t.Value {
			case queryKey:
				optyp = model.OperationTypeQuery
			case mutationKey:
				optyp = model.OperationTypeMutation
			default:
				return nil, errors.Errorf(`unknown operation type '%s'`, t.Value)
			}
		}
	}

	def := model.NewOperationDefinition(optyp)

	switch t := pctx.peek(); t.Type {
	case NAME:
		pctx.advance()
		def.SetName(t.Value)
	}

	switch t := pctx.peek(); t.Type {
	case PAREN_L:
		vdef, err := pctx.parseVariableDefinitions()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse query variable definitions`)
		}
		def.AddVariableDefinitions(vdef...)
	}

	switch t := pctx.peek(); t.Type {
	case AT:
		directives, err := pctx.parseDirectives()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse query directives`)
		}
		def.AddDirectives(directives...)
	}

	selections, err := pctx.parseSelectionSet()
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse query selection set`)
	}
	def.AddSelections(selections...)
	return def, nil
}

// VariableDefinitions:
//  ( VariableDefinition... )
func (pctx *parseCtx) parseVariableDefinitions() (model.VariableDefinitionList, error) {
	switch t := pctx.next(); t.Type {
	case PAREN_L:
	default:
		return nil, errors.Errorf(`expected PAREN_L, got %s`, t.Type)
	}

	var list model.VariableDefinitionList
	for loop := true; loop; {
		switch t := pctx.peek(); t.Type {
		case PAREN_R:
			loop = false
			continue
		}

		vdef, err := pctx.parseVariableDefinition()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse variable definition`)
		}
		list = append(list, vdef)
	}

	if pctx.next().Type != PAREN_R {
		return nil, errors.New(`expected PAREN_R`)
	}
	return list, nil
}

// Variable:
//    $  Name
// VariableDefinition:
//    Variable : Type DefaultValue?
// DefaultValue:
//    = Value
func (pctx *parseCtx) parseVariableDefinition() (*model.VariableDefinition, error) {
	switch t := pctx.next(); t.Type {
	case DOLLAR:
	default:
		return nil, unexpectedToken(t, `variable`, DOLLAR)
	}

	var name string
	switch t := pctx.next(); t.Type {
	case NAME:
		name = t.Value
	default:
		return nil, errors.Errorf(`variable: expected NAME, got %s`, t.Type)
	}

	switch t := pctx.next(); t.Type {
	case COLON:
	default:
		return nil, errors.Errorf(`variable: expected COLO, got %s`, t.Type)
	}

	typ, err := pctx.parseType()
	if err != nil {
		return nil, errors.Wrap(err, `variable: failed to parse type`)
	}

	vdef := model.NewVariableDefinition(name, typ)
	if pctx.peek().Type == EQUALS {
		pctx.advance()
		v, err := pctx.parseValue()
		if err != nil {
			return nil, errors.Wrap(err, `variable: failed to parse default value`)
		}
		vdef.SetDefaultValue(v)
	}

	return vdef, nil
}

// Type:
//   NamedType
//   ListType
//   NonNullType
//
// NamedType:
//   Name
//
// ListType:
//   [ Type ]
//
// NonNullType:
//   NamedType !
//   ListType !
func (pctx *parseCtx) parseType() (model.Type, error) {
	var typ model.Type
	var err error
	switch t := pctx.peek(); t.Type {
	case NAME:
		typ, err = pctx.parseNamedType()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse named type`)
		}
	case BRACKET_L:
		typ, err = pctx.parseListType()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse list type`)
		}
	default:
		return nil, errors.Errorf(`expected NamedType or ListType, got %s`, t.Type)
	}

	switch t := pctx.peek(); t.Type {
	case BANG:
		pctx.advance()
		typ.SetNullable(false)
	}
	return typ, nil
}

func (pctx *parseCtx) parseNamedType() (model.Type, error) {
	t := pctx.next()
	if t.Type != NAME {
		return nil, errors.Errorf(`expected Name for NamedType, got %s`, t.Type)
	}

	typ := model.NewNamedType(t.Value)
	if err := pctx.registerType(typ); err != nil {
		return nil, errors.Wrap(err, `failed to register type`)
	}

	return typ, nil
}

func (pctx *parseCtx) parseListType() (model.Type, error) {
	t := pctx.next()
	if t.Type != BRACKET_L {
		return nil, unexpectedToken(t, `list type`, BRACKET_L)
	}

	t = pctx.next()
	if t.Type != NAME {
		return nil, unexpectedToken(t, `list type`, NAME)
	}

	typ, err := pctx.lookupType(t.Value)
	if err != nil {
		typ = model.NewNamedType(t.Value)
		if err := pctx.registerType(typ); err != nil {
			return nil, errors.Wrap(err, `failed to register type`)
		}
	}

	switch t := pctx.next(); t.Type {
	case BRACKET_R:
	default:
		return nil, unexpectedToken(t, `list type`, BRACKET_R)
	}

	return model.NewListType(typ), nil
}

// ValueConst:
//   [~Const] Variable
//   IntValue
//   FloatValue
//   StringValue
//   BooleanValue
//   NullValue
//   EnumValue
//   ListValue [?Const]
//   ObjectValue [?Const]
func (pctx *parseCtx) parseValue() (model.Value, error) {
	switch t := pctx.peek(); t.Type {
	case DOLLAR:
		pctx.advance()
		switch t = pctx.next(); t.Type {
		case NAME:
			return model.NewVariable(t.Value), nil
		default:
			return nil, errors.Errorf(`value: expected NAME, got %s`, t.Type)
		}
	case INT:
		pctx.advance()
		return model.NewIntValue(t.Value)
	case FLOAT:
		pctx.advance()
		return model.NewFloatValue(t.Value)
	case STRING:
		pctx.advance()
		return model.NewStringValue(t.Value), nil
	case BRACE_L:
		return pctx.parseObjectValue()
	case NAME:
		pctx.advance()
		switch t.Value {
		case "true", "false":
			return model.NewBoolValue(t.Value)
		case "null":
			return model.NullValue{}, nil
		default:
			return model.NewEnumValue(t.Value), nil
		}
	default:
		return nil, errors.Errorf(`value: unexpected token %s`, t.Type)
	}
}

func (pctx *parseCtx) parseDirectives() (model.DirectiveList, error) {
	var directives model.DirectiveList
	for loop := true; loop; {
		switch t := pctx.peek(); t.Type {
		case AT:
			pctx.advance()
		default:
			loop = false
			continue
		}

		var name string
		switch t := pctx.next(); t.Type {
		case NAME:
			name = t.Value
		default:
			return nil, unexpectedToken(t, `directive`, NAME)
		}

		d := model.NewDirective(name)
		switch t := pctx.peek(); t.Type {
		case PAREN_L:
			arguments, err := pctx.parseArguments()
			if err != nil {
				return nil, errors.Wrap(err, `failed to parse arguments`)
			}
			d.AddArguments(arguments...)
		}

		directives.Add(d)
	}
	return directives, nil
}

// SelectionSet:
//   { Selection... }
func (pctx *parseCtx) parseSelectionSet() (model.SelectionSet, error) {
	switch t := pctx.next(); t.Type {
	case BRACE_L:
	default:
		return nil, errors.Errorf(`selection set: expected {, got %s`, t.Type)
	}

	var set model.SelectionSet
	for loop := true; loop; {
		switch t := pctx.peek(); t.Type {
		case BRACE_R:
			loop = false
			continue
		}

		sel, err := pctx.parseSelection()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse selection`)
		}
		set = append(set, sel)
	}

	switch t := pctx.next(); t.Type {
	case BRACE_R:
	default:
		return nil, errors.Errorf(`selection set: expected }, got %s`, t.Type)
	}
	return set, nil
}

// Selection:
//   Field
//   FragmentSpread
//   InlineFragment
func (pctx *parseCtx) parseSelection() (model.Selection, error) {
	switch t := pctx.peek(); t.Type {
	case SPREAD:
		pctx.advance()
		return pctx.parseFragmentSpreadOrInlineFragment()
	default:
		return pctx.parseField()
	}
}

func (pctx *parseCtx) parseField() (*model.Field, error) {
	var name string
	var alias string
	switch t := pctx.next(); t.Type {
	case NAME:
		name = t.Value
	default:
		return nil, errors.Errorf(`field: expected NAME, got %s`, t.Type)
	}

	switch t := pctx.peek(); t.Type {
	case COLON:
		pctx.advance()
		alias = name
		switch t = pctx.next(); t.Type {
		case NAME:
			name = t.Value
		default:
			return nil, errors.Errorf(`field: expected NAME, got %s`, t.Type)
		}
	}

	field := model.NewField(name)
	if len(alias) > 0 {
		field.SetAlias(alias)
	}

	switch t := pctx.peek(); t.Type {
	case PAREN_L:
		args, err := pctx.parseArguments()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse arguments`)
		}
		field.AddArguments(args...)
	}

	switch t := pctx.peek(); t.Type {
	case AT:
		directives, err := pctx.parseDirectives()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse directives`)
		}
		field.AddDirectives(directives...)
	}

	switch t := pctx.peek(); t.Type {
	case BRACE_L:
		set, err := pctx.parseSelectionSet()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse selection set`)
		}
		field.AddSelections(set...)
	}
	return field, nil
}

// FragmentSpread:
//   ... FragmentName Directives?
// FragmentName:
//   Name but not "on"
// InlineFragment
//   ... TypeCondition? Directives? SelectionSet
// TypeCondition:
//   on NamedType
func (pctx *parseCtx) parseFragmentSpreadOrInlineFragment() (model.Selection, error) {
	switch t := pctx.peek(); t.Type {
	case BRACE_L, AT:
		return pctx.parseInlineFragment()
	case NAME:
		switch t.Value {
		case "on":
			return pctx.parseInlineFragment()
		}
		// it's something else, then
		return pctx.parseFragmentSpread()
	default:
		return nil, errors.Errorf(`expected FragmentSpread or InlineFragment`)
	}
}

func (pctx *parseCtx) parseArguments() (model.ArgumentList, error) {
	switch t := pctx.next(); t.Type {
	case PAREN_L:
	default:
		return nil, unexpectedToken(t, `arguments`, PAREN_L)
	}

	var args model.ArgumentList

	for loop := true; loop; {
		switch t := pctx.peek(); t.Type {
		case PAREN_R:
			loop = false
			continue
		}

		var name string
		switch t := pctx.next(); t.Type {
		case NAME:
			name = t.Value
		default:
			return nil, unexpectedToken(t, `arguments`, NAME)
		}

		switch t := pctx.next(); t.Type {
		case COLON:
		default:
			return nil, unexpectedToken(t, `arguments`, COLON)
		}

		value, err := pctx.parseValue()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse value`)
		}

		args = append(args, model.NewArgument(name, value))

	}

	switch t := pctx.next(); t.Type {
	case PAREN_R:
	default:
		return nil, unexpectedToken(t, `arguments`, PAREN_R)
	}

	return args, nil
}

// FragmentSpread:
//   ... FragmentName Directives?
func (pctx *parseCtx) parseFragmentSpread() (*model.FragmentSpread, error) {
	// Assumes ... has already been consumed
	name, err := pctx.parseFragmentName()
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse fragment name`)
	}

	frag := model.NewFragmentSpread(name)

	switch t := pctx.peek(); t.Type {
	case AT:
		directives, err := pctx.parseDirectives()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse directives`)
		}
		frag.AddDirectives(directives...)
	}

	return frag, nil
}

// InlineFragment
//   ...TypeCondition? tDirectives? SelectionSet
func (pctx *parseCtx) parseInlineFragment() (*model.InlineFragment, error) {
	// Assumes ... has already been consumed
	var typ *model.NamedType
	switch t := pctx.peek(); t.Type {
	case NAME:
		if t.Value == onKey {
			var err error
			typ, err = pctx.parseTypeCondition()
			if err != nil {
				return nil, errors.Wrap(err, `failed to parse type condition`)
			}
		}
	}

	var directives model.DirectiveList
	switch t := pctx.peek(); t.Type {
	case AT:
		var err error
		directives, err = pctx.parseDirectives()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse directives`)
		}
	}

	selections, err := pctx.parseSelectionSet()
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse selection set`)
	}

	frag := model.NewInlineFragment()
	frag.AddSelections(selections...)
	frag.AddDirectives(directives...)
	frag.SetTypeCondition(typ)

	return frag, nil
}

// ObjectValue:
//   { ObjectField? }
// ObjectField:
//   Name : Value
func (pctx *parseCtx) parseObjectValue() (*model.ObjectValue, error) {
	switch t := pctx.next(); t.Type {
	case BRACE_L:
	default:
		return nil, unexpectedToken(t, `object value`, BRACE_L)
	}

	obj := model.NewObjectValue()
	for loop := true; loop; {
		switch t := pctx.peek(); t.Type {
		case BRACE_R:
			loop = false
			continue
		}

		field, err := pctx.parseObjectField()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse object field`)
		}

		obj.AddFields(field)
	}

	switch t := pctx.next(); t.Type {
	case BRACE_R:
	default:
		return nil, unexpectedToken(t, `object value`, BRACE_R)
	}
	return obj, nil
}

func (pctx *parseCtx) parseObjectField() (*model.ObjectField, error) {
	var name string
	switch t := pctx.next(); t.Type {
	case NAME:
		name = t.Value
	default:
		return nil, unexpectedToken(t, `object field`, NAME)
	}

	switch t := pctx.next(); t.Type {
	case COLON:
	default:
		return nil, unexpectedToken(t, `object field`, COLON)
	}

	v, err := pctx.parseValue()
	if err != nil {
		return nil, errors.Wrap(err, `object field: failed to parse value`)
	}
	return model.NewObjectField(name, v), nil
}

func (pctx *parseCtx) parseObjectTypeDefinition() (*model.ObjectTypeDefinition, error) {
	if _, err := consumeName(pctx, typeKey); err != nil {
		return nil, errors.Wrap(err, `object type`)
	}

	name, err := consumeName(pctx)
	if err != nil {
		return nil, errors.Wrap(err, `object type`)
	}

	var implName string
	if peekName(pctx, implementsKey) {
		if _, err := consumeName(pctx, implementsKey); err != nil {
			return nil, errors.Wrap(err, `object type`)
		}

		implName, err = consumeName(pctx)
		if err != nil {
			return nil, errors.Wrap(err, `object type`)
		}
	}

	if _, err := consumeToken(pctx, BRACE_L); err != nil {
		return nil, errors.Wrap(err, `object type`)
	}

	var fields []*model.ObjectTypeDefinitionField
	for loop := true; loop; {
		switch t := pctx.peek(); t.Type {
		case BRACE_R:
			loop = false
			continue
		}

		field, err := pctx.parseObjectTypeDefinitionField()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse object type field`)
		}
		fields = append(fields, field)
	}

	switch t := pctx.next(); t.Type {
	case BRACE_R:
	default:
		return nil, unexpectedToken(t, `object type`, BRACE_R)
	}

	def := model.NewObjectTypeDefinition(name)
	def.AddFields(fields...)
	if len(implName) > 0 {
		def.SetImplements(implName)
	}
	return def, nil
}

func (pctx *parseCtx) parseObjectTypeDefinitionField() (*model.ObjectTypeDefinitionField, error) {
	var name string
	switch t := pctx.next(); t.Type {
	case NAME:
		name = t.Value
	default:
		return nil, unexpectedToken(t, `object field`, NAME)
	}

	var arguments model.ObjectTypeDefinitionFieldArgumentList
	switch t := pctx.peek(); t.Type {
	case PAREN_L:
		var err error
		arguments, err = pctx.parseObjectTypeDefinitionFieldArguments()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse arguments`)
		}
	}

	switch t := pctx.next(); t.Type {
	case COLON:
	default:
		return nil, unexpectedToken(t, `object field`, COLON)
	}

	typ, err := pctx.parseType()
	if err != nil {
		return nil, errors.Wrap(err, `object field: failed to parse type`)
	}
	f := model.NewObjectTypeDefinitionField(name, typ)
	f.AddArguments(arguments...)
	return f, nil
}

func (pctx *parseCtx) parseObjectTypeDefinitionFieldArguments() (model.ObjectTypeDefinitionFieldArgumentList, error) {
	switch t := pctx.next(); t.Type {
	case PAREN_L:
	default:
		return nil, unexpectedToken(t, `object field arguments`, PAREN_L)
	}

	var args model.ObjectTypeDefinitionFieldArgumentList

	for loop := true; loop; {
		switch t := pctx.peek(); t.Type {
		case PAREN_R:
			loop = false
			continue
		}

		var name string
		switch t := pctx.next(); t.Type {
		case NAME:
			name = t.Value
		default:
			return nil, unexpectedToken(t, `object field arguments`, NAME)
		}

		switch t := pctx.next(); t.Type {
		case COLON:
		default:
			return nil, unexpectedToken(t, `object field arguments`, COLON)
		}

		typ, err := pctx.parseType()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse object field type`)
		}

		arg := model.NewObjectTypeDefinitionFieldArgument(name, typ)

		switch t := pctx.peek(); t.Type {
		case EQUALS:
			// we have default
			pctx.advance()
			value, err := pctx.parseValue()
			if err != nil {
				return nil, errors.Wrap(err, `failed to parse object field default value`)
			}
			arg.SetDefaultValue(value)
		}

		args = append(args, arg)
	}
	switch t := pctx.next(); t.Type {
	case PAREN_R:
	default:
		return nil, unexpectedToken(t, `object field arguments`, PAREN_R)
	}

	return args, nil
}

func (pctx *parseCtx) parseEnumDefinition() (*model.EnumDefinition, error) {
	switch t := pctx.next(); t.Type {
	case NAME:
		if t.Value != enumKey {
			return nil, unexpectedName(t, `enum`, enumKey)
		}
	default:
		return nil, unexpectedToken(t, `enum`, NAME)
	}

	var name string
	switch t := pctx.next(); t.Type {
	case NAME:
		name = t.Value
		// check for valid name
	default:
		return nil, unexpectedToken(t, `enum`, NAME)
	}

	switch t := pctx.next(); t.Type {
	case BRACE_L:
	default:
		return nil, unexpectedToken(t, `enum`, BRACE_L)
	}

	var elements []*model.EnumElement
	for loop := true; loop; {
		switch t := pctx.peek(); t.Type {
		case BRACE_R:
			loop = false
			continue
		case NAME:
			pctx.advance()
			elements = append(elements, model.NewEnumElement(t.Value))
		default:
			return nil, unexpectedToken(t, `enum`, NAME)
		}
	}

	switch t := pctx.next(); t.Type {
	case BRACE_R:
	default:
		return nil, unexpectedToken(t, `enum`, BRACE_R)
	}

	def := model.NewEnumDefinition(name)
	def.AddElements(elements...)
	return def, nil
}

func (pctx *parseCtx) parseInterfaceDefinition() (*model.InterfaceDefinition, error) {
	if _, err := consumeName(pctx, interfaceKey); err != nil {
		return nil, errors.Wrap(err, `interface`)
	}

	name, err := consumeName(pctx)
	if err != nil {
		return nil, errors.Wrap(err, `interface`)
	}

	if _, err := consumeToken(pctx, BRACE_L); err != nil {
		return nil, errors.Wrap(err, `interface`)
	}

	var fields []*model.InterfaceField
	for loop := true; loop; {
		switch t := pctx.peek(); t.Type {
		case BRACE_R:
			loop = false
			continue
		}

		field, err := pctx.parseInterfaceDefinitionField()
		if err != nil {
			return nil, errors.Wrap(err, `interface`)
		}
		fields = append(fields, field)
	}

	if _, err := consumeToken(pctx, BRACE_R); err != nil {
		return nil, errors.Wrap(err, `interface`)
	}
	iface := model.NewInterfaceDefinition(name)
	iface.AddFields(fields...)
	return iface, nil
}

func (pctx *parseCtx) parseInterfaceDefinitionField() (*model.InterfaceField, error) {
	name, err := consumeName(pctx)
	if err != nil {
		return nil, errors.Wrap(err, `interface field`)
	}

	if _, err := consumeToken(pctx, COLON); err != nil {
		return nil, errors.Wrap(err, `interface field`)
	}

	typ, err := pctx.parseType()
	if err != nil {
		return nil, errors.Wrap(err, `interface field`)
	}

	return model.NewInterfaceField(name, typ), nil
}
