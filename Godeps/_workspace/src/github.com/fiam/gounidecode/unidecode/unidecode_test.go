package unidecode

import (
	"testing"
)

func testTransliteration(original string, decoded string, t *testing.T) {
	if r := Unidecode(original); r != decoded {
		t.Errorf("Expected '%s', got '%s'\n", decoded, r)
	}
}

func TestASCII(t *testing.T) {
	s := "ABCDEF"
	testTransliteration(s, s, t)
}

func TestKnosos(t *testing.T) {
	o := "Κνωσός"
	d := "Knosos"
	testTransliteration(o, d, t)
}

func TestBeiJing(t *testing.T) {
	o := "\u5317\u4EB0"
	d := "Bei Jing "
	testTransliteration(o, d, t)
}
