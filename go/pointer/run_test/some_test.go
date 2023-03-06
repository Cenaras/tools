package main

import (
	"fmt"
	"testing"

	"golang.org/x/tools/go/pointer"
)

func Test(t *testing.T) {
	// Don't assert @pointsto(t) since its label contains a fragile line number.
	run_engine("obj_sens.go", &pointer.KObjNHeap{K: 1, N: 0})
	t.Fail()
}

func TestAll(t *testing.T) {
	test := []string{"dynamic_sens.go", "myprog.go", "callsite_dup.go", "heap_cloning.go", "invoke_sens.go", "callsite_three_deep.go"}
	for _, s := range test {
		s := s
		t.Run(s, func(t *testing.T) {
			t.Parallel()
			run_engine(s, nil)
		})
	}
}

type S struct {
	i int
}

func TestMap(t *testing.T) {
	m := make(map[S]bool)
	p := S{1}
	m[p] = true
	q := S{1}
	fmt.Print(m[q])
	t.Fail()
}
