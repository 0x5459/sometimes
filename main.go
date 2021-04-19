package main

import (
	"sometimes/hir"
	"sometimes/vm"
	"sometimes/vm/assembly"
)

func main() {
	hirFuncBuilder := hir.NewFuncBuilder("main", []*hir.Binding{})
	varA := &hir.ExprVar{
		VarBinding: hir.NewBinding("a"),
	}
	varB := &hir.ExprVar{
		VarBinding: hir.NewBinding("b"),
	}
	varC := &hir.ExprVar{
		VarBinding: hir.NewBinding("c"),
	}

	hirFuncBuilder.Emit(&hir.ExprMutate{
		Lhs: varA,
		Rhs: &hir.ExprBinary{
			Lhs: &hir.ExprLiteral{Val: hir.NewValueInt(10)},
			Rhs: &hir.ExprLiteral{Val: hir.NewValueInt(20)},
			Op:  hir.OpAdd,
		},
	})

	hirFuncBuilder.Emit(&hir.ExprMutate{
		Lhs: varB,
		Rhs: &hir.ExprBinary{
			Lhs: &hir.ExprLiteral{Val: hir.NewValueInt(2)},
			Rhs: &hir.ExprLiteral{Val: hir.NewValueInt(3)},
			Op:  hir.OpMul,
		},
	})
	hirFuncBuilder.Emit(&hir.ExprMutate{
		Lhs: varC,
		Rhs: &hir.ExprBinary{
			Lhs: varA,
			Rhs: varB,
			Op:  hir.OpAdd,
		},
	})

	hirBuilder := hir.NewBuilder()
	hirBuilder.InsertFunc(hirFuncBuilder.Build(), true)
	asmCompiler := assembly.NewCompiler(hirBuilder.Build())
	asm := asmCompiler.Compile()
	program := vm.NewProgramFromAsm(asm)
	machine := vm.New(program, 256, 128)
	machine.Execute()
}
