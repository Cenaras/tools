package pointer

type nuutila struct {
	a        *analysis
	I        int
	D        map[nodeid]int
	R        map[nodeid]nodeid
	C        map[nodeid]struct{}
	S        []nodeid
	T        []nodeid
	InCycles map[nodeid]struct{}
}

func (nuu *nuutila) visitAll() {
	//for x, _ := range nuu.a.nodes {
	for x, _ := range nuu.a.work {
		if id := nodeid(x); nuu.a.find(id) == id && nuu.D[id] == 0 {
			nuu.visit(id)
		}
	}
	//nuu.a.work.Clear()
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
		if _, ok := nuu.C[w]; !ok {
			if nuu.D[nuu.R[v]] >= nuu.D[nuu.R[w]] {
				nuu.R[v] = nuu.R[w]
				nuu.InCycles[v] = struct{}{}
			}
		}
	}
	if nuu.R[v] == v {
		nuu.C[v] = struct{}{}
		for len(nuu.S) != 0 {
			w := nuu.S[len(nuu.S)-1]
			if nuu.D[w] <= nuu.D[v] {
				break
			} else {
				nuu.S = nuu.S[:len(nuu.S)-1]
				nuu.C[w] = struct{}{}
				nuu.R[w] = v
				nuu.InCycles[w] = struct{}{}
			}
		}
		nuu.T = append(nuu.T, v)
	} else {
		nuu.S = append(nuu.S, v)
	}
}
