package affine2d_test

import (
	"math"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-affine2d"
)

func TestTransformation_Inverse(t *testing.T) {
	for _, tc := range []struct {
		name     string
		t        *affine2d.Transformation
		expected *affine2d.Transformation
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
			assertInDelta(t, tc.expected.Float64s(), tc.t.Inverse().Float64s(), 0)
		})
	}
}

func TestTransformation_Transform(t *testing.T) {
	for _, tc := range []struct {
		name           string
		transformation *affine2d.Transformation
		v              []float64
		expected       []float64
	}{
		{
			name:           "zero_identity",
			transformation: affine2d.Identity(),
			v:              []float64{0, 0},
			expected:       []float64{0, 0},
		},
		{
			name:           "zero_rotate",
			transformation: affine2d.Rotate(math.Pi),
			v:              []float64{0, 0},
			expected:       []float64{0, 0},
		},
		{
			name:           "zero_scale",
			transformation: affine2d.Scale(2, 3),
			v:              []float64{0, 0},
			expected:       []float64{0, 0},
		},
		{
			name:           "zero_shear",
			transformation: affine2d.Shear(2, 3),
			v:              []float64{0, 0},
			expected:       []float64{0, 0},
		},
		{
			name:           "zero_translate",
			transformation: affine2d.Translate(2, 3),
			v:              []float64{0, 0},
			expected:       []float64{2, 3},
		},
		{
			name:           "unit_x_identity",
			transformation: affine2d.Identity(),
			v:              []float64{1, 0},
			expected:       []float64{1, 0},
		},
		{
			name:           "unit_x_rotate",
			transformation: affine2d.Rotate(math.Pi),
			v:              []float64{1, 0},
			expected:       []float64{-1, 0},
		},
		{
			name:           "unit_x_scale",
			transformation: affine2d.Scale(2, 3),
			v:              []float64{1, 0},
			expected:       []float64{2, 0},
		},
		{
			name:           "unit_x_shear",
			transformation: affine2d.Shear(2, 3),
			v:              []float64{1, 0},
			expected:       []float64{1, 3},
		},
		{
			name:           "unit_x_translate",
			transformation: affine2d.Translate(2, 3),
			v:              []float64{1, 0},
			expected:       []float64{3, 3},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertInDelta(t, tc.expected, tc.transformation.Transform(tc.v), 1e-15)
		})
	}
}

func TestTransformation_TransformDirection(t *testing.T) {
	for _, tc := range []struct {
		name           string
		transformation *affine2d.Transformation
		v              []float64
		expected       []float64
	}{
		{
			name:           "identity",
			transformation: affine2d.Identity(),
			v:              []float64{1, 0},
			expected:       []float64{1, 0},
		},
		{
			name:           "rotate",
			transformation: affine2d.Rotate(math.Pi / 2),
			v:              []float64{1, 0},
			expected:       []float64{0, 1},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertInDelta(t, tc.expected, tc.transformation.TransformDirection(tc.v), 1e-16)
		})
	}
}

func TestTransformation_Multiply(t *testing.T) {
	for _, tc := range []struct {
		name             string
		t                *affine2d.Transformation
		u                *affine2d.Transformation
		expectedFloat64s []float64
	}{
		{
			name:             "identity",
			t:                affine2d.Identity(),
			u:                affine2d.Identity(),
			expectedFloat64s: affine2d.Identity().Float64s(),
		},
		{
			name:             "identity_scale",
			t:                affine2d.Identity(),
			u:                affine2d.Scale(2, 3),
			expectedFloat64s: affine2d.Scale(2, 3).Float64s(),
		},
		{
			name:             "rotate_rotate",
			t:                affine2d.Rotate(math.Pi / 2),
			u:                affine2d.Rotate(math.Pi),
			expectedFloat64s: affine2d.Rotate(3 * math.Pi / 2).Float64s(),
		},
		{
			name:             "scale_identity",
			t:                affine2d.Scale(2, 3),
			u:                affine2d.Identity(),
			expectedFloat64s: affine2d.Scale(2, 3).Float64s(),
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
			expectedFloat64s: affine2d.Translate(6, 8).Float64s(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertInDelta(t, tc.expectedFloat64s, tc.t.Multiply(tc.u).Float64s(), 0)
		})
	}
}

func TestTransformation_Then(t *testing.T) {
	for _, tc := range []struct {
		name             string
		t                *affine2d.Transformation
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
			assertInDelta(t, tc.expectedFloat64s, tc.t.Float64s(), 1e-16)
		})
	}
}

func TestTransformation_TransformSlice(t *testing.T) {
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
