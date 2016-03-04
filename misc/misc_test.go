package misc

import "testing"

func TestIntValue(t *testing.T) {
	const in1, in2, in3, in4, in5, out Int = 10, 5, -8, 6, -3, -28
	if x := in1.Op("+", in2.Op("*", in3).Op("-", in4.Op("/", in5))); x != out {
		t.Errorf("in1 = %d, in2 = %d, in3 = %d, in4 = %d, in5 = %d: ",
			in1, in2, in3, in4, in5)
		t.Errorf("in1 + in2*in3 - in4/in5 [=%d] != %d",
			in1, in2, in3, in4, in5, x, out)
	}
}

func TestStrValue(t *testing.T) {
	const in1, in2 Str = "ra", "ja"
	const in3, in4, in5 Int = 3, 6, 3
	const out Str = "rajajajajajaja"
	if x := in1.Op("+", in2.Op("*", in3).Op("*", in4.Op("/", in5))); x != out {
		t.Errorf("in1 = %s, in2 = %s, in3 = %d, in4 = %d, in5 = %d: ",
			in1, in2, in3, in4, in5)
		t.Errorf("in1 + in2*in3*(in4/in5) [=%s] != %s", x, out)
	}
}

func TestCompose(t *testing.T) {
	const in, out1, out2 = 1, 1, 3
	var sqr Fn = func(x int) int { return x * x }
	var inc Fn = func(x int) int { return x + 1 }
	var dec Fn = func(x int) int { return x - 1 }

	x1 := Compose(inc, sqr)(dec)(in)
	x2 := Compose(dec, sqr)(inc)(in)

	if x1 != out1 || x2 != out2 {
		t.Errorf("inc(sqr(dec(%d))) [%d] != %d || dec(sqr(inc(%d))) [%d] != %d",
			in, x1, out1, in, x2, out2)
	}
}

func TestComposeFn(t *testing.T) {
	const in, out1, out2 = 1, 1, 3
	var sqr Fn = func(x int) int { return x * x }
	var inc Fn = func(x int) int { return x + 1 }
	var dec Fn = func(x int) int { return x - 1 }

	x1 := ComposeFn(inc, sqr, dec)(in)
	x2 := ComposeFn(dec, sqr, inc)(in)

	if x1 != out1 || x2 != out2 {
		t.Errorf("inc(sqr(dec(%d))) [%d] != %d || dec(sqr(inc(%d))) [%d] != %d",
			in, x1, out1, in, x2, out2)
	}
}