package pointer

import (
	"fmt"
)

type nuutila struct {
	a        *analysis
	I        int
	D        map[nodeid]int
	R        map[nodeid]nodeid
	C        nodeset
	S        []nodeid
	T        []nodeid
	InCycles nodeset
}

func (nuu *nuutila) visit(v nodeid) {
	nuu.I++
	nuu.D[v] = nuu.I
	nuu.R[v] = v
	var deltaSpace []int
	for _, x := range nuu.a.nodes[v].solve.copyTo.AppendTo(deltaSpace) {
		w := nuu.a.find(nodeid(x))
		if nuu.D[w] == 0 {
			nuu.visit(w)
		}
		if !nuu.C.Has(int(w)) {
			if nuu.D[nuu.R[v]] >= nuu.D[nuu.R[w]] {
				nuu.R[v] = nuu.R[w]
				nuu.InCycles.add(v)
			}
		}
	}
	if nuu.R[v] == v {
		nuu.C.add(v)
		for len(nuu.S) != 0 {
			w := nuu.S[len(nuu.S)-1]
			if nuu.D[w] <= nuu.D[v] {
				break
			} else {
				nuu.S = nuu.S[:len(nuu.S)-1]
				nuu.C.add(w)
				nuu.R[w] = v
				nuu.InCycles.add(w)
			}
		}
		nuu.T = append(nuu.T, v)
	} else {
		nuu.S = append(nuu.S, v)
	}
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

func unify(a *analysis, inCycles *nodeset, r map[nodeid]nodeid) {
	var stale nodeset
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

			if xsolve.pts.addAll(&ysolve.pts) {
				a.addWork(rep)
			}

			for _, w := range ysolve.copyTo.AppendTo(deltaSpace) {
				xsolve.copyTo.add(a.find(nodeid(w)))
				a.nodes[w].solve.pts.addAll(&xsolve.prevPTS)
			}
			xsolve.checkedLazy.addAll(&ysolve.checkedLazy)
			xsolve.complex = append(xsolve.complex, ysolve.complex...)
			if !xsolve.prevPTS.IsEmpty() {
				stale.add(rep)
			}

			y.solve = x.solve
			y.rep = rep
		}
	}
	for _, id := range stale.AppendTo(deltaSpace) {
		n := a.nodes[id]
		a.solveConstraints(n, &n.solve.prevPTS, false)
	}
}
