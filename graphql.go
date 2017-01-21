//go:generate go run internal/cmd/gentokens.go -- parser/tokens.go
//go:generate go run internal/cmd/genkinds.go -- parser/kinds.go
//go:generate go run internal/cmd/geniters.go -- parser/iterators.go

package graphql
