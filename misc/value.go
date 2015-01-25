package misc

import (
	"fmt"
)

type Value interface {
	String() string
	Op(op string, y Value) Value
}

type Int int

func (x Int) String() string { return fmt.Sprintf("%d", int(x)) }
func (x Int) Op(op string, y Value) Value {
	switch y := y.(type) { // type switch
	case Int:
		switch op {
		case "+":
			return x + y
		case "-":
			return x - y
		case "*":
			return x * y
		case "/":
			return x / y
		default:
			return Error(fmt.Sprintf("illegal operator: '%s %s %s'", x, op, y))
		}
	case Error:
		return y
	}
	return Error(fmt.Sprintf("illegal operand: '%s %s %s'", x, op, y))
}

// Error is a 'class' that 'implements' Value
type Error string

func (e Error) String() string              { return string(e) }
func (e Error) Op(op string, y Value) Value { return e }

type Str string

func (s Str) String() string { return string(s) }
func (s Str) Op(op string, y Value) Value {
	switch y := y.(type) { // type switch
	case Str:
		switch op {
		case "+":
			return s + y
		default:
			return Error(fmt.Sprintf("illegal operator: '%s %s %s'", s, op, y))
		}
	case Int:
		switch op {
		case "*":
			x := s
			for i := y - 1; i > 0; i-- {
				x += s
			}
			return x
		default:
			return Error(fmt.Sprintf("illegal operator: '%s %s %s'", s, op, y))
		}
	case Error:
		return y
	}
	return Error(fmt.Sprintf("illegal operand: '%s %s %s'", s, op, y))
}

func ValueTest1() {
	var x1, x2, x3, x4, x5 Int = 10, 5, -8, 6, -3

	z1 := x1.Op("+", x2.Op("*", x3).Op("-", x4.Op("/", x5)))

	fmt.Println("x1 = 10, x2 = 5, x3 = -8, x4 = 6, x5 = -3: x1 + x2*x3 - x4/x5 =", z1)
}

func ValueTest2() {
	var x1, x2 Str = "ra", "ja"
	var x3, x4, x5 Int = 3, 6, 3

	z1 := x1.Op("+", x2.Op("*", x3).Op("*", x4.Op("/", x5)))

	fmt.Println("x1 = 'ra', x2 = 'ja', x3 = 3, x4 = 6, x5 = 3: x1 + x2*x3*(x4/x5) =", z1)
}

func DumpValue() {
	fmt.Println("Value Test\n----------------")
	ValueTest1()
	ValueTest2()
	fmt.Println("----------------")
}
