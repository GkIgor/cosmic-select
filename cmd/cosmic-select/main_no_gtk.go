//go:build !gtk4

package main

import "fmt"

func main() {
	fmt.Println("COSMIC Select requires a GTK4 build; use: go run -tags gtk4 ./cmd/cosmic-select")
}
