package parser

import (
	"unicode"
	"unicode/utf8"
)

const eof = rune(0)

type Position struct {
	Offset int
	Line   int
	Column int
}

type Token struct {
	Type  TokenType
	Value string
	Pos   Position
}

type lrune struct {
	r rune
	w int
}

type Lexer struct {
	input     []byte
	maxpos    int
	peekCount int
	peekRunes [3]lrune
	cur       Position
	start     Position
}

func (l *Lexer) emit(tok *Token, tt TokenType) bool {
	peekOffset := 0
	for i := 0; i <= l.peekCount; i++ {
		peekOffset += l.peekRunes[i].w
	}

	if tt != IGNORABLE {
		tok.Type = tt
		tok.Value = string(l.input[l.start.Offset : l.cur.Offset-peekOffset])
		tok.Pos.Offset = l.start.Offset
		tok.Pos.Line = l.start.Line
		tok.Pos.Column = l.start.Column
	}

	l.start.Offset = l.cur.Offset - peekOffset
	l.start.Line = l.cur.Line
	l.start.Column = l.cur.Column

	if tt == IGNORABLE {
		return false
	}
	return true
}

func (l *Lexer) advance() {
	// if the current rune is a new line, we line++
	switch r := l.peek(); r {
	case '\n':
		l.cur.Line++
		l.cur.Column = 0
	case eof:
	default:
		l.cur.Column++
	}

	l.peekCount--
}

// becareful, we don't check for peekCount
func (l *Lexer) rewind() {
	l.peekCount++
}

func (l *Lexer) next() rune {
	r := l.peek()
	l.advance()
	return r
}

func (l *Lexer) peek() rune {
	if l.peekCount >= 0 {
		return l.peekRunes[l.peekCount].r
	}

	if l.cur.Offset >= l.maxpos {
		return eof
	}

	r, w := utf8.DecodeRune(l.input[l.cur.Offset:])
	l.peekCount++
	l.peekRunes[l.peekCount].r = r
	l.peekRunes[l.peekCount].w = w
	l.cur.Offset += w

	return r
}

func NewLexer(src []byte) *Lexer {
	l := &Lexer{}
	l.input = src
	l.maxpos = len(src)
	l.cur.Offset = 0
	l.cur.Line = 1
	l.cur.Column = 1
	l.start = l.cur
	l.peekCount = -1
	return l
}

func (l *Lexer) Next(tok *Token) bool {
	l.skipInsignificant()

	switch t := l.peek(); t {
	case eof:
		return l.emit(tok, EOF)
	case '!':
		l.advance()
		return l.emit(tok, BANG)
	case '$':
		l.advance()
		return l.emit(tok, DOLLAR)
	case '(':
		l.advance()
		return l.emit(tok, PAREN_L)
	case ')':
		l.advance()
		return l.emit(tok, PAREN_R)
	case ':':
		l.advance()
		return l.emit(tok, COLON)
	case '=':
		l.advance()
		return l.emit(tok, EQUALS)
	case '@':
		l.advance()
		return l.emit(tok, AT)
	case '[':
		l.advance()
		return l.emit(tok, BRACKET_L)
	case ']':
		l.advance()
		return l.emit(tok, BRACKET_R)
	case '{':
		l.advance()
		return l.emit(tok, BRACE_L)
	case '|':
		l.advance()
		return l.emit(tok, PIPE)
	case '}':
		l.advance()
		return l.emit(tok, BRACE_R)
	case '.':
		if !l.runSpread() {
			return l.emit(tok, ILLEGAL)
		}
		return l.emit(tok, SPREAD)
	case '"':
		if !l.runString() {
			return l.emit(tok, ILLEGAL)
		}
		return l.emit(tok, STRING)
	default:
		typ, ok := l.lexValue()
		if !ok {
			return l.emit(tok, ILLEGAL)
		}
		return l.emit(tok, typ)
	}
	return l.emit(tok, ILLEGAL)
}

func (l *Lexer) lexValue() (TokenType, bool) {
	r := l.peek()
	switch {
	case unicode.IsDigit(r):
		return l.lexNumber()
	case r == '-' || r == '+':
		l.advance()
		if unicode.IsDigit(l.peek()) {
			l.rewind()
			return l.lexNumber()
		}
		return ILLEGAL, false
	case r == '"':
		return l.lexString()
	default:
		if !l.runName() {
			return ILLEGAL, false
		}
		return NAME, true
	}
}

func (l *Lexer) runDigits() bool {
	if !unicode.IsDigit(l.next()) {
		return false
	}

	for unicode.IsDigit(l.peek()) {
		l.advance()
	}
	return true
}

func (l *Lexer) lexString() (TokenType, bool) {
	if !l.runString() {
		return ILLEGAL, false
	}

	return STRING, true
}

func (l *Lexer) lexNumber() (TokenType, bool) {
	r := l.next()
	switch r {
	case '-', '+':
		r = l.next()
	}

	var typ TokenType
	if !l.runDigits() {
		return ILLEGAL, false
	}
	typ = INT

	if l.peek() == '.' {
		l.advance()
		if !l.runDigits() {
			return ILLEGAL, false
		}
		typ = FLOAT
	}

	switch l.peek() {
	case 'e', 'E':
		typ = FLOAT
		l.advance()
		switch l.next() {
		case '-', '+':
		default:
			return ILLEGAL, false
		}
		if !l.runDigits() {
			return ILLEGAL, false
		}
	}

	return typ, true
}

func (l *Lexer) skipInsignificant() {
	for {
		switch l.peek() {
		case '\t', ' ', '\n', '\r', ',':
			l.advance()
		default:
			l.emit(nil, IGNORABLE)
			return
		}
	}
}

// ...
func (l *Lexer) runSpread() bool {
	for i := 0; i < 3; i++ {
		if l.next() != '.' {
			return false
		}
	}
	return true
}

// [_A-Za-z][_0-9A-Za-z]*
func (l *Lexer) runName() bool {
	r := l.next()
	switch {
	case r == 0x5f: // _
	case 0x41 <= r && r <= 0x5a: // A-Z
	case 0x61 <= r && r <= 0x7a: // a-z
	default:
		return false
	}

	for {
		r := l.peek()
		switch {
		case r == 0x5f: // _
		case 0x30 <= r && r <= 0x39: // 0-9
		case 0x41 <= r && r <= 0x5a: // A-Z
		case 0x61 <= r && r <= 0x7a: // a-z
		default:
			return true
		}
		l.advance()
	}
	return false
}

func (l *Lexer) runString() bool {
	if l.next() != '"' {
		return false
	}

	for loop := true; loop; {
		switch l.peek() {
		case '"':
			// bail out of the for
			loop = false
			continue
		case '\\':
			if !l.runEscapeSequence() {
				return false
			}
		case '\n', '\r':
			return false
		default:
			l.advance()
		}
	}

	if l.next() != '"' {
		return false
	}
	return true
}

func (l *Lexer) runEscapeSequence() bool {
	if l.next() != '\\' {
		return false
	}

	switch l.peek() {
	case 'u':
		l.advance()
		for i := 0; i < 4; i++ {
			r := l.next()
			switch {
			case 0x30 <= r && r <= 0x39: // 0-9
			case 0x41 <= r && r <= 0x46: // A-F
			case 0x61 <= r && r <= 0x66: // a-f
			default:
				return false
			}
		}
		return true
	case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
		l.advance()
		return true
	}
	return false
}
