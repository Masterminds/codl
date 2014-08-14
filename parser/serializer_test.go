package parser

import (
	"os"
	"testing"
	"strings"
)

func TestSerialize(t *testing.T) {
	doc := `IMPORT foo bar baz
	ROUTE "test" "TEST"
		DOES cmd1 "first"
			USING p1 defval FROM cxt:p1
			USING p2 defval2 FROM cxt:p2
		DOES cmd2 "second"
	ROUTE "foo" "BAR"`
	input := strings.NewReader(doc)

	h, err := Parse(input)
	if err != nil {
		t.Errorf("Surprise! Error: %s", err)
	}

	reg := h.(Registry)
	ser := NewSerializer(os.Stdout, reg)
	if err := ser.Write(); err != nil {
		t.Errorf("Failed to serialize: %s", err)
	}
}
