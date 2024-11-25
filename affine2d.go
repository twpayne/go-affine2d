// Package affine2d implements 2D affine transformations.
package affine2d

import "math"

// A Transform is a 2D affine transform.
type Transform struct {
	m [6]float64
}

// Delta returns the transform that transforms the origin to origin and the unit
// X vector to unitX.
func Delta(origin, unitX []float64) *Transform {
	dx := unitX[0] - origin[0]
	dy := unitX[1] - origin[1]
	return &Transform{
		m: [6]float64{
			dx, -dy, origin[0],
			dy, dx, origin[1],
		},
	}
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

// NewTransform returns a new transform with the given coefficients.
func NewTransform(m [6]float64) *Transform {
	return &Transform{
		m: m,
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

// Float64Array returns the coefficients of t's matrix in row order.
func (t *Transform) Float64Array() [6]float64 {
	return t.m
}

// Float64Slice returns the coefficients of t's matrix in row order. The returned
// slice must not be modified.
func (t *Transform) Float64Slice() []float64 {
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
func (t *Transform) Transform(p []float64) []float64 {
	x, y := t.TransformXY(p[0], p[1])
	return []float64{x, y}
}

// TransformInPlace transforms a single vector in place.
func (t *Transform) TransformInPlace(p []float64) []float64 {
	p[0], p[1] = t.TransformXY(p[0], p[1])
	return p
}

// TransformDirection transforms the direction v.
func (t *Transform) TransformDirection(v []float64) []float64 {
	return []float64{
		t.m[0]*v[0] + t.m[1]*v[1],
		t.m[3]*v[0] + t.m[4]*v[1],
	}
}

// TransformSlice transforms a slice of vectors.
func (t *Transform) TransformSlice(ps [][]float64) [][]float64 {
	result := make([][]float64, len(ps))
	for i, p := range ps {
		result[i] = t.Transform(p)
	}
	return result
}

// TransformSliceInPlace transforms a slice of vectors in place.
func (t *Transform) TransformSliceInPlace(ps [][]float64) [][]float64 {
	for i := range ps {
		t.TransformInPlace(ps[i])
	}
	return ps
}

// TransformXY transforms an X-Y coordinate.
func (t *Transform) TransformXY(x, y float64) (float64, float64) {
	return t.m[0]*x + t.m[1]*y + t.m[2], t.m[3]*x + t.m[4]*y + t.m[5]
}

// Translate returns a new transform which is t then a translate.
func (t *Transform) Translate(tx, ty float64) *Transform {
	// FIXME optimize matrix multiplication as we know the structure of the translate matrix
	return t.Then(Translate(tx, ty))
}
