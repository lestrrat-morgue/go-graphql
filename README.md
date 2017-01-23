# go-graphql

GraphQL for Go

[![Build Status](https://travis-ci.org/lestrrat/go-graphql.png?branch=master)](https://travis-ci.org/lestrrat/go-graphql)

[![GoDoc](https://godoc.org/github.com/lestrrat/go-graphql?status.svg)](https://godoc.org/github.com/lestrrat/go-graphql)

# MOTIVATION

* I really didn't like the implementation of https://github.com/graphql-go/graphql, and couldn't see a way to change it by just sending a few PRs

# STATUS

* Can parse queries
* Can parse schemas
* Can format queries/schemas
* Can create queries/schemas programatically, via raw [models](./model) or [DSL](./dsl)
* Can traverse queries/schemas using [visitor](./visitor)
* TODO: validation
* TODO: actual dispatching

## BENCHMARK

```
% make bench 
go test -tags bench -benchmem -benchtime=5s -bench .
BenchmarkParseOfficial-4      200000         45166 ns/op       18344 B/op        516 allocs/op
BenchmarkParseLestrrat-4      200000         43567 ns/op        7554 B/op        255 allocs/op
BenchmarkParseNeelance-4      200000         37186 ns/op       12379 B/op        291 allocs/op
PASS
ok      github.com/lestrrat/go-graphql  26.462s
```
