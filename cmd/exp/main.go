package main

import (
	"context"
	"fmt"
	"strings"
)

type ctxKey string

const (
	colorKey ctxKey = "color"
)

func main() {

	ctx := context.Background()
	ctx = context.WithValue(ctx, colorKey, "blue")

	value := ctx.Value(colorKey)
	str, ok := value.(string)
	if !ok {
		panic("BAD TYPE CONV")
	}
	fmt.Println(strings.HasPrefix(str, "b"))
}
