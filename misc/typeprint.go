package misc

import (
	"fmt"
)

func DumpTypePrint(args ...interface{}) {
	fmt.Println("DumpTypePrint\n----------------")
	for i, arg := range args {
		if i > 0 {
			fmt.Printf(" ")
		}

		switch a := arg.(type) { // type switch
		case int:
			fmt.Printf("%d", a)
		case float64:
			fmt.Printf("%.8f", a)
		case float32:
			fmt.Printf("%.4f", a)
		case string:
			fmt.Printf(a)
		default:
			fmt.Printf("???")
		}
		if i+1 == len(args) {
			fmt.Printf("\n")
		}
	}
	fmt.Println("----------------")
}
