package controllers

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_increasedContextSize(t *testing.T) {
	scenarios := []struct {
		name           string
		size           uint64
		step           uint64
		expectedResult uint64
	}{
		{
			name:           "increase by one",
			size:           3,
			step:           1,
			expectedResult: 4,
		},
		{
			name:           "increase by a bigger step",
			size:           3,
			step:           4,
			expectedResult: 7,
		},
		{
			name:           "increase from zero",
			size:           0,
			step:           1,
			expectedResult: 1,
		},
		{
			name:           "a zero step is a no-op",
			size:           5,
			step:           0,
			expectedResult: 5,
		},
		{
			name:           "increase exactly up to the maximum",
			size:           math.MaxUint64 - 4,
			step:           4,
			expectedResult: math.MaxUint64,
		},
		{
			name:           "increase past the maximum clamps instead of overflowing",
			size:           math.MaxUint64 - 1,
			step:           4,
			expectedResult: math.MaxUint64,
		},
		{
			name:           "increase at the maximum stays at the maximum",
			size:           math.MaxUint64,
			step:           1,
			expectedResult: math.MaxUint64,
		},
		{
			name:           "a huge step clamps at the maximum",
			size:           1,
			step:           math.MaxUint64,
			expectedResult: math.MaxUint64,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expectedResult, increasedContextSize(s.size, s.step))
		})
	}
}

func Test_decreasedContextSize(t *testing.T) {
	scenarios := []struct {
		name           string
		size           uint64
		step           uint64
		min            uint64
		expectedResult uint64
	}{
		{
			name:           "decrease by one",
			size:           3,
			step:           1,
			min:            0,
			expectedResult: 2,
		},
		{
			name:           "decrease by a bigger step",
			size:           11,
			step:           4,
			min:            0,
			expectedResult: 7,
		},
		{
			name:           "decrease down to zero",
			size:           1,
			step:           1,
			min:            0,
			expectedResult: 0,
		},
		{
			name:           "decrease exactly down to the minimum",
			size:           5,
			step:           2,
			min:            3,
			expectedResult: 3,
		},
		{
			name:           "decrease past the minimum clamps at the minimum",
			size:           4,
			step:           2,
			min:            3,
			expectedResult: 3,
		},
		{
			name:           "decrease at the minimum stays at the minimum",
			size:           3,
			step:           1,
			min:            3,
			expectedResult: 3,
		},
		{
			name:           "decrease from below the minimum snaps up to the minimum",
			size:           1,
			step:           1,
			min:            3,
			expectedResult: 3,
		},
		{
			name:           "a step bigger than the size clamps instead of underflowing",
			size:           3,
			step:           4,
			min:            0,
			expectedResult: 0,
		},
		{
			name:           "a zero step is a no-op",
			size:           5,
			step:           0,
			min:            0,
			expectedResult: 5,
		},
		{
			name:           "a zero step still snaps up to the minimum",
			size:           1,
			step:           0,
			min:            3,
			expectedResult: 3,
		},
		{
			name:           "decrease from the maximum",
			size:           math.MaxUint64,
			step:           1,
			min:            0,
			expectedResult: math.MaxUint64 - 1,
		},
		{
			name:           "a huge step clamps at the minimum",
			size:           5,
			step:           math.MaxUint64,
			min:            2,
			expectedResult: 2,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.Equal(t, s.expectedResult, decreasedContextSize(s.size, s.step, s.min))
		})
	}
}
