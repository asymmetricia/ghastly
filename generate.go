//go:generate mkdir -p search
//go:generate go run github.com/mna/pigeon -o search/search.go search.peg
//go:generate go run golang.org/x/tools/cmd/goimports -w search/search.go
package main
