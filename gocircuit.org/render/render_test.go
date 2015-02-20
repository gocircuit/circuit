package render

import (
	"testing"
)

func TestRender(t *testing.T) {
	if Render("{{.X}}", A{"X": 1}) != "1" {
		t.Errorf("mismatch")
	}
}

func TestHtml(t *testing.T) {
	println(RenderHtml("hola", "amiga"))
}
