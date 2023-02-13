// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pointer

// This file defines the internal (context-sensitive) call graph.

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/ssa"
)

type cgnode struct {
	fn            *ssa.Function
	obj           nodeid    // start of this contour's object block
	sites         []context // ordered list of callsites within this function
	callercontext context   // where called from, if known; nil for shared contours
}

// contour returns a description of this node's contour.
func (n *cgnode) contour() string {
	if n.callercontext == nil {
		return "shared contour"
	}
	if true {
		return fmt.Sprintf("as called from %s", "test")
	}
	return fmt.Sprintf("as called from intrinsic (targets=n%d)", n.callercontext.Targets())
}

func (n *cgnode) String() string {
	return fmt.Sprintf("cg%d:%s", n.obj, n.fn)
}

// A callsite represents a single call site within a cgnode;
// it is implicitly context-sensitive.
// callsites never represent calls to built-ins;
// they are handled as intrinsics.
type callsite struct {
	targets nodeid              // pts(·) contains objects for dynamically called functions
	instr   ssa.CallInstruction // the call instruction; nil for synthetic/intrinsic
}

func (c *callsite) String() string {
	if c.instr != nil {
		return c.instr.Common().Description()
	}
	return "synthetic function call"
}

// pos returns the source position of this callsite, or token.NoPos if implicit.
func (c *callsite) pos() token.Pos {
	if c.instr != nil {
		return c.instr.Pos()
	}
	return token.NoPos
}

type context interface {
	Targets() nodeid
	SetTargets(targets nodeid)
	Instr() ssa.CallInstruction
	NewContext(site ssa.CallInstruction) context
	ShouldUseContext(fn *ssa.Function) bool
	String() string
	Pos() token.Pos
}

type callsiteContext struct {
	targets nodeid
	instr   ssa.CallInstruction
}

func (c *callsiteContext) Targets() nodeid {
	return c.targets
}

func (c *callsiteContext) SetTargets(targets nodeid) {
	c.targets = targets
}

func (c *callsiteContext) Instr() ssa.CallInstruction {
	return c.instr
}

func (c *callsiteContext) ShouldUseContext(fn *ssa.Function) bool {
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

func (c *callsiteContext) String() string {
	return "CallSite Instruction TODO"
}

func (c *callsiteContext) Pos() token.Pos {
	if c.instr != nil {
		return c.instr.Pos()
	}
	return token.NoPos
}

func EmptyContext() *callsiteContext {
	return &callsiteContext{}
}

func (c *callsiteContext) NewContext(site ssa.CallInstruction) context {
	return &callsiteContext{instr: site}
}
