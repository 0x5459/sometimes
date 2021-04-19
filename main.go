package main

import (
	"sometimes/hir"
	"sometimes/vm"
	"sometimes/vm/assembly"
)

func main() {
	hirBuilder := hir.NewFuncBuilder("main", []*hir.Binding{})
	varA := &hir.ExprVar{
		VarBinding: hir.NewBinding("a"),
	}
	varB := &hir.ExprVar{
		VarBinding: hir.NewBinding("b"),
	}
	varC := &hir.ExprVar{
		VarBinding: hir.NewBinding("c"),
	}

	hirBuilder.Emit(&hir.ExprMutate{
		Lhs: varA,
		Rhs: &hir.ExprBinary{
			Lhs: &hir.ExprLiteral{Val: hir.NewValueInt(10)},
			Rhs: &hir.ExprLiteral{Val: hir.NewValueInt(20)},
			Op:  hir.OpAdd,
		},
	})

	hirBuilder.Emit(&hir.ExprMutate{
		Lhs: varB,
		Rhs: &hir.ExprBinary{
			Lhs: &hir.ExprLiteral{Val: hir.NewValueInt(2)},
			Rhs: &hir.ExprLiteral{Val: hir.NewValueInt(3)},
			Op:  hir.OpMul,
		},
	})
	hirBuilder.Emit(&hir.ExprMutate{
		Lhs: varC,
		Rhs: &hir.ExprBinary{
			Lhs: varA,
			Rhs: varB,
			Op:  hir.OpAdd,
		},
	})
	hirProgram := hirBuilder.Build()
	asmCompiler := assembly.NewCompiler(hirProgram)
	asm := asmCompiler.Compile()
	program := vm.NewProgramFromAsm(asm)
	machine := vm.New(program, 256, 128)
	machine.Execute()
}
