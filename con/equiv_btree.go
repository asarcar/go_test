package main

import (
	"code.google.com/p/go-tour/tree"
	"fmt"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	if t == nil {
		close(ch)
		return
	}

	chl := make(chan int, 10)
	chr := make(chan int, 10)
	go Walk(t.Left, chl)
	go Walk(t.Right, chr)

	for v := range chl {
		ch <- v
	}
	ch <- t.Value
	for v := range chr {
		ch <- v
	}

	close(ch)

	return
}

func DumpTree(t *tree.Tree) string {
	if t == nil {
		return ""
	}
	var sV, sL, sR string
	sL = DumpTree(t.Left)
	sV = fmt.Sprintf("%v ", t.Value)
	sR = DumpTree(t.Right)

	return sL + sV + sR
}

func DumpChannel(ch chan int) (str string) {
	for v := range ch {
		str += fmt.Sprintf("%v ", v)
	}
	return str
}

func StuffChannel(v int, ch chan int) {
	for i := 1; i <= 10; i++ {
		ch <- v * i
	}
	close(ch)
	return
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	v1, v2 := 0, 0
	ok1, ok2 := true, true

	for ok1 && ok2 {
		v1, ok1 = <-ch1
		v2, ok2 = <-ch2
		if ok1 != ok2 || v1 != v2 {
			return false
		}
	}
	return true
}

func main() {
	ch := make(chan int, 10)
	// go StuffChannel(2, ch)
	// go Walk(&tree.Tree{&tree.Tree{nil, 5, nil}, 10, &tree.Tree{nil, 15, nil}}, ch)
	go Walk(tree.New(2), ch)
	// fmt.Println("Tree1:", DumpTree(tree.New(1)))
	fmt.Println("Tree(2):", DumpChannel(ch))
	fmt.Println("Tree(1) eq Tree(1):", Same(tree.New(1), tree.New(1)))
	fmt.Println("Tree(1) eq Tree(2):", Same(tree.New(1), tree.New(2)))
}
