package parser

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testlex(src []byte, tokens ...TokenType) (string, func(t *testing.T)) {
	return string(src), func(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ch := make(chan *Token)
	go lex(ctx, src, ch)

	for i, expected := range tokens {
		select {
		case <-ctx.Done():
			t.Errorf("timeout reached")
			return
		case tok := <-ch:
			t.Logf("%s", tok)
			if !assert.Equal(t, expected, tok.Type, "token #%d type should match (expected %s, got %s)", i+1, expected, tok.Type) {
				return
			}
		}
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
