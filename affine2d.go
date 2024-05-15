// Package affine2d implements 2D affine transformations.
package affine2d

import "math"

// A Transform is a 2D affine transform.
type Transform struct {
	m [6]float64
}

// Identity returns a new identity transform.
func Identity() *Transform {
	return &Transform{
		m: [6]float64{
			1, 0, 0,
			0, 1, 0,
		},
	}
}

// Rotate returns a new rotate transform.
func Rotate(theta float64) *Transform {
	cosTheta := math.Cos(theta)
	sinTheta := math.Sin(theta)
	return &Transform{
		m: [6]float64{
			cosTheta, -sinTheta, 0,
			sinTheta, cosTheta, 0,
		},
	}
}

// Scale returns a new scale transform.
func Scale(sx, sy float64) *Transform {
	return &Transform{
		m: [6]float64{
			sx, 0, 0,
			0, sy, 0,
		},
	}
}

// Shear returns a new shear transform.
func Shear(sx, sy float64) *Transform {
	return &Transform{
		m: [6]float64{
			1, sx, 0,
			sy, 1, 0,
		},
	}
}

// Translate returns a new translate transform.
func Translate(tx, ty float64) *Transform {
	return &Transform{
		m: [6]float64{
			1, 0, tx,
			0, 1, ty,
		},
	}
}

// Float64s returns the coefficients of t's matrix in row order. The returned
// slice must not be modified.
func (t *Transform) Float64s() []float64 {
	return t.m[:]
}

func (t *Transform) Inverse() *Transform {
	d := t.m[0]*t.m[4] - t.m[1]*t.m[3]
	return &Transform{
		m: [6]float64{
			t.m[4] / d, -t.m[1] / d, (t.m[1]*t.m[5] - t.m[2]*t.m[4]) / d,
			-t.m[3] / d, t.m[0] / d, (t.m[2]*t.m[3] - t.m[0]*t.m[5]) / d,
		},
	}
}

// Multiply returns a new transform which is u then t.
func (t *Transform) Multiply(u *Transform) *Transform {
	return &Transform{
		m: [6]float64{
			t.m[0]*u.m[0] + t.m[1]*u.m[3], t.m[0]*u.m[1] + t.m[1]*u.m[4], t.m[0]*u.m[2] + t.m[1]*u.m[5] + t.m[2],
			t.m[3]*u.m[0] + t.m[4]*u.m[3], t.m[3]*u.m[1] + t.m[4]*u.m[4], t.m[3]*u.m[2] + t.m[4]*u.m[5] + t.m[5],
		},
	}
}

// Rotate returns a new transform which is t then a rotation.
func (t *Transform) Rotate(theta float64) *Transform {
	// FIXME optimize matrix multiplication as we know the structure of the rotate matrix
	return t.Then(Rotate(theta))
}

// Scale returns a new transform which is t then a scale.
func (t *Transform) Scale(sx, sy float64) *Transform {
	// FIXME optimize matrix multiplication as we know the structure of the scale matrix
	return t.Then(Scale(sx, sy))
}

// Then returns a transform that is t then u.
func (t *Transform) Then(u *Transform) *Transform {
	return u.Multiply(t)
}

// Transform transforms a single vector.
func (t *Transform) Transform(v []float64) []float64 {
	return []float64{
		t.m[0]*v[0] + t.m[1]*v[1] + t.m[2],
		t.m[3]*v[0] + t.m[4]*v[1] + t.m[5],
	}
}

// TransformDirection transforms the direction v.
func (t *Transform) TransformDirection(v []float64) []float64 {
	return []float64{
		t.m[0]*v[0] + t.m[1]*v[1],
		t.m[3]*v[0] + t.m[4]*v[1],
	}
}

// TransformSlice transforms a slice of vectors.
func (t *Transform) TransformSlice(vs [][]float64) [][]float64 {
	result := make([][]float64, 0, len(vs))
	for _, v := range vs {
		result = append(result, t.Transform(v))
	}
	return result
}

// Translate returns a new transform which is t then a translate.
func (t *Transform) Translate(tx, ty float64) *Transform {
	// FIXME optimize matrix multiplication as we know the structure of the translate matrix
	return t.Then(Translate(tx, ty))
}
