package vm

import (
	"fmt"
	"math"
	"sometimes/vm/value"
)

type arithOperator struct {
	intFunc   func(int, int) int
	floatFunc func(float64, float64) float64
}

var arithOperators = []arithOperator{
	OpAdd - op_arith_start: arithOperator{intFunc: _iadd, floatFunc: _fadd},
	OpSub - op_arith_start: arithOperator{intFunc: _isub, floatFunc: _fsub},
	OpMul - op_arith_start: arithOperator{intFunc: _imul, floatFunc: _fmul},
	OpDiv - op_arith_start: arithOperator{intFunc: _idiv, floatFunc: _fdiv},
	OpMod - op_arith_start: arithOperator{intFunc: _imod, floatFunc: _fmod},
	OpNeg - op_arith_start: arithOperator{intFunc: _ineg, floatFunc: _fneg},
}

type ArithInstruction interface {
	Instruction
	isArith()
}

type (
	BinaryArithInstruction interface {
		ArithInstruction
		isBinaryArith()
	}
	UnaryArithInstruction interface {
		ArithInstruction
		isUnaryArith()
	}
)

func (*InstrAdd) isArith() {}
func (*InstrSub) isArith() {}
func (*InstrMul) isArith() {}
func (*InstrDiv) isArith() {}
func (*InstrMod) isArith() {}
func (*InstrNeg) isArith() {}

func (*InstrAdd) isBinaryArith() {}
func (*InstrSub) isBinaryArith() {}
func (*InstrMul) isBinaryArith() {}
func (*InstrDiv) isBinaryArith() {}
func (*InstrMod) isBinaryArith() {}
func (*InstrNeg) isUnaryArith()  {}

func arith(op ArithInstruction, x, y value.Value) value.Value {
	if _, ok := x.(value.NumberValue); !ok || x.Type() != y.Type() {
		panic(unsupportedOperandError(op.Op(), x, y))
	}

	f := arithOperators[op.Op()-op_arith_start]

	switch a := x.(type) {
	case (*value.Int):
		b := y.(*value.Int)
		return &value.Int{Val: f.intFunc(a.Val, b.Val)}
	case (*value.Float):
		b := y.(*value.Float)
		return &value.Float{Val: f.floatFunc(a.Val, b.Val)}
	}

	// never
	return &value.Nil{}
}

func _iadd(x, y int) int         { return x + y }
func _fadd(x, y float64) float64 { return x + y }

func _isub(x, y int) int         { return x - y }
func _fsub(x, y float64) float64 { return x - y }

func _imul(x, y int) int         { return x * y }
func _fmul(x, y float64) float64 { return x * y }

func _idiv(x, y int) int         { return x / y }
func _fdiv(x, y float64) float64 { return x / y }

func _imod(x, y int) int { return x % y }

var _fmod = math.Mod

func _ineg(x, _ int) int         { return -x }
func _fneg(x, _ float64) float64 { return -x }

var logicOperators = []func(x, y value.Value) bool{
	OpEq - op_logic_start:  _eq,
	OpNE - op_logic_start:  _ne,
	OpGT - op_logic_start:  _gt,
	OpLT - op_logic_start:  _lt,
	OpGTE - op_logic_start: _gte,
	OpLTE - op_logic_start: _lte,
	OpNot - op_logic_start: _not,
	OpAnd - op_logic_start: _and,
	OpOr - op_logic_start:  _or,
}

type LogicInstruction interface {
	Instruction
	isLogic()
}

type (
	BinaryLogicInstruction interface {
		LogicInstruction
		isBinaryLogic()
	}
	UnaryLogicInstruction interface {
		LogicInstruction
		isUnaryLogic()
	}
)

func (*InstrEq) isLogic()  {}
func (*InstrNE) isLogic()  {}
func (*InstrGT) isLogic()  {}
func (*InstrLT) isLogic()  {}
func (*InstrGTE) isLogic() {}
func (*InstrLTE) isLogic() {}
func (*InstrAnd) isLogic() {}
func (*InstrOr) isLogic()  {}
func (*InstrNot) isLogic() {}

func (*InstrEq) isBinaryLogic()  {}
func (*InstrNE) isBinaryLogic()  {}
func (*InstrGT) isBinaryLogic()  {}
func (*InstrLT) isBinaryLogic()  {}
func (*InstrGTE) isBinaryLogic() {}
func (*InstrLTE) isBinaryLogic() {}
func (*InstrAnd) isBinaryLogic() {}
func (*InstrOr) isBinaryLogic()  {}
func (*InstrNot) isUnaryLogic()  {}

func logic(op LogicInstruction, x, y value.Value) bool {
	f := logicOperators[op.Op()]
	return f(x, y)
}

func _eq(x, y value.Value) bool {
	switch a := x.(type) {
	case (*value.Int):
	case (*value.Float):
	case (*value.Char):
		switch b := y.(type) {
		case (*value.Int):
		case (*value.Float):
		case (*value.Char):
			return a.Val == b.Val
		}
	case (*value.Boolean):
		if b, ok := y.(*value.Boolean); ok {
			return a.Val == b.Val
		}
	case (*value.Nil):
		_, ok := y.(*value.Nil)
		return ok
	}
	panic(unsupportedOperandError(OpEq, x, y))
}

func _ne(x, y value.Value) bool {
	return !_eq(x, y)
}

func _gt(x, y value.Value) bool {
	switch a := x.(type) {
	case (*value.Int):
	case (*value.Float):
	case (*value.Char):
		switch b := y.(type) {
		case (*value.Int):
		case (*value.Float):
		case (*value.Char):
			return a.Val > b.Val
		}
	}
	panic(unsupportedOperandError(OpGT, x, y))
}

func _lt(x, y value.Value) bool {
	switch a := x.(type) {
	case (*value.Int):
	case (*value.Float):
	case (*value.Char):
		switch b := y.(type) {
		case (*value.Int):
		case (*value.Float):
		case (*value.Char):
			return a.Val < b.Val
		}
	}
	panic(unsupportedOperandError(OpLT, x, y))
}

func _gte(x, y value.Value) bool {
	switch a := x.(type) {
	case (*value.Int):
	case (*value.Float):
	case (*value.Char):
		switch b := y.(type) {
		case (*value.Int):
		case (*value.Float):
		case (*value.Char):
			return a.Val >= b.Val
		}
	}
	panic(unsupportedOperandError(OpGTE, x, y))
}

func _lte(x, y value.Value) bool {
	switch a := x.(type) {
	case (*value.Int):
	case (*value.Float):
	case (*value.Char):
		switch b := y.(type) {
		case (*value.Int):
		case (*value.Float):
		case (*value.Char):
			return a.Val <= b.Val
		}
	}
	panic(unsupportedOperandError(OpLTE, x, y))
}

func _not(x, _ value.Value) bool {
	if a, ok := x.(*value.Boolean); ok {
		return !a.Val
	}
	panic(unsupportedOperandError(OpEq, x, &value.Nil{}))
}

func _and(x, y value.Value) bool {
	if a, xIsBool := x.(*value.Boolean); xIsBool {
		if b, bIsBool := y.(*value.Boolean); bIsBool {
			return a.Val && b.Val
		}
	}
	panic(unsupportedOperandError(OpAnd, x, y))
}

func _or(x, y value.Value) bool {
	if a, xIsBool := x.(*value.Boolean); xIsBool {
		if b, bIsBool := y.(*value.Boolean); bIsBool {
			return a.Val && b.Val
		}
	}
	panic(unsupportedOperandError(OpOr, x, y))
}

func unsupportedOperandError(op Op, lhs, rhs value.Value) error {
	panic(fmt.Errorf("unsupported operand type for `%s`: lhs: `%s` rhs: `%s`",
		op.String(), lhs.Type().String(), rhs.Type().String()))
}
