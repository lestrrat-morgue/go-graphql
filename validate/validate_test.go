package validate_test

import (
	"testing"
	"time"

	"github.com/lestrrat/go-graphql/parser"
	"github.com/lestrrat/go-graphql/schema"
	"github.com/lestrrat/go-graphql/validate"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestValidate(t *testing.T) {
	t.Run("Cannot spread fragment within itself", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		src := `{
  hero {
    ...NameAndAppearancesAndFriends
  }
}

fragment NameAndAppearancesAndFriends on Character {
  name
  appearsIn
  friends {
    ...NameAndAppearancesAndFriends
  }
}`
		p := parser.New()
		doc, err := p.ParseString(ctx, src)
		if !assert.NoError(t, err, "p.Parse should succed") {
			return
		}

		if !assert.Error(t, validate.Validate(ctx, schema.StarWars, doc), "document should fail to validate") {
			return
		}
	})
}
