package pointer

// Define functions on Graph, give arguments to methods

type UFNode struct {
	solverState *solverState // Solver state of the underlying node represented.
	parent      *UFNode      // The representative of this node - nil if root
	rank        int          // The rank of the node
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

func (ufNode *UFNode) union(other *UFNode) {
	x := ufNode.findParent()
	y := other.findParent()

	if x == y {
		return
	}

	// Set the root of this subset to be node with lowest rank.
	if x.rank < y.rank {
		x, y = y, x
	}

	y.parent = x
	if x.rank == y.rank {
		x.rank += 1
	}
}

// Do a bunched unify instead of a set of nodes to unify rather than this.
func (ufNode *UFNode) unify(other *UFNode) {
	x := ufNode.find()
	y := other.find()

	x.pts.UnionWith(&y.pts.Sparse)
	x.complex = append(x.complex, y.complex...) // TODO: Dedupe
	//if y.complex != nil && !x.prevPTS.IsEmpty()

}
