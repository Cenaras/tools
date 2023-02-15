package pointer

import (
	"go/token"

	"golang.org/x/tools/go/ssa"
)

type kCallsiteContext struct {
	targets nodeid
	instr   []ssa.CallInstruction
}

func (c *kCallsiteContext) Targets() nodeid {
	return c.targets
}

func (c *kCallsiteContext) SetTargets(targets nodeid) {
	c.targets = targets
}

func (c *kCallsiteContext) Instr() ssa.CallInstruction {
	if len(c.instr) != 0 {
		return c.instr[len(c.instr)-1]
	} else {
		return nil
	}
}

func (c *kCallsiteContext) ShouldUseContext(fn *ssa.Function, a *analysis) bool {
	if a.findIntrinsic(fn) != nil {
		return true // treat intrinsics context-sensitively
	}
	if len(fn.Blocks) != 1 {
		return false // too expensive
	}
	blk := fn.Blocks[0]
	if len(blk.Instrs) > 10 {
		return false // too expensive
	}
	if fn.Synthetic != "" && (fn.Pkg == nil || fn != fn.Pkg.Func("init")) {
		return true // treat synthetic wrappers context-sensitively
	}
	return true
}

func (c *kCallsiteContext) String() string {
	if len(c.instr) != 0 {
		return c.instr[len(c.instr)-1].Common().Description()
	}
	return "synthetic function call"
}

func (c *kCallsiteContext) Pos() token.Pos {
	if len(c.instr) != 0 {
		return c.instr[len(c.instr)-1].Pos()
	}
	return token.NoPos
}

func (c *kCallsiteContext) HashString(fn *ssa.Function) string {
	str := fn.String()
	for _, ins := range c.instr {
		str = str + ins.String()
	}
	return str
}

func EmptyContext() *kCallsiteContext {
	return &kCallsiteContext{}
}

func (c *kCallsiteContext) NewContext(site ssa.CallInstruction) context {
	k := 2
	instrs := append(c.instr, site)
	if len(instrs) > k {
		instrs = instrs[1:]
	}
	return &kCallsiteContext{instr: instrs}
}

func (c *kCallsiteContext) IsEmpty() bool {
	return len(c.instr) == 0
}
