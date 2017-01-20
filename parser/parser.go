package parser

import (
	"fmt"

	"github.com/lestrrat/go-graphql/model"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	enumKey       = "enum"
	falseKey      = "false"
	fragmentKey   = "fragment"
	implementsKey = "implements"
	inputKey      = "input"
	interfaceKey  = "interface"
	mutationKey   = "mutation"
	nullKey       = "null"
	onKey         = "on"
	queryKey      = "query"
	trueKey       = "true"
	typeKey       = "type"
	unionKey      = "union"
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

func peekToken(pctx *parseCtx, typ TokenType) bool {
	switch t := pctx.peek(); t.Type {
	case typ:
		return true
	default:
		return false
	}
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
	pctx.lexsrc = NewLexer(src)
	pctx.peekCount = -1
	pctx.peekTokens = [3]Token{}
	pctx.types = make(map[string]*model.NamedType)

	doc, err := pctx.parseDocument()
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse document`)
	}
	return doc, nil
}

type parseCtx struct {
	context.Context

	lexsrc     *Lexer
	peekCount  int
	peekTokens [3]Token
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
		default:
		}
		if !pctx.lexsrc.Next(&pctx.peekTokens[pctx.peekCount+1]) {
			return &eofToken
		}
		pctx.peekCount++
	}
	return &pctx.peekTokens[pctx.peekCount]
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
				typ, err := pctx.parseObjectDefinition()
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
			case unionKey:
				union, err := pctx.parseUnionDefinition()
				if err != nil {
					return nil, errors.Wrap(err, `failed to parse union definition`)
				}
				doc.AddDefinitions(union)
			case inputKey:
				input, err := pctx.parseInputDefinition()
				if err != nil {
					return nil, errors.Wrap(err, `failed to parse input definition`)
				}
				doc.AddDefinitions(input)
			default:
				return nil, unexpectedName(t, `document`, queryKey, mutationKey, fragmentKey, typeKey, enumKey, interfaceKey, unionKey)
			}
		default:
			return nil, unexpectedToken(t, `document`)
		}
	}
	return nil, errors.New("error for now")
}

func (pctx *parseCtx) parseTypeCondition() (*model.NamedType, error) {
	if _, err := consumeName(pctx, onKey); err != nil {
		return nil, errors.Wrap(err, `type condition`)
	}

	typ, err := pctx.parseNamedType()
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse named type`)
	}
	return typ.(*model.NamedType), nil
}

func (pctx *parseCtx) parseFragmentName() (string, error) {
	if peekName(pctx, onKey) {
		return "", errors.New(`fragment name: illegal fragment name "on"`)
	}
	return consumeName(pctx)
}

// FragmentDefinition:
//   fragment FragmentName TypeCondition Directives? SelectionSet
// FragmentName:
//   Name but not on
func (pctx *parseCtx) parseFragmentDefinition() (*model.FragmentDefinition, error) {
	t, err := consumeToken(pctx, NAME)
	if err != nil {
		return nil, errors.Wrap(err, `fragment definition`)
	}
	switch t.Value {
	case fragmentKey:
	default:
		return nil, syntaxErr(t, `expected "fragment", but got %s`, t.Value)
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
	if peekToken(pctx, AT) {
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
		name, err := consumeName(pctx)
		if err != nil {
			return nil, errors.Wrap(err, `operation definition`)
		}
		switch name {
		case queryKey:
			optyp = model.OperationTypeQuery
		case mutationKey:
			optyp = model.OperationTypeMutation
		default:
			return nil, errors.Errorf(`unknown operation type '%s'`, name)
		}
	}

	def := model.NewOperationDefinition(optyp)
	if t := pctx.peek(); t.Type == NAME {
		pctx.advance()
		def.SetName(t.Value)
	}

	if peekToken(pctx, PAREN_L) {
		vdef, err := pctx.parseVariableDefinitions()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse query variable definitions`)
		}
		def.AddVariableDefinitions(vdef...)
	}

	if peekToken(pctx, AT) {
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
	if _, err := consumeToken(pctx, PAREN_L); err != nil {
		return nil, errors.Wrap(err, `variable`)
	}

	var list model.VariableDefinitionList
	for loop := true; loop; {
		if peekToken(pctx, PAREN_R) {
			loop = false
			continue
		}

		vdef, err := pctx.parseVariableDefinition()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse variable definition`)
		}
		list = append(list, vdef)
	}

	if _, err := consumeToken(pctx, PAREN_R); err != nil {
		return nil, errors.Wrap(err, `variable`)
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
	if _, err := consumeToken(pctx, DOLLAR); err != nil {
		return nil, errors.Wrap(err, `variable`)
	}

	name, err := consumeName(pctx)
	if err != nil {
		return nil, errors.Wrap(err, `variable`)
	}

	if _, err := consumeToken(pctx, COLON); err != nil {
		return nil, errors.Wrap(err, `variable`)
	}

	typ, err := pctx.parseType()
	if err != nil {
		return nil, errors.Wrap(err, `variable: failed to parse type`)
	}

	vdef := model.NewVariableDefinition(name, typ)
	if peekToken(pctx, EQUALS) {
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

	if peekToken(pctx, BANG) {
		pctx.advance()
		typ.SetNullable(false)
	}
	return typ, nil
}

func (pctx *parseCtx) parseNamedType() (model.Type, error) {
	typname, err := consumeName(pctx)
	if err != nil {
		return nil, errors.Wrap(err, `named type`)
	}

	typ := model.NewNamedType(typname)
	if err := pctx.registerType(typ); err != nil {
		return nil, errors.Wrap(err, `failed to register type`)
	}

	return typ, nil
}

func (pctx *parseCtx) parseListType() (model.Type, error) {
	if _, err := consumeToken(pctx, BRACKET_L); err != nil {
		return nil, errors.Wrap(err, `list type`)
	}

	typname, err := consumeName(pctx)
	if err != nil {
		return nil, errors.Wrap(err, `list type`)
	}

	typ, err := pctx.lookupType(typname)
	if err != nil {
		typ = model.NewNamedType(typname)
		if err := pctx.registerType(typ); err != nil {
			return nil, errors.Wrap(err, `failed to register type`)
		}
	}

	if _, err := consumeToken(pctx, BRACKET_R); err != nil {
		return nil, errors.Wrap(err, `list type`)
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
		name, err := consumeName(pctx)
		if err != nil {
			return nil, errors.Wrap(err, `value`)
		}
		return model.NewVariable(name), nil
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
		name := t.Value

		switch name {
		case trueKey, falseKey:
			return model.NewBoolValue(name)
		case nullKey:
			return model.NullValue{}, nil
		default:
			return model.NewEnumValue(name), nil
		}
	default:
		return nil, errors.Errorf(`value: unexpected token %s`, t.Type)
	}
}

func (pctx *parseCtx) parseDirectives() (model.DirectiveList, error) {
	var directives model.DirectiveList
	for loop := true; loop; {
		if !peekToken(pctx, AT) {
			loop = false
			continue
		}
		pctx.advance()

		name, err := consumeName(pctx)
		if err != nil {
			return nil, errors.Wrap(err, `directive`)
		}

		d := model.NewDirective(name)
		if peekToken(pctx, PAREN_L) {
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
	if _, err := consumeToken(pctx, BRACE_L); err != nil {
		return nil, errors.Wrap(err, `selection set`)
	}

	var set model.SelectionSet
	for loop := true; loop; {
		if peekToken(pctx, BRACE_R) {
			loop = false
			continue
		}

		sel, err := pctx.parseSelection()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse selection`)
		}
		set = append(set, sel)
	}

	if _, err := consumeToken(pctx, BRACE_R); err != nil {
		return nil, errors.Wrap(err, `selection set`)
	}
	return set, nil
}

// Selection:
//   Field
//   FragmentSpread
//   InlineFragment
func (pctx *parseCtx) parseSelection() (model.Selection, error) {
	if peekToken(pctx, SPREAD) {
		pctx.advance()
		return pctx.parseFragmentSpreadOrInlineFragment()
	}
	return pctx.parseField()
}

func (pctx *parseCtx) parseField() (*model.Field, error) {
	var name string
	var alias string
	name, err := consumeName(pctx)
	if err != nil {
		return nil, errors.Wrap(err, `field`)
	}

	if peekToken(pctx, COLON) {
		pctx.advance()
		alias = name
		name, err = consumeName(pctx)
		if err != nil {
			return nil, errors.Wrap(err, `field`)
		}
	}

	field := model.NewField(name)
	if len(alias) > 0 {
		field.SetAlias(alias)
	}

	if peekToken(pctx, PAREN_L) {
		args, err := pctx.parseArguments()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse arguments`)
		}
		field.AddArguments(args...)
	}

	if peekToken(pctx, AT) {
		directives, err := pctx.parseDirectives()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse directives`)
		}
		field.AddDirectives(directives...)
	}

	if peekToken(pctx, BRACE_L) {
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
		case onKey:
			return pctx.parseInlineFragment()
		}
		// it's something else, then
		return pctx.parseFragmentSpread()
	default:
		return nil, errors.Errorf(`expected FragmentSpread or InlineFragment`)
	}
}

func (pctx *parseCtx) parseArguments() (model.ArgumentList, error) {
	if _, err := consumeToken(pctx, PAREN_L); err != nil {
		return nil, errors.Wrap(err, `arguments`)
	}

	var args model.ArgumentList

	for loop := true; loop; {
		if peekToken(pctx, PAREN_R) {
			loop = false
			continue
		}

		name, err := consumeName(pctx)
		if err != nil {
			return nil, errors.Wrap(err, `arguments`)
		}

		if _, err := consumeToken(pctx, COLON); err != nil {
			return nil, errors.Wrap(err, `arguments`)
		}

		value, err := pctx.parseValue()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse value`)
		}

		args = append(args, model.NewArgument(name, value))

	}

	if _, err := consumeToken(pctx, PAREN_R); err != nil {
		return nil, errors.Wrap(err, `arguments`)
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

	if peekToken(pctx, AT) {
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
	if peekName(pctx, onKey) {
		var err error
		typ, err = pctx.parseTypeCondition()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse type condition`)
		}
	}

	var directives model.DirectiveList
	if peekToken(pctx, AT) {
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
	if _, err := consumeToken(pctx, BRACE_L); err != nil {
		return nil, errors.Wrap(err, `object value`)
	}

	obj := model.NewObjectValue()
	for loop := true; loop; {
		if peekToken(pctx, BRACE_R) {
			loop = false
			continue
		}

		field, err := pctx.parseObjectField()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse object field`)
		}

		obj.AddFields(field)
	}

	if _, err := consumeToken(pctx, BRACE_R); err != nil {
		return nil, errors.Wrap(err, `object value`)
	}
	return obj, nil
}

func (pctx *parseCtx) parseObjectField() (*model.ObjectField, error) {
	name, err := consumeName(pctx)
	if err != nil {
		return nil, errors.Wrap(err, `object field`)
	}

	if _, err := consumeToken(pctx, COLON); err != nil {
		return nil, errors.Wrap(err, `object field`)
	}

	v, err := pctx.parseValue()
	if err != nil {
		return nil, errors.Wrap(err, `object field: failed to parse value`)
	}
	return model.NewObjectField(name, v), nil
}

func (pctx *parseCtx) parseObjectDefinition() (*model.ObjectDefinition, error) {
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

	var fields []*model.ObjectFieldDefinition
	for loop := true; loop; {
		if peekToken(pctx, BRACE_R) {
			loop = false
			continue
		}

		field, err := pctx.parseObjectFieldDefinition()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse object type field`)
		}
		fields = append(fields, field)
	}

	if _, err := consumeToken(pctx, BRACE_R); err != nil {
		return nil, errors.Wrap(err, `object type`)
	}

	def := model.NewObjectDefinition(name)
	def.AddFields(fields...)
	if len(implName) > 0 {
		def.SetImplements(implName)
	}
	return def, nil
}

func (pctx *parseCtx) parseObjectFieldDefinition() (*model.ObjectFieldDefinition, error) {
	name, err := consumeName(pctx)
	if err != nil {
		return nil, errors.Wrap(err, `object field`)
	}

	var arguments model.ObjectFieldArgumentDefinitionList
	if peekToken(pctx, PAREN_L) {
		var err error
		arguments, err = pctx.parseObjectFieldArgumentDefinitions()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse arguments`)
		}
	}

	if _, err := consumeToken(pctx, COLON); err != nil {
		return nil, errors.Wrap(err, `object field`)
	}

	typ, err := pctx.parseType()
	if err != nil {
		return nil, errors.Wrap(err, `object field: failed to parse type`)
	}
	f := model.NewObjectFieldDefinition(name, typ)
	f.AddArguments(arguments...)
	return f, nil
}

func (pctx *parseCtx) parseObjectFieldArgumentDefinitions() (model.ObjectFieldArgumentDefinitionList, error) {
	if _, err := consumeToken(pctx, PAREN_L); err != nil {
		return nil, errors.Wrap(err, `object field arguments`)
	}

	var args model.ObjectFieldArgumentDefinitionList

	for loop := true; loop; {
		if peekToken(pctx, PAREN_R) {
			loop = false
			continue
		}

		name, err := consumeName(pctx)
		if err != nil {
			return nil, errors.Wrap(err, `object field arguments`)
		}

		if _, err := consumeToken(pctx, COLON); err != nil {
			return nil, errors.Wrap(err, `object field arguments`)
		}

		typ, err := pctx.parseType()
		if err != nil {
			return nil, errors.Wrap(err, `failed to parse object field type`)
		}

		arg := model.NewObjectFieldArgumentDefinition(name, typ)

		if peekToken(pctx, EQUALS) {
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

	if _, err := consumeToken(pctx, PAREN_R); err != nil {
		return nil, errors.Wrap(err, `object field arguments`)
	}

	return args, nil
}

func (pctx *parseCtx) parseEnumDefinition() (*model.EnumDefinition, error) {
	if _, err := consumeName(pctx, enumKey); err != nil {
		return nil, errors.Wrap(err, `enum`)
	}

	name, err := consumeName(pctx)
	if err != nil {
		return nil, errors.Wrap(err, `enum`)
	}

	if _, err := consumeToken(pctx, BRACE_L); err != nil {
		return nil, errors.Wrap(err, `enum`)
	}

	var elements []*model.EnumElement
	for loop := true; loop; {
		if peekToken(pctx, BRACE_R) {
			loop = false
			continue
		}

		elem, err := consumeName(pctx)
		if err != nil {
			return nil, errors.Wrap(err, `enum`)
		}
		elements = append(elements, model.NewEnumElement(elem))
	}

	if _, err := consumeToken(pctx, BRACE_R); err != nil {
		return nil, errors.Wrap(err, `enum`)
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
		if peekToken(pctx, BRACE_R) {
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

func (pctx *parseCtx) parseUnionDefinition() (*model.UnionDefinition, error) {
	if _, err := consumeName(pctx, unionKey); err != nil {
		return nil, errors.Wrap(err, `union`)
	}

	name, err := consumeName(pctx)
	if err != nil {
		return nil, errors.Wrap(err, `union`)
	}

	if _, err := consumeToken(pctx, EQUALS); err != nil {
		return nil, errors.Wrap(err, `union`)
	}

	union := model.NewUnionDefinition(name)

	typ, err := pctx.parseType()
	if err != nil {
		return nil, errors.Wrap(err, `union`)
	}

	var types []model.Type
	types = append(types, typ)

	for loop := true; loop; {
		if !peekToken(pctx, PIPE) {
			loop = false
			continue
		}
		pctx.advance()

		typ, err := pctx.parseType()
		if err != nil {
			return nil, errors.Wrap(err, `union`)
		}
		types = append(types, typ)
	}
	union.AddTypes(types...)

	return union, nil
}

func (pctx *parseCtx) parseInputDefinition() (*model.InputDefinition, error) {
	if _, err := consumeName(pctx, inputKey); err != nil {
		return nil, errors.Wrap(err, `input`)
	}

	name, err := consumeName(pctx)
	if err != nil {
		return nil, errors.Wrap(err, `input`)
	}

	if _, err := consumeToken(pctx, BRACE_L); err != nil {
		return nil, errors.Wrap(err, `input`)
	}

	var fields []*model.InputFieldDefinition
	for loop := true; loop; {
		if peekToken(pctx, BRACE_R) {
			loop = false
			continue
		}

		field, err := pctx.parseInputDefinitionField()
		if err != nil {
			return nil, errors.Wrap(err, `input`)
		}
		fields = append(fields, field)
	}

	if _, err := consumeToken(pctx, BRACE_R); err != nil {
		return nil, errors.Wrap(err, `input`)
	}
	iface := model.NewInputDefinition(name)
	iface.AddFields(fields...)
	return iface, nil
}

func (pctx *parseCtx) parseInputDefinitionField() (*model.InputFieldDefinition, error) {
	name, err := consumeName(pctx)
	if err != nil {
		return nil, errors.Wrap(err, `input field`)
	}

	if _, err := consumeToken(pctx, COLON); err != nil {
		return nil, errors.Wrap(err, `input field`)
	}

	typ, err := pctx.parseType()
	if err != nil {
		return nil, errors.Wrap(err, `input field`)
	}

	def := model.NewInputFieldDefinition(name)
	def.SetType(typ)
	return def, nil
}
