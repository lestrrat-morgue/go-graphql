# go-graphql

GraphQL for Go

[![Build Status](https://travis-ci.org/lestrrat/go-graphql.png?branch=master)](https://travis-ci.org/lestrrat/go-graphql)

[![GoDoc](https://godoc.org/github.com/lestrrat/go-graphql?status.svg)](https://godoc.org/github.com/lestrrat/go-graphql)

# MOTIVATION

* I really didn't like the implementation of https://github.com/graphql-go/graphql, and couldn't see a way to change it by just sending a few PRs

# STATUS

## PARSING (DONE)

* Query Document
  * Query
  * Mutation
  * Fragments
  * Directives
* Type Definitions
  * Object
  * Enum
  * Interface
  * Union
  * Input
  * Schema

## DSL

We have a DSL to build the schema. See this file [dsl_test.go](./dsl/dsl_test.go)

## BENCHMARK

```
% go test -tags bench -benchmem -bench .
BenchmarkParseOfficial-4       10000        100539 ns/op       41000 B/op       1147 allocs/op
BenchmarkParseLestrrat-4       20000         73134 ns/op       15453 B/op        583 allocs/op
PASS
ok      github.com/lestrrat/go-graphql  3.262s
```
