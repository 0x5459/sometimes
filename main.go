package main

import (
	"fmt"
	"sometimes/hir"
	"sometimes/vm"
	"sometimes/vm/assembly"
)

func main() {

	hirBuilder := hir.NewBuilder()
	hirBuilder.InsertFunc(funcMain(), true)
	hirBuilder.InsertFunc(funcTest(), false)
	hirBuilder.InsertFunc(funcSum(), false)
	asmCompiler := assembly.NewCompiler(hirBuilder.Build())
	asm := asmCompiler.Compile()
	fmt.Println(asm.String())

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
		Args:   []hir.Expr{&hir.ExprLiteral{Val: hir.NewValueFloat(3.14159)}, varC},
	})

	hirFuncBuilder.Emit(&hir.ExprCall{
		Callee: &hir.ExprVar{VarBinding: hir.NewBinding("sum")},
		Args:   []hir.Expr{&hir.ExprLiteral{Val: hir.NewValueInt(100)}},
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

	hirFuncBuilder.Emit(&hir.ExprReturn{Expr: &hir.ExprIf{
		Cond: &hir.ExprBinary{
			Lhs: varA,
			Rhs: &hir.ExprLiteral{Val: hir.NewValueInt(37)},
			Op:  hir.OpGT,
		},
		Body: &hir.ExprBlock{
			Body: []hir.Expr{&hir.ExprBinary{
				Lhs: varA,
				Rhs: varB,
				Op:  hir.OpAdd,
			}},
		},
		Else: &hir.ExprBlock{
			Body: []hir.Expr{&hir.ExprBinary{
				Lhs: varA,
				Rhs: varB,
				Op:  hir.OpMul,
			}},
		},
	},
	})
	return hirFuncBuilder.Build()
}

func funcSum() *hir.ExprFunction {

	varN := &hir.ExprVar{VarBinding: hir.NewBinding("n")}
	varS := &hir.ExprVar{VarBinding: hir.NewBinding("s")}
	varI := &hir.ExprVar{VarBinding: hir.NewBinding("i")}
	hirFuncBuilder := hir.NewFuncBuilder("sum", []*hir.Binding{varN.VarBinding})
	// s = 0
	hirFuncBuilder.Emit(&hir.ExprMutate{
		Lhs: varS,
		Rhs: &hir.ExprLiteral{Val: hir.NewValueInt(0)},
	})

	// i = 1
	hirFuncBuilder.Emit(&hir.ExprMutate{
		Lhs: varI,
		Rhs: &hir.ExprLiteral{Val: hir.NewValueInt(1)},
	})

	hirFuncBuilder.Emit(&hir.ExprLoop{
		Cond: &hir.ExprBinary{
			Lhs: varI,
			Rhs: varN,
			Op:  hir.OpLTE,
		},
		Body: &hir.ExprBlock{
			Body: []hir.Expr{
				// s+=i
				&hir.ExprMutate{
					Lhs: varS,
					Rhs: &hir.ExprBinary{
						Lhs: varS,
						Rhs: varI,
						Op:  hir.OpAdd,
					},
				},
				// i++
				&hir.ExprMutate{
					Lhs: varI,
					Rhs: &hir.ExprBinary{
						Lhs: varI,
						Rhs: &hir.ExprLiteral{Val: hir.NewValueInt(1)},
						Op:  hir.OpAdd,
					},
				},
			},
		},
	})

	hirFuncBuilder.Emit(&hir.ExprReturn{Expr: varS})
	return hirFuncBuilder.Build()
}
