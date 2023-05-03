//go:build ignore
// +build ignore

package main

var a, b int

type S struct{}

func id[T any](x *T) *T {
	return x
}

func param() {
	var c *int = id(&a)
	var d *S = id(&S{})
	var e *int = id(&b)
	print(c)
	print(d)
	print(e)
}

func main() {
	param()
}
