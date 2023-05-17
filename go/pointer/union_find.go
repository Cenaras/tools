package pointer

import (
	"fmt"
)

func (a *analysis) find(x nodeid) nodeid {
	xn := a.nodes[x]
	rep := xn.rep
	if rep != x {
		rep = a.find(rep) // simple path compression
		xn.rep = rep
	}
	return rep
}

func unify(a *analysis, inCycles map[nodeid]struct{}, r map[nodeid]nodeid) {
	//var stale nodeset
	var deltaSpace []int
	for id := range inCycles {
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
