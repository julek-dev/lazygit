package gocui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// topEdge returns what is drawn on the top edge of the view's frame.
func topEdge(g *Gui, v *View) string {
	edge := ""
	for x := v.x0; x <= v.x1; x++ {
		s, _, _ := g.screen.Get(x, v.y0)
		edge += s
	}
	return edge
}

func TestSubtitleLeavesTheTitleAlone(t *testing.T) {
	g := newTestGui(t)
	v, _ := g.SetView("main", 0, 0, 40, 5, 0)
	v.Title = "Unstaged changes"
	v.Subtitle = "(ignoring whitespace)"
	assert.NoError(t, g.flush())

	assert.Equal(t,
		/* EXPECTED:
		"┌─Unstaged changes─(ignoring white…─────┐",
		ACTUAL: */
		"┌─Unstaged cha(ignoring whitespace)─────┐",
		topEdge(g, v))
}
