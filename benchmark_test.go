// +build bench

package graphql_test

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"testing"

	official "github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/source"
	lestrrat "github.com/lestrrat/go-graphql/parser"
	"golang.org/x/net/context"
)

func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()
}

const src = `query {
  me {
    name
  }
}

query HeroNameAndFriends($episode: Episode) {
  hero(episode: $episode) {
    name
    friends {
      name
    }
  }
}

query {
  human(id: "1000") {
    name
    height
  }
}

query {
  human(id: "1000") {
    name
    height(unit: FOOT)
  }
}

query {
  empireHero: hero(episode: EMPIRE) {
    name
  }
  jediHero: hero(episode: JEDI) {
    name
  }
}

query {
  leftComparison: hero(episode: EMPIRE) {
    ...comparisonFields
  }
  rightComparison: hero(episode: JEDI) {
    ...comparisonFields
  }
}

fragment comparisonFields on Character {
  name
  appearsIn
  friends {
    name
  }
}

query Hero($episode: Episode, $withFriends: Boolean!) {
  hero(episode: $episode) {
    name
    friends @include(if: $withFriends) {
      name
    }
  }
}

mutation CreateReviewForEpisode($ep: Episode!, $review: ReviewInput!) {
  createReview(episode: $ep, review: $review) {
    stars
    commentary
  }
}

query HeroForEpisode($ep: Episode!) {
  hero(episode: $ep) {
    name
    ... on Droid {
      primaryFunction
    }
    ... on Human {
      height
    }
  }
}

query {
  search(text: "an") {
    __typename
    ... on Human {
      name
    }
    ... on Droid {
      name
    }
    ... on Starship {
      name
    }
  }
}

query {
  nearestThing(location: {
    lon: 12.43
    lat: -53.211
  })
}

type Person1 {
  name: String
  age: Int
  picture: Url
  relationship: Person
}

type Person2 {
  name: String
  picture(size: Int): Url
}`

func BenchmarkParseOfficial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := official.Parse(official.ParseParams{
			Source: &source.Source{
				Body: []byte(src),
				Name: "benchmark",
			},
		})
		if err != nil {
			b.Logf("error: %s", err)
			return
		}
	}
}

func BenchmarkParseLestrrat(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p := lestrrat.New()
	for i := 0; i < b.N; i++ {
		_, err := p.Parse(ctx, []byte(src))
		if err != nil {
			b.Logf("error: %s", err)
			return
		}
	}
}
