package lexer

import "testing"

func TestLexer1(t *testing.T) {
	const name, in = "test1", "abc 123"
	out := []item{item{itemText, "abc 123"}}
	i := 0
	ch := lex(name, in)
	for item := range ch {
		if item.typ == itemEOF {
			break
		}
		if item.typ != out[i].typ || item.val != out[i].val {
			t.Errorf("token# %d: Lex(%s) = {%d, %s} want {%d, %s}",
				i, in, item.typ, item.val, out[i].typ, out[i].val)
		}
		i++
	}
}

func TestLexer2(t *testing.T) {
	const name, in = "test2", "abc 123 {{234}} def 456"
	out := []item{
		item{itemText, "abc 123 "},
		item{itemLeftMeta, leftMetaT},
		item{itemNumber, "234"},
		item{itemRightMeta, rightMetaT},
		item{itemText, " def 456"},
	}
	i := 0
	ch := lex(name, in)
	for item := range ch {
		if item.typ == itemEOF {
			break
		}
		if item.typ != out[i].typ || item.val != out[i].val {
			t.Errorf("token# %d: Lex(%s) = {%d, %s} want {%d, %s}",
				i, in, item.typ, item.val, out[i].typ, out[i].val)
		}
		i++
	}
}

func TestLexer3(t *testing.T) {
	const name, in = "test3", "{{234 |  456|   78.25}}"
	out := []item{
		item{itemLeftMeta, leftMetaT},
		item{itemNumber, "234"},
		item{itemSpace, " "},
		item{itemPipe, pipeT},
		item{itemSpace, "  "},
		item{itemNumber, "456"},
		item{itemPipe, pipeT},
		item{itemSpace, "   "},
		item{itemNumber, "78.25"},
		item{itemRightMeta, rightMetaT},
	}
	i := 0
	ch := lex(name, in)
	for item := range ch {
		if item.typ == itemEOF {
			break
		}
		if item.typ != out[i].typ || item.val != out[i].val {
			t.Errorf("token# %d: Lex(%s) = {%d, %s} want {%d, %s}",
				i, in, item.typ, item.val, out[i].typ, out[i].val)
		}
		i++
	}
}

func TestLexerSync4(t *testing.T) {
	const name, in = "test3", "{{234 |  456|   78.25}}"
	out := []item{
		item{itemLeftMeta, leftMetaT},
		item{itemNumber, "234"},
		item{itemSpace, " "},
		item{itemPipe, pipeT},
		item{itemSpace, "  "},
		item{itemNumber, "456"},
		item{itemPipe, pipeT},
		item{itemSpace, "   "},
		item{itemNumber, "78.25"},
		item{itemRightMeta, rightMetaT},
	}
	i := 0
	l := lexSync(name, in)
	for item := l.nextItem(); item.typ != itemEOF; item = l.nextItem() {
		if item.typ != out[i].typ || item.val != out[i].val {
			t.Errorf("token# %d: Lex(%s) = {%d, %s} want {%d, %s}",
				i, in, item.typ, item.val, out[i].typ, out[i].val)
		}
		i++
	}
}
