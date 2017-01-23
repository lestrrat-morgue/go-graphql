.PHONY: test bench generate

generate:
	go generate

test:
	go test ./...

bench:
	go test -tags bench -benchmem -benchtime=5s -bench .