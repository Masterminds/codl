package parser

import (
	"os"
	"testing"
	"strings"
)

func TestSerialize(t *testing.T) {
	doc := `IMPORT
		github.com/Masterminds/cookoo/web
		github.com/Masterminds/cookoo/cli

	ROUTE "test" "TEST"
		DOES web.Flush "first"
			USING p1 defval FROM cxt:p1 cxt:p2
			USING p2 «1» FROM cxt:p2
		DOES cli.ParseArgs "CMD"
	ROUTE "foo" "BAR"
		DOES web.Flush cmd3`
	input := strings.NewReader(doc)

	h, err := Parse(input)
	if err != nil {
		t.Errorf("Surprise! Error: %s", err)
	}

	reg := h.(Registry)
	ser := NewSerializer("test", "serializertest", os.Stdout, reg)
	if err := ser.Write(); err != nil {
		t.Errorf("Failed to serialize: %s", err)
	}
}
