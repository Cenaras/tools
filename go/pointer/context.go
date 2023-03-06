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
}

type Context interface {
	String() string
}
type HeapContext interface {
	String() string
}

type defaultContext struct {
	instr ssa.CallInstruction
}

func (c *defaultContext) String() string {
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
	return &defaultContext{}
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
	return &defaultContext{instr: callLabel}
}

func (cs *defaultContextStrategy) TreatStaticInvoke() bool {
	return false
}
