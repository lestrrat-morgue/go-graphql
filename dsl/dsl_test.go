package dsl_test

import (
	"bytes"
	"testing"

	. "github.com/lestrrat/go-graphql/dsl"
	"github.com/lestrrat/go-graphql/format"
	"github.com/stretchr/testify/assert"
)

func TestStarWars(t *testing.T) {
	var episodeEnum = Enum(
		Name(`Episode`),
		Description(`One of the films in the Star Wars Trilogy`),
		EnumValue(
			`NEWHOPE`,
			IntValue(4),
			Description(`Released in 1977.`),
		),
		EnumValue(
			`EMPIRE`,
			IntValue(5),
			Description(`Released in 1980.`),
		),
		EnumValue(
			`JEDI`,
			IntValue(6),
			Description(`Released in 1983.`),
		),
	)

	var characterInterfaceDef = Interface(`Character`)
	var humanTypeDef = Object(`Human`)
	var droidTypeDef = Object(`Droid`)
	var queryTypeDef = Object(`Query`)

	var humanType = humanTypeDef.Configure(
		Implements(characterInterfaceDef.Type()),
		ObjectField(
			`id`,
			NotNull(String()),
			Description(`The id of the human.`),
		),
		ObjectField(
			`name`,
			String(),
			Description(`The name of the human.`),
		),
		ObjectField(
			`friends`,
			List(characterInterfaceDef.Type()),
			Description(`'The friends of the human, or an empty list if they have none.`),
			// resolve: human => getFriends(human),
		),
		ObjectField(
			`appearsIn`,
			List(episodeEnum),
			Description(`Which movies they appear in.`),
		),
		ObjectField(
			`homePlanet`,
			String(),
			Description(`The home planet of the human, or null if unknown.`),
		),
		ObjectField(
			`secretBackstory`,
			String(),
			Description(`Where are they from and how they came to be who they are.`),
			// resolve() {
			//  throw new Error('secretBackstory is secret.');
			//},
		),
	).Type()

	var droidType = droidTypeDef.Configure(
		Implements(characterInterfaceDef.Type()),
		Description(`A mechanical creature in the Star Wars universe.`),
		ObjectField(
			`id`,
			NotNull(String()),
			Description(`The id of the droid.`),
		),
		ObjectField(
			`name`,
			String(),
			Description(`The name of the droid.`),
		),
		ObjectField(
			`friends`,
			List(characterInterfaceDef.Type()),
			Description(`'The friends of the droid, or an empty list if they have none.`),
			// resolve: droid => getFriends(droid),
		),
		ObjectField(
			`appearsIn`,
			List(episodeEnum),
			Description(`Which movies they appear in.`),
		),
		ObjectField(
			`homePlanet`,
			String(),
			Description(`The home planet of the droid, or null if unknown.`),
		),
		ObjectField(
			`secretBackstory`,
			String(),
			Description(`Where are they from and how they came to be who they are.`),
			// resolve() {
			//  throw new Error('secretBackstory is secret.');
			//},
		),
		ObjectField(
			`primaryFunction`,
			String(),
			Description(`The primary function of the droid.`),
		),
	).Type()

	var characterInterface = characterInterfaceDef.Configure(
		Description(`A character in the Star Wars Trilogy`),
		InterfaceField(
			`id`,
			NotNull(String()),
			Description(`The id of the character.`),
		),
		InterfaceField(
			`name`,
			String(),
			Description(`The name of the character.`),
		),
		InterfaceField(
			`friends`,
			List(characterInterfaceDef.Type()),
			Description(`The friends of the character, or an empty list if they have none.`),
		),
		InterfaceField(
			`appearsIn`,
			List(episodeEnum),
			Description(`Which movies they appear in.`),
		),
		InterfaceField(
			`secretBackstory`,
			String(),
			Description(`All secrets about their past.`),
		),
	).Type()

	/*
	     resolveType(character) {
	       if (character.type === 'Human') {
	         return humanType;
	       }
	       if (character.type === 'Droid') {
	         return droidType;
	       }
	     }
	   });
	*/

	var queryType = queryTypeDef.Configure(
		ObjectField(
			`hero`,
			characterInterface,
			ObjectFieldArgument(
				`episode`,
				episodeEnum,
				Description(`If omitted, returns the hero of the whole saga. If provided, returns the hero of that particular episode.`),
			),
			//resolve: (root, { episode }) => getHero(episode),
		),
		ObjectField(
			`human`,
			humanType,
			ObjectFieldArgument(
				`id`,
				NotNull(String()),
				Description(`id of the human`),
			),
			// resolve: (root, { id }) => getHuman(id),
		),
		ObjectField(
			`droid`,
			droidType,
			ObjectFieldArgument(
				`id`,
				NotNull(String()),
				Description(`id of the droid`),
			),
			// resolve: (root, { id }) => getDroid(id),
		),
	).Type()

	schema := Schema(
		episodeEnum,
		characterInterface,
		humanType,
		droidType,
		queryType,
	)

	var buf bytes.Buffer
	if !assert.NoError(t, format.GraphQL(&buf, schema), "format.GraphQL succeeds") {
		return
	}
	t.Logf("%s", buf.String())
}
