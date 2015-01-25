package misc

import (
	"flag"
	"fmt"
)

type TZ uint32

const (
	bit0, mask0 = 1 << iota, 1<<iota - 1
	bit1, mask1
	bit2, mask2
)

func DumpFlags() {
	fmt.Println("DumpFlags\n-------------")

	bPtr := flag.Bool("b", false, "no\n")
	iPtr := flag.Int("i", 10, "value\n")
	sPtr := flag.String("s", "str", "string\n")

	flag.Parse()

	s := ""
	for i := 0; i < flag.NArg(); i++ {
		if i > 0 {
			s += " "
		}
		s += flag.Arg(i)
	}

	fmt.Println("bit0", bit0, "mask0", mask0)
	fmt.Println("bit1", bit1, "mask1", mask1)
	fmt.Println("bit2", bit2, "mask2", mask2)
	fmt.Printf("Bool b %v: Int i %v: String s '%s'\n", *bPtr, *iPtr, *sPtr)
	fmt.Print("Flags '", s, "'")
	fmt.Println(": Tail ", flag.Args())

	fmt.Println("-----------")
}
