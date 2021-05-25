package main

import (
	"sometimes/lexer"
	"sometimes/parser"
	"sometimes/visitor"
	"sometimes/vm"
	"sometimes/vm/assembly"
)

func main() {
	code :=
		`
// 计算第 n 项斐波拉契数列小程序
fn main() {
	let a=1, b=1, t=0;
	let i=2, n=7;
	loop (i<n) {
		t = a;
		a = b;
		b = t + a;
		i += 1;
	};
	print(b);
}
`
	parser := parser.NewParser(lexer.NewTokenCursor(lexer.NewSrcCursor([]byte(code))))
	vis := visitor.NewVistor()
	prog := vis.Visit(parser.Parse())

	asmCompiler := assembly.NewCompiler(prog)
	asm := asmCompiler.Compile()
	// fmt.Println(asm.String())

	program := vm.NewProgramFromAsm(asm)
	machine := vm.New(program, 256, 128)
	machine.Execute()
}
