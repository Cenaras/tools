package pointer

import (
	"strconv"

	"golang.org/x/tools/go/ssa"
)

type Insens struct {
}

type InsensContext struct{}

func (c *InsensContext) String() string {
	return ""
}

func (cs *Insens) Record(value ssa.Value, ctx Context) HeapContext {
	return &InsensContext{}
}
func (cs *Insens) Merge(value ssa.Value, hctx HeapContext, callLabel ssa.CallInstruction, ctx Context) Context {
	return &InsensContext{}
}
func (cs *Insens) MergeStatic(callLabel ssa.CallInstruction, ctx Context) Context {
	return &InsensContext{}
}
func (cs *Insens) EmptyContext() Context {
	return &InsensContext{}
}
func (cs *Insens) EmptyHeapContext() HeapContext {
	return &InsensContext{}
}

type KCallNHeap struct {
	K int
	N int
}

type KCallNHeapContext struct {
	instrs []ssa.CallInstruction
}

func (c *KCallNHeapContext) String() string {
	str := ""
	for _, ins := range c.instrs {
		str = str + ins.String() + strconv.Itoa(int(ins.Pos()))
	}
	return str
}

func (cs *KCallNHeap) Record(value ssa.Value, ctx Context) HeapContext {
	instrs := ctx.(*KCallNHeapContext).instrs
	if len(instrs)-cs.N <= 0 {
		return ctx
	} else {
		return &KCallNHeapContext{instrs: instrs[len(instrs)-cs.N:]}
	}
}
func (cs *KCallNHeap) Merge(value ssa.Value, hctx HeapContext, callLabel ssa.CallInstruction, ctx Context) Context {
	return cs.MergeStatic(callLabel, ctx)
}
func (cs *KCallNHeap) MergeStatic(callLabel ssa.CallInstruction, ctx Context) Context {
	instrs := append(ctx.(*KCallNHeapContext).instrs, callLabel)
	if len(instrs) > cs.K {
		instrs = instrs[1:]
	}
	return &KCallNHeapContext{instrs: instrs}
}
func (cs *KCallNHeap) EmptyContext() Context {
	return &KCallNHeapContext{}
}
func (cs *KCallNHeap) EmptyHeapContext() HeapContext {
	return &KCallNHeapContext{}
}
