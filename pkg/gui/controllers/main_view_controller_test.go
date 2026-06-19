package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextChangeRegionIdx(t *testing.T) {
	// true marks a changed line; consecutive changed lines form one region, so
	// regions start at indices 5 and 10.
	changes := []bool{
		false, false, false, false, false,
		true, true,
		false, false, false,
		true, true,
	}
	isChange := func(i int) bool { return changes[i] }

	scenarios := []struct {
		name        string
		fromIdx     int
		forward     bool
		expectedIdx int
		expectFound bool
	}{
		{"forward from the top lands on the first region", 0, true, 5, true},
		{"forward from a region start lands on the next region", 5, true, 10, true},
		{"forward from within a region lands on the next region", 6, true, 10, true},
		{"forward past the last region finds nothing", 10, true, 0, false},
		{"backward from the bottom lands on the last region", 11, false, 10, true},
		{"backward from a region start lands on the previous region", 10, false, 5, true},
		{"backward from within a region lands on its start", 6, false, 5, true},
		{"backward before the first region finds nothing", 5, false, 0, false},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			idx, found := nextChangeRegionIdx(len(changes), s.fromIdx, s.forward, isChange)
			assert.Equal(t, s.expectFound, found)
			if s.expectFound {
				assert.Equal(t, s.expectedIdx, idx)
			}
		})
	}

	t.Run("a view with no changes finds nothing", func(t *testing.T) {
		_, found := nextChangeRegionIdx(3, 0, true, func(int) bool { return false })
		assert.False(t, found)
	})
}

func TestIsDiffChangeMarker(t *testing.T) {
	scenarios := []struct {
		line     string
		expected bool
	}{
		{"+added line", true},
		{"-removed line", true},
		{" context line", false},
		{"@@ -1,4 +1,4 @@", false},
		{"+++ b/file", false},
		{"--- a/file", false},
		{"diff --git a/file b/file", false},
		{"", false},
	}

	for _, s := range scenarios {
		t.Run(s.line, func(t *testing.T) {
			assert.Equal(t, s.expected, isDiffChangeMarker(s.line))
		})
	}
}
