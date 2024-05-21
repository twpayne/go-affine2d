package affine2d_test

import (
	"math"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-affine2d"
)

func TestTransform_Inverse(t *testing.T) {
	for _, tc := range []struct {
		name     string
		t        *affine2d.Transform
		expected *affine2d.Transform
	}{
		{
			name:     "identity",
			t:        affine2d.Identity(),
			expected: affine2d.Identity(),
		},
		{
			name:     "rotate",
			t:        affine2d.Rotate(math.Pi / 2),
			expected: affine2d.Rotate(-math.Pi / 2),
		},
		{
			name:     "scale",
			t:        affine2d.Scale(2, 4),
			expected: affine2d.Scale(0.5, 0.25),
		},
		{
			name:     "translate",
			t:        affine2d.Translate(2, 3),
			expected: affine2d.Translate(-2, -3),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertInDelta(t, tc.expected.Float64Slice(), tc.t.Inverse().Float64Slice(), 0)
		})
	}
}

func TestTransform_Float64Array(t *testing.T) {
	assert.Equal(t, [6]float64{1, 2, 3, 4, 5, 6}, affine2d.NewTransform([6]float64{1, 2, 3, 4, 5, 6}).Float64Array())
}

func TestTransform_New(t *testing.T) {
	for _, tc := range []struct {
		name      string
		transform *affine2d.Transform
		expected  [6]float64
	}{
		{
			name:      "new",
			transform: affine2d.NewTransform([6]float64{1, 2, 3, 4, 5, 6}),
			expected: [6]float64{
				1, 2, 3,
				4, 5, 6,
			},
		},
		{
			name:      "delta",
			transform: affine2d.Delta([]float64{0, 1}, []float64{2, 3}),
			expected: [6]float64{
				2, -2, 0,
				2, 2, 1,
			},
		},
		{
			name:      "identity",
			transform: affine2d.Identity(),
			expected: [6]float64{
				1, 0, 0,
				0, 1, 0,
			},
		},
		{
			name:      "rotate",
			transform: affine2d.Rotate(math.Pi),
			expected: [6]float64{
				-1, 0, 0,
				0, -1, 0,
			},
		},
		{
			name:      "scale",
			transform: affine2d.Scale(2, 3),
			expected: [6]float64{
				2, 0, 0,
				0, 3, 0,
			},
		},
		{
			name:      "translate",
			transform: affine2d.Translate(2, 3),
			expected: [6]float64{
				1, 0, 2,
				0, 1, 3,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertInDelta(t, tc.expected[:], tc.transform.Float64Slice(), 1e-15)
		})
	}
}

func TestTransform_Transform(t *testing.T) {
	for _, tc := range []struct {
		name      string
		transform *affine2d.Transform
		p         []float64
		expected  []float64
	}{
		{
			name:      "zero_identity",
			transform: affine2d.Identity(),
			p:         []float64{0, 0},
			expected:  []float64{0, 0},
		},
		{
			name:      "zero_rotate",
			transform: affine2d.Rotate(math.Pi),
			p:         []float64{0, 0},
			expected:  []float64{0, 0},
		},
		{
			name:      "zero_scale",
			transform: affine2d.Scale(2, 3),
			p:         []float64{0, 0},
			expected:  []float64{0, 0},
		},
		{
			name:      "zero_shear",
			transform: affine2d.Shear(2, 3),
			p:         []float64{0, 0},
			expected:  []float64{0, 0},
		},
		{
			name:      "zero_translate",
			transform: affine2d.Translate(2, 3),
			p:         []float64{0, 0},
			expected:  []float64{2, 3},
		},
		{
			name:      "unit_x_identity",
			transform: affine2d.Identity(),
			p:         []float64{1, 0},
			expected:  []float64{1, 0},
		},
		{
			name:      "unit_x_rotate",
			transform: affine2d.Rotate(math.Pi),
			p:         []float64{1, 0},
			expected:  []float64{-1, 0},
		},
		{
			name:      "unit_x_scale",
			transform: affine2d.Scale(2, 3),
			p:         []float64{1, 0},
			expected:  []float64{2, 0},
		},
		{
			name:      "unit_x_shear",
			transform: affine2d.Shear(2, 3),
			p:         []float64{1, 0},
			expected:  []float64{1, 3},
		},
		{
			name:      "unit_x_translate",
			transform: affine2d.Translate(2, 3),
			p:         []float64{1, 0},
			expected:  []float64{3, 3},
		},
		{
			name:      "origin_delta",
			transform: affine2d.Delta([]float64{0, 1}, []float64{2, 3}),
			p:         []float64{0, 0},
			expected:  []float64{0, 1},
		},
		{
			name:      "unit_x_delta",
			transform: affine2d.Delta([]float64{0, 1}, []float64{2, 3}),
			p:         []float64{1, 0},
			expected:  []float64{2, 3},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertInDelta(t, tc.expected, tc.transform.Transform(tc.p), 1e-15)
		})
	}
}

func TestTransform_TransformDirection(t *testing.T) {
	for _, tc := range []struct {
		name      string
		transform *affine2d.Transform
		v         []float64
		expected  []float64
	}{
		{
			name:      "identity",
			transform: affine2d.Identity(),
			v:         []float64{1, 0},
			expected:  []float64{1, 0},
		},
		{
			name:      "rotate",
			transform: affine2d.Rotate(math.Pi / 2),
			v:         []float64{1, 0},
			expected:  []float64{0, 1},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertInDelta(t, tc.expected, tc.transform.TransformDirection(tc.v), 1e-16)
		})
	}
}

func TestTransform_Multiply(t *testing.T) {
	for _, tc := range []struct {
		name             string
		t                *affine2d.Transform
		u                *affine2d.Transform
		expectedFloat64s []float64
	}{
		{
			name:             "identity",
			t:                affine2d.Identity(),
			u:                affine2d.Identity(),
			expectedFloat64s: affine2d.Identity().Float64Slice(),
		},
		{
			name:             "identity_scale",
			t:                affine2d.Identity(),
			u:                affine2d.Scale(2, 3),
			expectedFloat64s: affine2d.Scale(2, 3).Float64Slice(),
		},
		{
			name:             "rotate_rotate",
			t:                affine2d.Rotate(math.Pi / 2),
			u:                affine2d.Rotate(math.Pi),
			expectedFloat64s: affine2d.Rotate(3 * math.Pi / 2).Float64Slice(),
		},
		{
			name:             "scale_identity",
			t:                affine2d.Scale(2, 3),
			u:                affine2d.Identity(),
			expectedFloat64s: affine2d.Scale(2, 3).Float64Slice(),
		},
		{
			name: "scale_translate",
			t:    affine2d.Scale(2, 3),
			u:    affine2d.Translate(4, 5),
			expectedFloat64s: []float64{
				2, 0, 8,
				0, 3, 15,
			},
		},
		{
			name: "translate_scale",
			t:    affine2d.Translate(4, 5),
			u:    affine2d.Scale(2, 3),
			expectedFloat64s: []float64{
				2, 0, 4,
				0, 3, 5,
			},
		},
		{
			name:             "translate_translate",
			t:                affine2d.Translate(2, 3),
			u:                affine2d.Translate(4, 5),
			expectedFloat64s: affine2d.Translate(6, 8).Float64Slice(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertInDelta(t, tc.expectedFloat64s, tc.t.Multiply(tc.u).Float64Slice(), 0)
		})
	}
}

func TestTransform_Then(t *testing.T) {
	for _, tc := range []struct {
		name             string
		t                *affine2d.Transform
		expectedFloat64s []float64
	}{
		{
			name: "identity",
			t:    affine2d.Identity(),
			expectedFloat64s: []float64{
				1, 0, 0,
				0, 1, 0,
			},
		},
		{
			name: "rotate_then_translate",
			t:    affine2d.Rotate(math.Pi/2).Translate(2, 3),
			expectedFloat64s: []float64{
				0, -1, 2,
				1, 0, 3,
			},
		},
		{
			name: "scale_then_translate",
			t:    affine2d.Scale(2, 3).Translate(4, 5),
			expectedFloat64s: []float64{
				2, 0, 4,
				0, 3, 5,
			},
		},
		{
			name: "translate_then_rotate",
			t:    affine2d.Translate(2, 3).Rotate(math.Pi / 2),
			expectedFloat64s: []float64{
				0, -1, -3,
				1, 0, 2,
			},
		},
		{
			name: "translate_then_scale",
			t:    affine2d.Translate(4, 5).Scale(2, 3),
			expectedFloat64s: []float64{
				2, 0, 8,
				0, 3, 15,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertInDelta(t, tc.expectedFloat64s, tc.t.Float64Slice(), 1e-16)
		})
	}
}

func TestTransform_TransformSlice(t *testing.T) {
	unitSquare := [][]float64{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
	}
	actual := affine2d.Scale(2, 3).TransformSlice(unitSquare)
	assert.Equal(t, [][]float64{
		{0, 0},
		{2, 0},
		{2, 3},
		{0, 3},
	}, actual)
}

func assertInDelta(t *testing.T, expected, actual []float64, maxDelta float64) {
	t.Helper()
	assert.Equal(t, len(expected), len(actual))
	for i := range expected {
		if actualDelta := math.Abs(expected[i] - actual[i]); actualDelta > maxDelta {
			t.Fatalf("Expected values to be within %e but |%e-%e| is %e at position %d", maxDelta, expected[i], actual[i], actualDelta, i)
		}
	}
}
