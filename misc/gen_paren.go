package main

import(
    "flag"
    "fmt"
)

// Given n pairs of parentheses, write a function to generate all combinations of well-formed parentheses.
// For example, given n = 3, a solution set is:
// [
//   "((()))",
//   "(()())",
//   "(())()",
//   "()(())",
//   "()()()"
// ]

type Par struct {
     n int
}

func (p *Par) Parens(s string, xtra_open int) [] string {
     slen := len(s)
     if slen + xtra_open > 2*p.n || xtra_open < 0 {
     	panic(fmt.Sprintf("string-so-far=\"%s\"[%d], xtra_open=%d, n=%d", s, slen, xtra_open, p.n))
     }
     
     // all open and closed bracket combinations have been discovered
     if slen == 2*p.n {
        return []string{s}
     }
     
     // open bracket and closed bracket options are matched
     // we can only start with another open bracket
     if xtra_open == 0 {
       	return p.Parens(s + "(", xtra_open + 1)
     } 

     // all open brackets exhausted, only closed brackets
     if slen + xtra_open >= 2*p.n {
        return p.Parens(s + ")", xtra_open - 1)
     
     } 

     // both and closed brackets
     s1s := p.Parens(s + ")", xtra_open - 1)
     s2s := p.Parens(s + "(", xtra_open + 1)
     return append(s1s, s2s...) 
} 

func (p *Par) String() string {
     strs := p.Parens("", 0) 
     s := "[\n" 

     for _, str := range strs {
     	 s = s + "  " + str + "\n"
     }

     return s + "]\n"
}

func main() {
     nPtr := flag.Int("n", 2, "number of brackets")
     flag.Parse()
     p := Par{*nPtr}
     fmt.Printf("n=%d: brackets %v", *nPtr, p.String())
}