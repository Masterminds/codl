package parser

import (
	"testing"
	"strings"
)

func TestSimpleParse(t *testing.T) {
	doc := `IMPORT foo`
	input := strings.NewReader(doc)

	h, err := Parse(input)
	if err != nil {
		t.Errorf("Surprise! Error: %s", err)
	}

	handy := h.(*handler)

	if len(handy.imports) != 1 {
		t.Errorf("Expected 1 import, got %d", len(handy.imports))
	}

	if handy.imports[0] != "`foo`" {
		t.Errorf("Expected \"foo\", got %s", handy.imports[0])
	}

}
