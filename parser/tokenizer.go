package parser

import (
	"io"
	"fmt"
	"bufio"
	"unicode"
	"strings"
)

type EventHandler interface {
	Error(error)
	Literal(string)
	Strval(string)

	Import()
	Includes()
	Route()
	Using()
	Does()
	From()
}

type Tokenizer struct {
	input *bufio.Reader
	lastErr error
	event EventHandler
}

func (z *Tokenizer) Next() {
	z.consumeSpace()

	b, _, err := z.input.ReadRune()
	if err != nil {
		z.err(err)
		return
	}

	switch b {
	case '`':
		z.literal()
	case '«':
		z.altLiteral()
	case '"':
		z.dquote()
	case '\'':
		z.squote()
	//case ' ', '\t', '\n', '\r', '\v', '\f', 0x85 /* NEL */, 0xA0 /* NBSP */:
		// consume whitespace.
	default:
		z.word(b)
	}
}

func (z *Tokenizer) err(e error) {
	z.event.Error(e)
	z.lastErr = e
}

func (z *Tokenizer) literal() {
	str, err := z.input.ReadString('`')
	if err != nil {
		z.err(err)
		return
	}
	z.event.Literal(strings.TrimSuffix(str, "`"))
}
func (z *Tokenizer) altLiteral() {
	str, err := z.input.ReadString('»')
	if err != nil {
		z.err(err)
		return
	}
	z.event.Literal(strings.TrimSuffix(str, "»"))
}

func (z *Tokenizer) dquote() {
	str, err := z.input.ReadString('"')
	if err != nil {
		z.err(err)
		return
	}
	z.event.Strval(strings.TrimSuffix(str, `"`))
}
func (z *Tokenizer) squote() {
	str, err := z.input.ReadString('"')
	if err != nil {
		z.err(err)
		return
	}
	z.event.Strval(strings.TrimSuffix(str, `'`))
}

var (
	mport = "MPORT"
	nclude = "NCLUDES"
	oute = "OUTE"
	sing = "SING"
	rom = "ROM"
	oes = "OES"
)

func (z *Tokenizer) word(b rune) {

	switch b {
	case 'I':
		if z.peekMatch(mport) {
			z.imports()
			return
		} else if z.peekMatch(nclude) {
			z.include()
			return
		}
		//z.input.UnreadRune()
		z.bareword([]rune{b})
		return
	case 'R': // ROUTE
		if z.peekMatch(oute) {
			z.route()
			return
		}

		z.bareword([]rune{b})
		return

	case 'U': // USING
		if z.peekMatch(sing) {
			z.using()
			return
		}
		z.bareword([]rune{b})
		return
	case 'D': // DOES
		if z.peekMatch(oes) {
			z.does()
			return
		}
		z.bareword([]rune{b})
		return
	case 'F': // FROM
		if z.peekMatch(rom) {
			z.from()
			return
		}
		z.bareword([]rune{b})
		return
	default:
		z.bareword([]rune{b})
	}
}

func (z *Tokenizer) peekMatch(word string) bool {
	size := len(word)
	p, err := z.input.Peek(size + 1)
	if err != nil && err != io.EOF {
		fmt.Printf("Received peek error. Please report: %s\n", err)
	}
	matches := len(p) >= size && string(p[0:size]) == word

	if len(p) == size + 1 && !unicode.IsSpace(rune(p[size])) {
		return false
	}

	if matches {
		// throw-away buffer.
		buf := make([]byte, size)
		z.input.Read(buf)
	}

	return matches
}

func (z *Tokenizer) bareword(prepend []rune) {
	/*
	if unread {
		if err := z.input.UnreadRune(); err != nil {
			fmt.Printf("Error unreading: %s\n", err)
		}
	}
	*/
	buf := prepend
	r, _, err := z.input.ReadRune()
	for {
		if err != nil {
			if len(buf) > 0 {
				z.event.Strval(string(buf))
			}
			z.event.Error(err)
			return
		} else if unicode.IsSpace(r) {
			z.event.Strval(string(buf))
			// And consume the space?
			return
		}
		buf = append(buf, r)
		r, _, err = z.input.ReadRune()
	}

	z.event.Strval(string(buf))
}

func (z *Tokenizer) consumeSpace() {
	r, _, err := z.input.ReadRune()
	for {
		if err != nil {
			z.event.Error(err)
			return
		} else if !unicode.IsSpace(r) {
			z.input.UnreadRune()
			return
		}
		r, _, err = z.input.ReadRune()
	}
}

func (z *Tokenizer) imports() {
	z.event.Import()
}

func (z *Tokenizer) include() {
	z.event.Includes()
}

func (z *Tokenizer) from() {
	z.event.From()
}

func (z *Tokenizer) does() {
	z.event.Does()
}

func (z *Tokenizer) using() {
	z.event.Using()
}

func (z *Tokenizer) route() {
	z.event.Route()
}


func NewTokenizer(input io.Reader, e EventHandler) *Tokenizer {
	z := Tokenizer{
		input: bufio.NewReader(input),
		event: e,
	}

	return &z
}

