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
go test -tags bench -benchmem -bench .
BenchmarkParseOfficial-4       10000        102123 ns/op       41000 B/op       1147 allocs/op
BenchmarkParseLestrrat-4       20000         77163 ns/op       16530 B/op        599 allocs/op
PASS
ok      github.com/lestrrat/go-graphql  3.399s
```
