package pointer

import (
	"strconv"

	"golang.org/x/tools/go/ssa"
)

type ContextStrategy interface {
	Record(value ssa.Value, ctx Context) HeapContext
	Merge(value ssa.Value, hctx HeapContext, callLabel ssa.CallInstruction, ctx Context) Context
	MergeStatic(callLabel ssa.CallInstruction, ctx Context) Context
	EmptyContext() Context
	EmptyHeapContext() HeapContext
	TreatStaticInvoke() bool
	ShouldUseContext(fn *ssa.Function, a *analysis) bool
}

type Context interface {
	String() string
}
type HeapContext interface {
	String() string
}

type DefaultContext struct {
	instr ssa.CallInstruction
}

func (c *DefaultContext) String() string {
	if c.instr != nil {
		return c.instr.String() + strconv.Itoa(int(c.instr.Pos()))
	} else {
		return ""
	}
}

type defaultHeapContext struct {
}

func (c *defaultHeapContext) String() string {
	return ""
}

type defaultContextStrategy struct {
}

func (cs *defaultContextStrategy) EmptyContext() Context {
	return &DefaultContext{}
}

func (cs *defaultContextStrategy) EmptyHeapContext() HeapContext {
	return &defaultHeapContext{}
}

func (cs *defaultContextStrategy) Record(obj ssa.Value, ctx Context) HeapContext {
	return ctx
}

func (cs *defaultContextStrategy) Merge(obj ssa.Value, hctx HeapContext, callLabel ssa.CallInstruction, ctx Context) Context {
	return cs.EmptyContext()
}

func (cs *defaultContextStrategy) MergeStatic(callLabel ssa.CallInstruction, ctx Context) Context {
	return &DefaultContext{instr: callLabel}
}

func (cs *defaultContextStrategy) TreatStaticInvoke() bool {
	return false
}

func (cs *defaultContextStrategy) ShouldUseContext(fn *ssa.Function, a *analysis) bool {
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
	for _, instr := range blk.Instrs {
		switch instr := instr.(type) {
		case ssa.CallInstruction:
			// Disallow function calls (except to built-ins)
			// because of the danger of unbounded recursion.
			if _, ok := instr.Common().Value.(*ssa.Builtin); !ok {
				return false
			}
		}
	}
	return true
}
