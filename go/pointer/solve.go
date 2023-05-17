// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pointer

// This file defines a naive Andersen-style solver for the inclusion
// constraint system.

import (
	"fmt"
	"go/types"
)

type solverState struct {
	complex []*waveConstraint // complex constraints attached to this node
	copyTo  nodeset           // simple copy constraint edges
	pts     nodeset           // points-to set of this node
	prevPTS nodeset           // pts(n) in previous iteration (for difference propagation)
}

type waveConstraint struct {
	constraint constraint
	cache      nodeset
}

func (a *analysis) puSolve() {
	start("Solving")
	if a.log != nil {
		fmt.Fprintf(a.log, "\n\n==== Solving constraints\n\n")
	}
	//first := true
	i := 0
	for {
		//start := time.Now()
		a.processNewConstraints()
		//fmt.Fprintf(os.Stdout, "Elapsed time for new constraints: %f\n", time.Since(start).Seconds())
		/*
			if first {
				first = false
				for id, _ := range a.nodes {
					a.cycleCandidates.add(nodeid(id))
				}
			}
		*/
		//start = time.Now()
		//Detect and collapse cycles
		nuu := &nuutila{a: a, I: 0, D: make(map[nodeid]int), R: make(map[nodeid]nodeid), C: make(map[nodeid]struct{}), InCycles: make(map[nodeid]struct{})}
		nuu.visitAll()
		//fmt.Fprintf(os.Stdout, "Elapsed time for detect cycles: %f\n", time.Since(start).Seconds())
		//start = time.Now()
		unify(a, nuu.InCycles, nuu.R)
		//fmt.Fprintf(os.Stdout, "Elapsed time for collapse cycles: %f\n", time.Since(start).Seconds())

		//start = time.Now()
		// Wave propagation
		t := nuu.T
		for len(t) != 0 {
			v := t[len(t)-1]
			t = t[:len(t)-1]
			if _, ok := a.work[v]; !ok {
				continue
			}
			nsolve := a.nodes[v].solve
			var diff nodeset
			diff.Difference(&nsolve.pts.Sparse, &nsolve.prevPTS.Sparse)
			nsolve.prevPTS.Copy(&nsolve.pts.Sparse)
			for _, w := range nsolve.copyTo.AppendTo(a.deltaSpace) {
				if a.nodes[nodeid(w)].solve.pts.addAll(&diff) {
					a.addWork(nodeid(w))
				}
			}
		}

		//fmt.Fprintf(os.Stdout, "Elapsed time for label propagation: %f\n", time.Since(start).Seconds())

		//start = time.Now()
		//var changed bool = false
		complexWork := a.work
		a.work = make(map[nodeid]struct{})
		for n, _ := range complexWork {
			nsolve := a.nodes[n].solve
			for _, c := range nsolve.complex {
				var diff nodeset
				diff.Difference(&nsolve.pts.Sparse, &c.cache.Sparse)
				c.cache.Copy(&nsolve.pts.Sparse)
				if a.log != nil {
					fmt.Fprintf(a.log, "\t\tconstraint %s\n", c.constraint)
				}
				c.constraint.solve(a, &diff)
			}
		}
		//fmt.Fprintf(os.Stdout, "Elapsed time for complex constraints: %f\n", time.Since(start).Seconds())
		if len(a.work) == 0 && len(a.constraints) == 0 {
			break
		}
		i++
		//fmt.Fprintf(os.Stdout, "Loop iteration %d\n", i)

	}

	if !a.nodes[0].solve.pts.IsEmpty() {
		panic(fmt.Sprintf("pts(0) is nonempty: %s", &a.nodes[0].solve.pts))
	}

	// Release working state (but keep final PTS).
	for _, n := range a.nodes {
		solve := n.solve
		solve.complex = nil
		solve.copyTo.Clear()
		solve.prevPTS.Clear()
	}

	if a.log != nil {
		fmt.Fprintf(a.log, "Solver done\n")

		// Dump solution.
		for i, n := range a.nodes {
			if !n.solve.pts.IsEmpty() {
				fmt.Fprintf(a.log, "pts(n%d) = %s : %s\n", i, &n.solve.pts, n.typ)
			}
		}
	}
}

// processNewConstraints takes the new constraints from a.constraints
// and adds them to the graph, ensuring
// that new constraints are applied to pre-existing labels and
// that pre-existing constraints are applied to new labels.
func (a *analysis) processNewConstraints() {
	// Take the slice of new constraints.
	// (May grow during call to solveConstraints.)
	constraints := a.constraints
	a.constraints = nil

	// Initialize points-to sets from addr-of (base) constraints.
	for _, c := range constraints {
		if c, ok := c.(*addrConstraint); ok {
			dst := a.nodes[c.dst]
			if dst.solve.pts.add(c.src) {
				// Populate the worklist with nodes that point to
				// something initially (due to addrConstraints) and
				// have other constraints attached.
				// (A no-op in round 1.)
				a.addWork(c.dst)
			}
		}
	}

	// Attach simple (copy) and complex constraints to nodes.
	//var stale nodeset
	for _, c := range constraints {
		var id nodeid
		switch c := c.(type) {
		case *addrConstraint:
			// base constraints handled in previous loop
			continue
		case *copyConstraint:
			// simple (copy) constraint
			id = c.src
			a.nodes[id].solve.copyTo.add(c.dst)
			a.nodes[c.dst].solve.pts.addAll(&a.nodes[id].solve.prevPTS)
			a.addWork(c.dst)
		default:
			// complex constraint
			id = c.ptr()
			solve := a.nodes[id].solve
			solve.complex = append(solve.complex, &waveConstraint{constraint: c})
			a.addWork(id)
			//a.waveConstraints = append(a.waveConstraints, &waveConstraint{constraint: c})
		}
		/*
			if n := a.nodes[id]; !n.solve.pts.IsEmpty() {
				if !n.solve.prevPTS.IsEmpty() {
					stale.add(id)
				}
				a.addWork(id)
			}
		*/
	}
	/*
		// Apply new constraints to pre-existing PTS labels.
		var space [50]int
		for _, id := range stale.AppendTo(space[:0]) {
			n := a.nodes[nodeid(id)]
			a.solveConstraints(n, &n.solve.prevPTS)
		}
	*/
}

// solveConstraints applies each resolution rule attached to node n to
// the set of labels delta.  It may generate new constraints in
// a.constraints.
func (a *analysis) solveConstraints(n *node, delta *nodeset) {
	if delta.IsEmpty() {
		return
	}

	// Process complex constraints dependent on n.
	for _, c := range n.solve.complex {
		if a.log != nil {
			fmt.Fprintf(a.log, "\t\tconstraint %s\n", c.constraint)
		}
		//c.solve(a, delta)
	}

	// Process copy constraints.
	var copySeen nodeset
	for _, x := range n.solve.copyTo.AppendTo(a.deltaSpace) {
		mid := nodeid(x)
		if copySeen.add(mid) {
			if a.nodes[mid].solve.pts.addAll(delta) {
				a.addWork(mid)
			}
		}
	}
}

// addLabel adds label to the points-to set of ptr and reports whether the set grew.
func (a *analysis) addLabel(ptr, label nodeid) bool {
	b := a.nodes[ptr].solve.pts.add(label)
	if b && a.log != nil {
		fmt.Fprintf(a.log, "\t\tpts(n%d) += n%d\n", ptr, label)
	}
	return b
}

func (a *analysis) addWork(id nodeid) {
	a.work[a.find(id)] = struct{}{}
	if a.log != nil {
		fmt.Fprintf(a.log, "\t\twork: n%d\n", id)
	}
}

// onlineCopy adds a copy edge.  It is called online, i.e. during
// solving, so it adds edges and pts members directly rather than by
// instantiating a 'constraint'.
//
// The size of the copy is implicitly 1.
// It returns true if pts(dst) changed.
func (a *analysis) onlineCopy(dst, src nodeid) bool {
	if dst != src {
		if nsrc := a.nodes[src]; nsrc.solve.copyTo.add(dst) {
			if a.log != nil {
				fmt.Fprintf(a.log, "\t\t\tdynamic copy n%d <- n%d\n", dst, src)
			}

			//a.addWork(dst)
			return a.nodes[dst].solve.pts.addAll(&nsrc.solve.prevPTS)
		}
	}
	return false
}

// Returns sizeof.
// Implicitly adds nodes to worklist.
//
// TODO(adonovan): now that we support a.copy() during solving, we
// could eliminate onlineCopyN, but it's much slower.  Investigate.
func (a *analysis) onlineCopyN(dst, src nodeid, sizeof uint32) uint32 {
	for i := uint32(0); i < sizeof; i++ {
		if a.onlineCopy(dst, src) {
			a.addWork(dst)
		}
		src++
		dst++
	}
	return sizeof
}

func (c *loadConstraint) solve(a *analysis, delta *nodeset) {
	var changed bool
	for _, x := range delta.AppendTo(a.deltaSpace) {
		k := nodeid(x)
		koff := k + nodeid(c.offset)
		if a.onlineCopy(c.dst, koff) {
			changed = true
		}
	}
	if changed {
		a.addWork(c.dst)
	}
}

func (c *storeConstraint) solve(a *analysis, delta *nodeset) {
	for _, x := range delta.AppendTo(a.deltaSpace) {
		k := nodeid(x)
		koff := k + nodeid(c.offset)
		if a.onlineCopy(koff, c.src) {
			a.addWork(koff)
		}
	}
}

func (c *offsetAddrConstraint) solve(a *analysis, delta *nodeset) {
	dst := a.nodes[c.dst]
	for _, x := range delta.AppendTo(a.deltaSpace) {
		k := nodeid(x)
		if dst.solve.pts.add(k + nodeid(c.offset)) {
			a.addWork(c.dst)
		}
	}
}

func (c *typeFilterConstraint) solve(a *analysis, delta *nodeset) {
	for _, x := range delta.AppendTo(a.deltaSpace) {
		ifaceObj := nodeid(x)
		tDyn, _, indirect := a.taggedValue(ifaceObj)
		if indirect {
			// TODO(adonovan): we'll need to implement this
			// when we start creating indirect tagged objects.
			panic("indirect tagged object")
		}

		if types.AssignableTo(tDyn, c.typ) {
			if a.addLabel(c.dst, ifaceObj) {
				a.addWork(c.dst)
			}
		}
	}
}

func (c *untagConstraint) solve(a *analysis, delta *nodeset) {
	predicate := types.AssignableTo
	if c.exact {
		predicate = types.Identical
	}
	for _, x := range delta.AppendTo(a.deltaSpace) {
		ifaceObj := nodeid(x)
		tDyn, v, indirect := a.taggedValue(ifaceObj)
		if indirect {
			// TODO(adonovan): we'll need to implement this
			// when we start creating indirect tagged objects.
			panic("indirect tagged object")
		}

		if predicate(tDyn, c.typ) {
			// Copy payload sans tag to dst.
			//
			// TODO(adonovan): opt: if tDyn is
			// nonpointerlike we can skip this entire
			// constraint, perhaps.  We only care about
			// pointers among the fields.
			a.onlineCopyN(c.dst, v, a.sizeof(tDyn))
		}
	}
}

func (c *invokeConstraint) solve(a *analysis, delta *nodeset) {
	for _, x := range delta.AppendTo(a.deltaSpace) {
		ifaceObj := nodeid(x)
		tDyn, v, indirect := a.taggedValue(ifaceObj)
		if indirect {
			// TODO(adonovan): we may need to implement this if
			// we ever apply invokeConstraints to reflect.Value PTSs,
			// e.g. for (reflect.Value).Call.
			panic("indirect tagged object")
		}

		// Look up the concrete method.
		fn := a.prog.LookupMethod(tDyn, c.method.Pkg(), c.method.Name())
		if fn == nil {
			panic(fmt.Sprintf("n%d: no ssa.Function for %s", c.iface, c.method))
		}
		sig := fn.Signature

		fnObj := a.globalobj[fn] // dynamic calls use shared contour
		if fnObj == 0 {
			// a.objectNode(fn) was not called during gen phase.
			panic(fmt.Sprintf("a.globalobj[%s]==nil", fn))
		}

		// Make callsite's fn variable point to identity of
		// concrete method.  (There's no need to add it to
		// worklist since it never has attached constraints.)
		a.addLabel(c.params, fnObj)

		// Extract value and connect to method's receiver.
		// Copy payload to method's receiver param (arg0).
		arg0 := a.funcParams(fnObj)
		recvSize := a.sizeof(sig.Recv().Type())
		a.onlineCopyN(arg0, v, recvSize)

		src := c.params + 1 // skip past identity
		dst := arg0 + nodeid(recvSize)

		// Copy caller's argument block to method formal parameters.
		paramsSize := a.sizeof(sig.Params())
		a.onlineCopyN(dst, src, paramsSize)
		src += nodeid(paramsSize)
		dst += nodeid(paramsSize)

		// Copy method results to caller's result block.
		resultsSize := a.sizeof(sig.Results())
		a.onlineCopyN(src, dst, resultsSize)
	}
}

func (c *addrConstraint) solve(a *analysis, delta *nodeset) {
	panic("addr is not a complex constraint")
}

func (c *copyConstraint) solve(a *analysis, delta *nodeset) {
	panic("copy is not a complex constraint")
}
