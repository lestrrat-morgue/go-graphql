package parser

import (
	"unicode"
	"unicode/utf8"

	"golang.org/x/net/context"
)

const eof = rune(0)

type Position struct {
	Line int
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

type lexCtx struct {
	context.Context
	col int
	input     []byte
	line int
	maxpos    int
	out       chan *Token
	peekCount int
	peekRunes [3]lrune
	pos       int
	start     int
}

func (lctx *lexCtx) emit(tt TokenType) {
	offset := 0
	for i := 0; i <= lctx.peekCount; i++ {
		offset += lctx.peekRunes[i].w
	}

	if tt != IGNORABLE {
		var tok Token
		tok.Type = tt
		tok.Value = string(lctx.input[lctx.start : lctx.pos-offset])
		tok.Pos.Line = lctx.line
		tok.Pos.Column = lctx.col

		select {
		case <-lctx.Done():
		case lctx.out <- &tok:
		}
	}
	lctx.start = lctx.pos - offset
}

func (lctx *lexCtx) advance() {
  // if the current rune is a new line, we line++
  switch r := lctx.peek(); r {
  case '\n':
    lctx.line++
    lctx.col = 0
  case eof:
  default:
    lctx.col++
  }

	lctx.peekCount--
}

// becareful, we don't check for peekCount
func (l *lexCtx) rewind() {
	l.peekCount++
}

func (l *lexCtx) next() rune {
	r := l.peek()
	l.advance()
	return r
}

func (l *lexCtx) peek() rune {
	if l.peekCount >= 0 {
		return l.peekRunes[l.peekCount].r
	}

	if l.pos >= l.maxpos {
		return eof
	}

	r, w := utf8.DecodeRune(l.input[l.pos:])
	l.peekCount++
	l.peekRunes[l.peekCount].r = r
	l.peekRunes[l.peekCount].w = w
	l.pos += w

	return r
}

func lex(ctx context.Context, src []byte, ch chan *Token) {
	defer close(ch)

	var lctx lexCtx
	lctx.Context = ctx
	lctx.out = ch
	lctx.input = src
	lctx.line = 1
	lctx.maxpos = len(src)
	lctx.pos = 0
	lctx.peekCount = -1

	for {
		select {
		case <-lctx.Done():
			return
		default:
		}

		lctx.skipInsignificant()

		switch t := lctx.peek(); t {
		case eof:
			lctx.emit(EOF)
			return
		case '!':
			lctx.advance()
			lctx.emit(BANG)
		case '$':
			lctx.advance()
			lctx.emit(DOLLAR)
		case '(':
			lctx.advance()
			lctx.emit(PAREN_L)
		case ')':
			lctx.advance()
			lctx.emit(PAREN_R)
		case ':':
			lctx.advance()
			lctx.emit(COLON)
		case '=':
			lctx.advance()
			lctx.emit(EQUALS)
		case '@':
			lctx.advance()
			lctx.emit(AT)
		case '[':
			lctx.advance()
			lctx.emit(BRACKET_L)
		case ']':
			lctx.advance()
			lctx.emit(BRACKET_R)
		case '{':
			lctx.advance()
			lctx.emit(BRACE_L)
		case '|':
			lctx.advance()
			lctx.emit(PIPE)
		case '}':
			lctx.advance()
			lctx.emit(BRACE_R)
		case '.':
			if !lctx.runSpread() {
				lctx.emit(ILLEGAL)
				return
			}
			lctx.emit(SPREAD)
		case '"':
			if !lctx.runString() {
				lctx.emit(ILLEGAL)
				return
			}
			lctx.emit(STRING)
		default:
			if !lctx.lexValue() {
				lctx.emit(ILLEGAL)
			}
		}
	}
}

func (lctx *lexCtx) lexValue() bool {
	r := lctx.peek()
	switch {
	case unicode.IsDigit(r):
		return lctx.lexNumber()
	case r == '-' || r == '+':
		lctx.advance()
		if unicode.IsDigit(lctx.peek()) {
			lctx.rewind()
			return lctx.lexNumber()
		}
		return false
	case r == '"':
		return lctx.lexString()
	default:
		if !lctx.runName() {
			return false
		}
		lctx.emit(NAME)
		return true
	}
}

func (lctx *lexCtx) runDigits() bool {
	if !unicode.IsDigit(lctx.next()) {
		return false
	}

	for unicode.IsDigit(lctx.peek()) {
		lctx.advance()
	}
	return true
}

func (lctx *lexCtx) lexString() bool {
	return false
}

func (lctx *lexCtx) lexNumber() bool {
	r := lctx.next()
	switch r {
	case '-', '+':
		r = lctx.next()
	}

	var typ TokenType
	if !lctx.runDigits() {
		return false
	}
	typ = INT

	if lctx.peek() == '.' {
		lctx.advance()
		if !lctx.runDigits() {
			return false
		}
		typ = FLOAT
	}

	switch lctx.peek() {
	case 'e', 'E':
		typ = FLOAT
		lctx.advance()
		switch lctx.next() {
		case '-', '+':
		default:
			return false
		}
		if !lctx.runDigits() {
			return false
		}
	}

	lctx.emit(typ)
	return true
}

func (lctx *lexCtx) skipInsignificant() {
	for {
		switch lctx.peek() {
		case '\t', ' ', '\n', '\r', ',':
			lctx.advance()
		default:
			lctx.emit(IGNORABLE)
			return
		}
	}
}

// ...
func (lctx *lexCtx) runSpread() bool {
	for i := 0; i < 3; i++ {
		if lctx.next() != '.' {
			return false
		}
	}
	return true
}

// [_A-Za-z][_0-9A-Za-z]*
func (lctx *lexCtx) runName() bool {
	r := lctx.next()
	switch {
	case r == 0x5f: // _
	case 0x41 <= r && r <= 0x5a: // A-Z
	case 0x61 <= r && r <= 0x7a: // a-z
	default:
		return false
	}

	for {
		r := lctx.peek()
		switch {
		case r == 0x5f: // _
		case 0x30 <= r && r <= 0x39: // 0-9
		case 0x41 <= r && r <= 0x5a: // A-Z
		case 0x61 <= r && r <= 0x7a: // a-z
		default:
			return true
		}
		lctx.advance()
	}
	return false
}

func (lctx *lexCtx) runString() bool {
	if lctx.next() != '"' {
		return false
	}

	for loop := true; loop; {
		switch lctx.peek() {
		case '"':
			// bail out of the for
			loop = false
			continue
		case '\\':
			if !lctx.runEscapeSequence() {
				return false
			}
		case '\n', '\r':
			return false
		default:
			lctx.advance()
		}
	}

	if lctx.next() != '"' {
		return false
	}
	return true
}

func (lctx *lexCtx) runEscapeSequence() bool {
	if lctx.next() != '\\' {
		return false
	}

	switch lctx.peek() {
	case 'u':
		lctx.advance()
		for i := 0; i < 4; i++ {
			r := lctx.next()
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
		lctx.advance()
		return true
	}
	return false
}
