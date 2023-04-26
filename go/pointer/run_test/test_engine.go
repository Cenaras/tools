// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"sort"

	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func main() {
	filename := os.Args[2]
	run_engine(filename, nil)
}

func run_engine(filename string, contextStrategy pointer.ContextStrategy) {

	var conf loader.Config

	// Parse the input file, a string.
	// (Command-line tools should use conf.FromArgs.)

	file, err := conf.ParseFile(filename, nil)
	if err != nil {
		fmt.Print(err) // parse error
		return
	}

	// Create single-file main package and import its dependencies.
	conf.CreateFromFiles("main", file)

	iprog, err := conf.Load()
	if err != nil {
		fmt.Print(err) // type error in some package
		return
	}

	// Create SSA-form program representation.
	prog := ssautil.CreateProgram(iprog, ssa.InstantiateGenerics)
	mainPkg := prog.Package(iprog.Created[0].Pkg)

	// Build SSA code for bodies of all functions in the whole program.
	prog.Build()

	//var log bytes.Buffer
	// Configure the pointer analysis to build a call-graph.
	config := &pointer.Config{
		Mains:           []*ssa.Package{mainPkg},
		BuildCallGraph:  true,
		Log:             os.Stdout,
		ContextStrategy: contextStrategy,
	}

	allfun := ssautil.AllFunctions(prog)
	for fun := range allfun {
		for _, param := range fun.Params {
			if pointer.CanPoint(param.Type()) {
				config.AddQuery(param)
			}
		}
		for _, fv := range fun.FreeVars {
			if pointer.CanPoint(fv.Type()) {
				config.AddQuery(fv)
			}
		}

		for _, b := range fun.Blocks {
			for _, instr := range b.Instrs {
				switch v := instr.(type) {
				case *ssa.Call:
					common := v.Common()
					if pointer.CanPoint(v.Type()) {
						config.AddQuery(v)
					}
					if pointer.CanPoint(common.Value.Type()) {
						config.AddQuery(common.Value)
					}
				case *ssa.Range:
				case ssa.Value:
					if pointer.CanPoint(v.Type()) {
						config.AddQuery(v)
					}
				}
			}
		}
	}
	// Run the pointer analysis.
	result, err := pointer.Analyze(config)
	if err != nil {
		panic(err) // internal error in pointer analysis
	}

	// Find edges originating from the main package.
	// By converting to strings, we de-duplicate nodes
	// representing the same function due to context sensitivity.
	var edges []string
	callgraph.GraphVisitEdges(result.CallGraph, func(edge *callgraph.Edge) error {
		caller := edge.Caller.Func
		if false {
			edges = append(edges, fmt.Sprint(caller, " --> ", edge.Callee.Func))
		}
		return nil
	})

	// Print the edges in sorted order.
	sort.Strings(edges)
	for _, edge := range edges {
		fmt.Println(edge)
	}
	fmt.Println()

	for v, q := range result.Queries {

		fmt.Printf("%s at %s may point to:\n", v.String(), prog.Fset.Position(v.Pos()))
		var labels []string
		for _, l := range q.PointsTo().Labels() {
			label := fmt.Sprintf("  %s: %s", prog.Fset.Position(l.Pos()), l)
			labels = append(labels, label)
		}
		sort.Strings(labels)
		for _, label := range labels {
			fmt.Println(label)
		}
		fmt.Println("")
	}

}
