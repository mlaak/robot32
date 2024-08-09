package translator

import (
	"testing"
)

// Not much to test here at the moment, but lets leave the skeleton here
func TestTranslator(t *testing.T) {

	s := TR(0, "Hello")

	if s != "Hello" {
		t.Fatalf("TestTranslator Fail 1 - expected to het back Hello but go " + s)
	}
}
