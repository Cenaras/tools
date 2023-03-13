package pointer

import (
	"go/token"
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

type KObjNHeap struct {
	K int
	N int
}

type KObjNHeapContext struct {
	allocs []ssa.Value
}

func (c *KObjNHeapContext) String() string {
	str := ""
	for _, all := range c.allocs {
		str = str + all.String() + strconv.Itoa(int(all.Pos()))
	}
	return str
}

func (cs *KObjNHeap) Record(value ssa.Value, ctx Context) HeapContext {
	allocs := ctx.(*KObjNHeapContext).allocs
	if len(allocs)-cs.N <= 0 {
		return ctx
	} else {
		return &KObjNHeapContext{allocs: allocs[len(allocs)-cs.N:]}
	}
}
func (cs *KObjNHeap) Merge(value ssa.Value, hctx HeapContext, callLabel ssa.CallInstruction, ctx Context) Context {
	allocs := append(hctx.(*KObjNHeapContext).allocs, value)
	if len(allocs) > cs.K {
		allocs = allocs[1:]
	}
	return &KObjNHeapContext{allocs: allocs}
}
func (cs *KObjNHeap) MergeStatic(callLabel ssa.CallInstruction, ctx Context) Context {
	return ctx
}
func (cs *KObjNHeap) EmptyContext() Context {
	return &KObjNHeapContext{}
}
func (cs *KObjNHeap) EmptyHeapContext() HeapContext {
	return &KObjNHeapContext{}
}

func (cs *KObjNHeap) TreatStaticInvoke() bool {
	return true
}

// Uniform hybrid context strategies
type U1Obj struct {
}

type U1ObjContext struct {
	heap ssa.Value
	invo ssa.CallInstruction
}

func (c *U1ObjContext) String() string {
	str := ""
	if c.heap != nil {
		str = str + c.heap.String() + strconv.Itoa(int(c.heap.Pos()))
	}
	if c.invo != nil {
		str = str + c.invo.String() + strconv.Itoa(int(c.invo.Pos()))
	}
	return str
}

func (cs *U1Obj) Record(value ssa.Value, ctx Context) HeapContext {
	return cs.EmptyHeapContext()
}
func (cs *U1Obj) Merge(value ssa.Value, hctx HeapContext, callLabel ssa.CallInstruction, ctx Context) Context {
	return &U1ObjContext{heap: value, invo: callLabel}
}
func (cs *U1Obj) MergeStatic(callLabel ssa.CallInstruction, ctx Context) Context {
	return &U1ObjContext{heap: ctx.(*U1ObjContext).heap, invo: callLabel}
}
func (cs *U1Obj) EmptyContext() Context {
	return &U1ObjContext{}
}
func (cs *U1Obj) EmptyHeapContext() HeapContext {
	return &U1ObjContext{}
}

type U2ObjH struct {
}

type U2ObjHContext struct {
	heap  ssa.Value
	heap2 ssa.Value
	invo  ssa.CallInstruction
}

type U2ObjHHeapContext struct {
	heap ssa.Value
}

func (c *U2ObjHContext) String() string {
	str := ""
	if c.heap != nil {
		str = str + c.heap.String() + strconv.Itoa(int(c.heap.Pos()))
	}
	if c.heap2 != nil {
		str = str + c.heap2.String() + strconv.Itoa(int(c.heap2.Pos()))
	}
	if c.invo != nil {
		str = str + c.invo.String() + strconv.Itoa(int(c.invo.Pos()))
	}
	return str
}

func (c *U2ObjHHeapContext) String() string {
	str := ""
	if c.heap != nil {
		str = str + c.heap.String() + strconv.Itoa(int(c.heap.Pos()))
	}
	return str
}

func (cs *U2ObjH) Record(value ssa.Value, ctx Context) HeapContext {
	return &U2ObjHHeapContext{heap: ctx.(*U2ObjHContext).heap}
}
func (cs *U2ObjH) Merge(value ssa.Value, hctx HeapContext, callLabel ssa.CallInstruction, ctx Context) Context {
	return &U2ObjHContext{heap: value, heap2: hctx.(*U2ObjHHeapContext).heap, invo: callLabel}
}
func (cs *U2ObjH) MergeStatic(callLabel ssa.CallInstruction, ctx Context) Context {
	return &U2ObjHContext{heap: ctx.(*U2ObjHContext).heap, heap2: ctx.(*U2ObjHContext).heap2, invo: callLabel}
}
func (cs *U2ObjH) EmptyContext() Context {
	return &U2ObjHContext{}
}
func (cs *U2ObjH) EmptyHeapContext() HeapContext {
	return &U2ObjHHeapContext{}
}

// Selective hybrid context strategies

type SB1Obj struct {
}

type SB1ObjContext struct {
	heap ssa.Value
	invo ssa.CallInstruction
}

func (c *SB1ObjContext) String() string {
	str := ""
	if c.heap != nil {
		str = str + c.heap.String() + strconv.Itoa(int(c.heap.Pos()))
	}
	if c.invo != nil {
		str = str + c.invo.String() + strconv.Itoa(int(c.invo.Pos()))
	}
	return str
}

func (cs *SB1Obj) Record(value ssa.Value, ctx Context) HeapContext {
	return cs.EmptyHeapContext()
}
func (cs *SB1Obj) Merge(value ssa.Value, hctx HeapContext, callLabel ssa.CallInstruction, ctx Context) Context {
	return &SB1ObjContext{heap: value, invo: nil}
}
func (cs *SB1Obj) MergeStatic(callLabel ssa.CallInstruction, ctx Context) Context {
	return &SB1ObjContext{heap: ctx.(*SB1ObjContext).heap, invo: callLabel}
}
func (cs *SB1Obj) EmptyContext() Context {
	return &SB1ObjContext{}
}
func (cs *SB1Obj) EmptyHeapContext() HeapContext {
	return &SB1ObjContext{}
}

type SA1Obj struct {
}

type SA1ObjContext struct {
	arg ContextArg
}

func (c *SA1ObjContext) String() string {
	str := ""
	if c.arg != nil {
		str = str + c.arg.String() + strconv.Itoa(int(c.arg.Pos()))
	}
	return str
}

func (cs *SA1Obj) Record(value ssa.Value, ctx Context) HeapContext {
	return cs.EmptyHeapContext()
}
func (cs *SA1Obj) Merge(value ssa.Value, hctx HeapContext, callLabel ssa.CallInstruction, ctx Context) Context {
	return &SA1ObjContext{arg: value}
}
func (cs *SA1Obj) MergeStatic(callLabel ssa.CallInstruction, ctx Context) Context {
	return &SA1ObjContext{arg: callLabel}
}
func (cs *SA1Obj) EmptyContext() Context {
	return &SA1ObjContext{}
}
func (cs *SA1Obj) EmptyHeapContext() HeapContext {
	return &SA1ObjContext{}
}

type ContextArg interface {
	String() string
	Pos() token.Pos
}

type S2ObjH struct {
}

type S2ObjHContext struct {
	heap ssa.Value
	arg2 ContextArg
	arg3 ContextArg
}

type S2ObjHHeapContext struct {
	heap ssa.Value
}

func (c *S2ObjHContext) String() string {
	str := ""
	if c.heap != nil {
		str = str + c.heap.String() + strconv.Itoa(int(c.heap.Pos()))
	}
	if c.arg2 != nil {
		str = str + c.arg2.String() + strconv.Itoa(int(c.arg2.Pos()))
	}
	if c.arg3 != nil {
		str = str + c.arg3.String() + strconv.Itoa(int(c.arg3.Pos()))
	}
	return str
}

func (c *S2ObjHHeapContext) String() string {
	str := ""
	if c.heap != nil {
		str = str + c.heap.String() + strconv.Itoa(int(c.heap.Pos()))
	}
	return str
}

func (cs *S2ObjH) Record(value ssa.Value, ctx Context) HeapContext {
	return &S2ObjHHeapContext{heap: ctx.(*S2ObjHContext).heap}
}
func (cs *S2ObjH) Merge(value ssa.Value, hctx HeapContext, callLabel ssa.CallInstruction, ctx Context) Context {
	return &S2ObjHContext{heap: value, arg2: hctx.(*S2ObjHHeapContext).heap, arg3: nil}
}
func (cs *S2ObjH) MergeStatic(callLabel ssa.CallInstruction, ctx Context) Context {
	return &S2ObjHContext{heap: ctx.(*S2ObjHContext).heap, arg2: callLabel, arg3: ctx.(*S2ObjHContext).arg2}
}
func (cs *S2ObjH) EmptyContext() Context {
	return &S2ObjHContext{}
}
func (cs *S2ObjH) EmptyHeapContext() HeapContext {
	return &S2ObjHHeapContext{}
}
