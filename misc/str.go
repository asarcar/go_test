package misc

import (
	"fmt"
	"unicode/utf8"
)

func DumpStr() {
	fmt.Println("DumpStr\n-----------")

	// Go source code is always UTF-8.
	// A string literal, absent byte-level escapes, always holds valid UTF-8 sequences.
	const literal_str = `A⌘A�A�A=A�A\A`
	fmt.Printf("plain string: %s; quoted-string: %q; quoted-unicode-string: %+q; hex bytes: % x\n", literal_str, literal_str, literal_str, literal_str)
	for index, runeValue := range literal_str {
		fmt.Printf("%#U starts at byte position %d\n", runeValue, index)
	}

	// A string holds arbitrary bytes.
	// No guarantee is made in Go that characters in strings are normalized:
	// \u00e0 is à which is same as \u0300 followed by \u0061 which is  ̀a
	const literal_str2 = "\\u00e0=\u00e0,\\u0061=\u0061,\\u0300\\u0061=\u0300\u0061"
	fmt.Printf("plain string: %s; quoted-string: %q; quoted-unicode-string: %+q; hex bytes: % x\n", literal_str2, literal_str2, literal_str2, literal_str2)
	for index, runeValue := range literal_str2 {
		fmt.Printf("%#U starts at byte position %d\n", runeValue, index)
	}

	const nihongo = "\xbd\xb2日本語"
	// Those sequences represent Unicode code points, called runes.
	fmt.Printf("plain string: %s: quoted-string: %q; quoted-unicode-string: %+q; hex bytes: % x\n", nihongo, nihongo, nihongo, nihongo)
	// Equalivalent to for index, runeValue range loop over string
	for i, w := 0, 0; i < len(nihongo); i += w {
		runeValue, width := utf8.DecodeRuneInString(nihongo[i:])
		fmt.Printf("%#U starts at byte position %d\n", runeValue, i)
		w = width
	}

	fmt.Println("-----------")
}
