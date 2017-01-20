package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testlex(src []byte, tokens ...TokenType) (string, func(t *testing.T)) {
	return string(src), func(t *testing.T) {
		l := NewLexer(src)
		seen := 0
		var tok Token
		for i, expected := range tokens {
			if !l.Next(&tok) {
				break
			}
			seen++
			t.Logf("%s", tok)
			if !assert.Equal(t, expected, tok.Type, "token #%d type should match (expected %s, got %s)", i+1, expected, tok.Type) {
				return
			}
		}

		if !assert.Equal(t, len(tokens), seen, "should see expected number of tokens") {
			return
		}
	}
}

func TestLex(t *testing.T) {
	t.Run(testlex([]byte("Hello"), NAME, EOF))
	t.Run(testlex([]byte("!"), BANG, EOF))
	t.Run(testlex([]byte("..."), SPREAD, EOF))
	t.Run(testlex([]byte("-123.142"), FLOAT, EOF))
	t.Run(testlex([]byte("+123.142"), FLOAT, EOF))
	t.Run(testlex([]byte("123.142"), FLOAT, EOF))
	t.Run(testlex([]byte("123e+142"), FLOAT, EOF))
	t.Run(testlex([]byte(`"Hello\u0020World"`), STRING, EOF))
}
