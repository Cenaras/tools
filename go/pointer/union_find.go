package pointer

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
}

// Do a bunched unify instead of a set of nodes to unify rather than this.
func unify(a *analysis, inCycles []nodeid, r map[nodeid]nodeid) {
	var stale nodeset
	for _, v := range inCycles {
		if v != r[v] {
			xsolve := a.nodes[r[v]].solve
			x := xsolve.find()
			ysolve := a.nodes[v].solve
			y := ysolve.find()

			if x.pts.UnionWith(&y.pts.Sparse) {
				a.addWork(x.id)
			}
			x.complex = append(x.complex, y.complex...) // TODO: Dedupe
			if len(y.complex) != 0 && !x.prevPTS.IsEmpty() {
				stale.add(x.id)
			}
			xsolve.union(ysolve)
		}
	}
	var deltaSpace []int
	for _, id := range stale.AppendTo(deltaSpace) {
		n := a.nodes[id]
		a.solveConstraints(n, &n.solve.find().prevPTS)
	}
}
