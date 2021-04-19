package main

import (
	"sometimes/hir"
	"sometimes/vm"
	"sometimes/vm/assembly"
)

func main() {

	hirBuilder := hir.NewBuilder()
	hirBuilder.InsertFunc(funcMain(), true)
	hirBuilder.InsertFunc(funcTest(), false)
	asmCompiler := assembly.NewCompiler(hirBuilder.Build())
	asm := asmCompiler.Compile()
	program := vm.NewProgramFromAsm(asm)
	machine := vm.New(program, 256, 128)
	machine.Execute()
}

func funcMain() *hir.ExprFunction {
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

	hirFuncBuilder.Emit(&hir.ExprCall{
		Callee: &hir.ExprVar{VarBinding: hir.NewBinding("test")},
		Args:   []hir.Expr{varB, varC},
	})
	hirFuncBuilder.Emit(&hir.ExprReturn{})
	return hirFuncBuilder.Build()
}

func funcTest() *hir.ExprFunction {
	varA := &hir.ExprVar{
		VarBinding: hir.NewBinding("a"),
	}
	varB := &hir.ExprVar{
		VarBinding: hir.NewBinding("b"),
	}
	hirFuncBuilder := hir.NewFuncBuilder("test", []*hir.Binding{varA.VarBinding, varB.VarBinding})

	varC := &hir.ExprVar{
		VarBinding: hir.NewBinding("c"),
	}

	hirFuncBuilder.Emit(&hir.ExprIf{
		Cond: &hir.ExprBinary{
			Lhs: varA,
			Rhs: &hir.ExprLiteral{Val: hir.NewValueInt(100)},
			Op:  hir.OpGT,
		},
		Body: &hir.ExprBlock{
			Body: []hir.Expr{&hir.ExprMutate{
				Lhs: varC,
				Rhs: &hir.ExprBinary{
					Lhs: varA,
					Rhs: varB,
					Op:  hir.OpAdd,
				},
			}},
		},
		Else: &hir.ExprBlock{
			Body: []hir.Expr{&hir.ExprMutate{
				Lhs: varC,
				Rhs: &hir.ExprBinary{
					Lhs: varA,
					Rhs: varB,
					Op:  hir.OpMul,
				},
			}},
		},
	})
	hirFuncBuilder.Emit(&hir.ExprReturn{Expr: varC})
	return hirFuncBuilder.Build()
}
