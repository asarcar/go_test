package misc

import (
//	"fmt"
)

// https://leetcode.com/problems/regular-expression-matching/
// Given an input string (s) and a pattern (p), implement regular expression matching with support for '.' and '*'.
//	'.' Matches any single character.
//	'*' Matches zero or more of the preceding element.
//	The matching should cover the entire input string (not partial).
// Note:
//	s could be empty and contains only lowercase letters a-z.
//	p could be empty and contains only lowercase letters a-z, and characters like . or *.

const (
	DOT  = '.'
	STAR = '*'
)

type RegExMatch struct {
	regexp string
}

func NewRegExMatch(s string) *RegExMatch {
	// validate regex
	return &RegExMatch{s}
}

// first byte is regexp token, second byte comes from string
func match_token(r, c byte) bool {
	match := (r == DOT || c == r)
	return match
}

func match_reg_exp(str string, str_pos int, regexp string, regexp_pos int) bool {
	len_str, len_regexp := len(str), len(regexp)

	// fmt.Printf("str \"%s\"[%v], reg \"%s\"[%v]\n",
	// 	 str, str_pos, regexp, regexp_pos)

	// Case 1: reached end of string as well as the regular expression = match is true
	if str_pos >= len_str && regexp_pos >= len_regexp {
		return true
	}

	// Case 2: reached end of regexp but not end of string = match is false
	if regexp_pos >= len_regexp {
		return false
	}

	// Case 3:
	// a. not end of regexp but end of string
	// b. neither end of regexp nor end of string

	// get current string token
	c := byte(0)
	if str_pos < len_str {
		c = str[str_pos]
	}

	// get current and next regexp token
	r_next := byte(0)
	r := regexp[regexp_pos]
	if regexp_pos < len_regexp-1 {
		r_next = regexp[regexp_pos+1]
	}

	// match for X* where X is a character
	if r_next == STAR {
		// X* is 2 characters - so jump 2 position of characters
		empty_match := match_reg_exp(str, str_pos, regexp, regexp_pos+2)
		if empty_match == true {
			return true
		}
		// regexp next token does not match - we've failed empty match

		// string and regexp current token test match
		if match_token(r, c) == false {
			return false
		}
		// next token matched : check more chars of str matches current regexp token
		return match_reg_exp(str, str_pos+1, regexp, regexp_pos)
	}

	// move to next regexp token on match
	// move ahead to next char in str and next token in regexp
	if match_token(r, c) == true {
		return match_reg_exp(str, str_pos+1, regexp, regexp_pos+1)
	}
	return false
}

func (pR *RegExMatch) Match(str string) bool {
	return match_reg_exp(str, 0, pR.regexp, 0)
}
