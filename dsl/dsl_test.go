package dsl_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/lestrrat/go-graphql/format"
	"github.com/lestrrat/go-graphql/schema"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestStarWars(t *testing.T) {
	const expected = `enum Episode {
  NEWHOPE
  EMPIRE
  JEDI
}

interface Character {
  id: String!
  name: String
  friends: [Character]
  appearsIn: [Episode]
  secretBackstory: String
}

type Human implements Character {
  id: String!
  name: String
  friends: [Character]
  appearsIn: [Episode]
  homePlanet: String
  secretBackstory: String
}

type Droid implements Character {
  id: String!
  name: String
  friends: [Character]
  appearsIn: [Episode]
  secretBackstory: String
  primaryFunction: String
}

type Query {
  hero(episode: Episode): Character
  human(id: String!): Human
  droid(id: String!): Droid
}

schema {
  query: Query
  types: [Episode, Character, Human, Droid]
}`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var buf bytes.Buffer
	if !assert.NoError(t, format.GraphQL(ctx, &buf, schema.StarWars), "format.GraphQL succeeds") {
		return
	}

	if !assert.Equal(t, expected, buf.String(), "generated schema matches") {
		t.Logf("expected:\n%s", expected)
		t.Logf("actual:\n%s", buf.String())
		return
	}
}
