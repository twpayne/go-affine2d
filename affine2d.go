// Package affine2d implements 2D affine transformations.
package affine2d

import "math"

// A Transformation is a 2D affine transformation.
type Transformation struct {
	m [6]float64
}

// Identity returns a new identity transformation.
func Identity() *Transformation {
	return &Transformation{
		m: [6]float64{
			1, 0, 0,
			0, 1, 0,
		},
	}
}

// Rotate returns a new rotate transformation.
func Rotate(theta float64) *Transformation {
	cosTheta := math.Cos(theta)
	sinTheta := math.Sin(theta)
	return &Transformation{
		m: [6]float64{
			cosTheta, -sinTheta, 0,
			sinTheta, cosTheta, 0,
		},
	}
}

// Scale returns a new scale transformation.
func Scale(sx, sy float64) *Transformation {
	return &Transformation{
		m: [6]float64{
			sx, 0, 0,
			0, sy, 0,
		},
	}
}

// Shear returns a new shear transformation.
func Shear(sx, sy float64) *Transformation {
	return &Transformation{
		m: [6]float64{
			1, sx, 0,
			sy, 1, 0,
		},
	}
}

// Translate returns a new translate transformation.
func Translate(tx, ty float64) *Transformation {
	return &Transformation{
		m: [6]float64{
			1, 0, tx,
			0, 1, ty,
		},
	}
}

// Float64s returns the coefficients of t's matrix in row order. The returned
// slice must not be modified.
func (t *Transformation) Float64s() []float64 {
	return t.m[:]
}

func (t *Transformation) Inverse() *Transformation {
	d := t.m[0]*t.m[4] - t.m[1]*t.m[3]
	return &Transformation{
		m: [6]float64{
			t.m[4] / d, -t.m[1] / d, (t.m[1]*t.m[5] - t.m[2]*t.m[4]) / d,
			-t.m[3] / d, t.m[0] / d, (t.m[2]*t.m[3] - t.m[0]*t.m[5]) / d,
		},
	}
}

// Multiply returns a new transformation which is u then t.
func (t *Transformation) Multiply(u *Transformation) *Transformation {
	return &Transformation{
		m: [6]float64{
			t.m[0]*u.m[0] + t.m[1]*u.m[3], t.m[0]*u.m[1] + t.m[1]*u.m[4], t.m[0]*u.m[2] + t.m[1]*u.m[5] + t.m[2],
			t.m[3]*u.m[0] + t.m[4]*u.m[3], t.m[3]*u.m[1] + t.m[4]*u.m[4], t.m[3]*u.m[2] + t.m[4]*u.m[5] + t.m[5],
		},
	}
}

// Rotate returns a new transformation which is t then a rotation.
func (t *Transformation) Rotate(theta float64) *Transformation {
	return t.Then(Rotate(theta))
}

// Scale returns a new transformation which is t then a scale.
func (t *Transformation) Scale(sx, sy float64) *Transformation {
	return t.Then(Scale(sx, sy))
}

// Then returns a transformation that is t then u.
func (t *Transformation) Then(u *Transformation) *Transformation {
	return u.Multiply(t)
}

// Transform transforms a single vector.
func (t *Transformation) Transform(v []float64) []float64 {
	return []float64{
		t.m[0]*v[0] + t.m[1]*v[1] + t.m[2],
		t.m[3]*v[0] + t.m[4]*v[1] + t.m[5],
	}
}

// TransformDirection transforms the direction v.
func (t *Transformation) TransformDirection(v []float64) []float64 {
	return []float64{
		t.m[0]*v[0] + t.m[1]*v[1],
		t.m[3]*v[0] + t.m[4]*v[1],
	}
}

// TransformSlice transforms a slice of vectors.
func (t *Transformation) TransformSlice(vs [][]float64) [][]float64 {
	result := make([][]float64, 0, len(vs))
	for _, v := range vs {
		result = append(result, t.Transform(v))
	}
	return result
}

// Translate returns a new transformation which is t then a translate.
func (t *Transformation) Translate(tx, ty float64) *Transformation {
	return t.Then(Translate(tx, ty))
}
