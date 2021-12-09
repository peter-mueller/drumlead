package parser

import (
	"testing"
)

func TestLexer(t *testing.T) {
	items := Lex(input)

	expected := []item{
		{typ: HASHTAG, value: "#"},
		{typ: HASHTAG, value: "#"},
		{typ: TEXT, value: "Hello"},
		{typ: TEXT, value: "asdf"},
		{typ: TEXT, value: "asdf"},
		{typ: TEXT, value: ""},
		{typ: HASHTAG, value: "#"},
		{typ: HASHTAG, value: "#"},
		{typ: TEXT, value: ""},
		{typ: HASHTAG, value: "#"},
		{typ: TEXT, value: "omg"},
	}

	for _, e := range expected {
		i, ok := <-items
		if !ok {
			t.Fatal()
		}
		if i.typ != e.typ {
			t.Fatalf("expected typ %d but got %d", e.typ, i.typ)
		}
		if i.value != e.value {
			t.Fatalf("expected value %s but got %s", e.value, i.value)
		}

		t.Logf("typ: %d value: %s", i.typ, e.value)
	}
}

const input = `## Hello
asdf
asdf

##
# omg

adsf

`
