package pointer

import (
	"fmt"
)

// Define functions on Graph, give arguments to methods

type UFNode struct {
	solverState *solverState // Solver state of the underlying node represented.
	parent      *UFNode      // The representative of this node - nil if root
}

func (uFNode *UFNode) find() *solverState {
	return uFNode.findParent().solverState
}

// Returns the representative of this node.
func (ufNode *UFNode) findParent() *UFNode {
	if ufNode.parent != nil {
		ufNode.parent = ufNode.parent.findParent()
		return ufNode.parent
	}
	// Somewhat inefficient as we essentially traverse twice, down --> top, top --> down.
	return ufNode
}

// ufNode becomes the parent of other
func (ufNode *UFNode) union(other *UFNode) {
	x := ufNode.findParent()
	y := other.findParent()

	if x == y {
		return
	}

	y.parent = x
	// No need to keep old solverstate
	y.solverState = nil
}

func (a *analysis) find(x nodeid) nodeid {
	xn := a.nodes[x]
	rep := xn.rep
	if rep != x {
		rep = a.find(rep) // simple path compression
		xn.rep = rep
	}
	return rep
}

// Do a bunched unify instead of a set of nodes to unify rather than this.
func unify(a *analysis, inCycles *nodeset, r map[nodeid]nodeid) {
	//var stale nodeset
	var deltaSpace []int
	for _, id := range inCycles.AppendTo(deltaSpace) {
		v := a.find(nodeid(id))
		rep := a.find(r[v])
		if v != rep {
			x := a.nodes[rep]
			xsolve := x.solve
			y := a.nodes[v]
			ysolve := y.solve
			if a.log != nil {
				fmt.Fprintf(a.log, "Unifying %d into %d\n", v, rep)
			}

			xsolve.pts.addAll(&ysolve.pts)

			for _, w := range ysolve.copyTo.AppendTo(deltaSpace) {
				xsolve.copyTo.add(a.find(nodeid(w)))
				a.nodes[w].solve.pts.addAll(&xsolve.prevPTS)
			}
			y.solve = x.solve
			y.rep = rep
		}
	}
}
