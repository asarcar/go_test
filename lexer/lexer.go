// Reference: https://golang.org/src/text/template/parse/lex.go
package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Template System
// Evaluation: {{.Title}}
// Constants and Functions: {{printf "%g: %#3X" 1.2+2i 123}}
// Control Structures {{range $s.Text}} {{.}} {{end}}

// itemType: type of lex items
type itemType int

const (
	itemError itemType = iota
	// itemDot
	itemEOF
	// itemElse
	// itemEnd
	// itemField
	// itemIdentifier
	// itemIf
	itemLeftMeta
	itemNumber
	itemPipe
	// itemRange
	// itemRawString
	itemRightMeta
	itemSpace
	// itemString
	itemText
)

const (
	EOF rune = -(iota + 1)
)

const (
	dotT       = "."
	elseT      = "else"
	ifT        = "if"
	leftMetaT  = "{{"
	pipeT      = "|"
	rightMetaT = "}}"
)

type item struct {
	typ itemType
	val string
}

type lexer struct {
	name  string
	input string
	state stateFn
	start int
	pos   int
	width int
	items chan item
}

type stateFn func(*lexer) stateFn

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// run lexes the input by executing state functions until state is nil
// emits token to the client in the process...
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

// Initialization routing kicks the lexing process in a go routine
func lex(name, input string) chan item {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item, 1),
	}
	go l.run() // Concurrently run state machine
	return l.items
}

// cannot run concurrent routine if lexer needed at init time
func lexSync(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		state: lexText,
		// lexer may emit multiple (<=2) tokens at one invocation
		items: make(chan item, 2),
	}
	return l
}

// nextItem: returns next item from the input
func (l *lexer) nextItem() item {
	for {
		select {
		case i := <-l.items:
			return i
		default:
			// Caller should not call nextItem once EOF is returned
			// Guard code just in case we have a rogue caller
			if l.state == nil {
				return item{itemEOF, l.input[l.pos:l.pos]}
			}
			l.state = l.state(l)
		}
	}
	panic("not reached")
}

// emit passes an item back to client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos // cursor moved to current position
}

func lexText(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], leftMetaT) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexLeftMeta
		}
		if l.next() == EOF {
			break
		}
	}
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func lexLeftMeta(l *lexer) stateFn {
	l.pos += len(leftMetaT)
	l.emit(itemLeftMeta)
	return lexInsideAction
}

func lexRightMeta(l *lexer) stateFn {
	l.pos += len(rightMetaT)
	l.emit(itemRightMeta)
	return lexText
}

func lexInsideAction(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], rightMetaT) {
			return lexRightMeta
		}
		switch r := l.next(); {
		case isEOF(r) || isEndOfLine(r):
			return l.errorf("unclosed action")
		case isPipe(r):
			l.emit(itemPipe)
		case unicode.IsSpace(r):
			return lexSpace
		// case isQuote(r):
		//		return lexQuote
		// case isRawQuote(r):
		//	return lexRawQuote
		// case isAlphaNumeric(r):
		// 	l.backup()
		// 	return lexIdentifier
		case isPlusMinus(r) || unicode.IsDigit(r):
			l.backup()
			return lexNumber
		}
	}
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func lexNumber(l *lexer) stateFn {
	l.accept("+-")
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	l.accept("i")
	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	l.emit(itemNumber)
	return lexInsideAction
}

func lexSpace(l *lexer) stateFn {
	for unicode.IsSpace(l.peek()) {
		l.next()
	}
	l.emit(itemSpace)
	return lexInsideAction
}

func isEndOfLine(r rune) bool {
	return r == '\n' || r == '\r'
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isPipe(r rune) bool {
	return r == '|'
}

func isQuote(r rune) bool {
	return r == '"'
}

func isRawQuote(r rune) bool {
	return r == '`'
}

func isPlusMinus(r rune) bool {
	return r == '+' || r == '-'
}

func isDollar(r rune) bool {
	return r == '$'
}

func isEOF(r rune) bool {
	return r == EOF
}
