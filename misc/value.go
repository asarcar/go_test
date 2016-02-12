package misc

import (
	"fmt"
)

type Value interface {
	String() string
	Op(op string, y Value) Value
}

type Bool bool

func (x Bool) String() string { return fmt.Sprintf("%v", bool(x)) }
func (x Bool) Op(op string, y Value) Value {
	return Error(fmt.Sprintf("illegal operator: '%s %s %s'", x, op, y))
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
		case "==":
			return Bool(x == y)
		case "!=":
			return Bool(x != y)
		case "<":
			return Bool(x < y)
		case ">":
			return Bool(x > y)
		case ">=":
			return Bool(x >= y)
		case "<=":
			return Bool(x <= y)
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
