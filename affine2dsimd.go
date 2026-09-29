//go:build goexperiment.simd

package affine2d

import "simd"

// TransformXYs transforms separate slices of X and Y coordinates.
func (t *Transform) TransformXYs(xs, ys []float64) ([]float64, []float64) {
	n := min(len(xs), len(ys))
	xResults := make([]float64, n)
	yResults := make([]float64, n)

	if n == 0 {
		return xResults, yResults
	}

	a := simd.BroadcastFloat64s(t.m[0])
	b := simd.BroadcastFloat64s(t.m[1])
	c := simd.BroadcastFloat64s(t.m[2])
	d := simd.BroadcastFloat64s(t.m[3])
	e := simd.BroadcastFloat64s(t.m[4])
	f := simd.BroadcastFloat64s(t.m[5])

	width := a.Len()
	i := 0

	for ; i+width <= n; i += width {
		x := simd.LoadFloat64s(xs[i:])
		y := simd.LoadFloat64s(ys[i:])

		transformedXs := x.Mul(a).MulAdd(y, b).Add(c)
		transformedYs := x.Mul(d).MulAdd(y, e).Add(f)

		transformedXs.Store(xResults[i:])
		transformedYs.Store(yResults[i:])
	}

	if i < n {
		x, _ := simd.LoadFloat64sPart(xs[i:])
		y, _ := simd.LoadFloat64sPart(ys[i:])

		transformedXs := x.Mul(a).MulAdd(y, b).Add(c)
		transformedYs := x.Mul(d).MulAdd(y, e).Add(f)

		transformedXs.StorePart(xResults[i:])
		transformedYs.StorePart(yResults[i:])
	}

	return xResults, yResults
}

// TransformXYsInPlace transforms separate slices of X and Y coordinates in place.
func (t *Transform) TransformXYsInPlace(xs, ys []float64) ([]float64, []float64) {
	n := min(len(xs), len(ys))
	if n == 0 {
		return xs, ys
	}

	a := simd.BroadcastFloat64s(t.m[0])
	b := simd.BroadcastFloat64s(t.m[1])
	c := simd.BroadcastFloat64s(t.m[2])
	d := simd.BroadcastFloat64s(t.m[3])
	e := simd.BroadcastFloat64s(t.m[4])
	f := simd.BroadcastFloat64s(t.m[5])

	width := a.Len()
	i := 0

	for ; i+width <= n; i += width {
		x := simd.LoadFloat64s(xs[i:])
		y := simd.LoadFloat64s(ys[i:])

		transformedXs := x.Mul(a).MulAdd(y, b).Add(c)
		transformedYs := x.Mul(d).MulAdd(y, e).Add(f)

		transformedXs.Store(xs[i:])
		transformedYs.Store(ys[i:])
	}

	if i < n {
		x, _ := simd.LoadFloat64sPart(xs[i:n])
		y, _ := simd.LoadFloat64sPart(ys[i:n])

		transformedXs := x.Mul(a).MulAdd(y, b).Add(c)
		transformedYs := x.Mul(d).MulAdd(y, e).Add(f)

		transformedXs.StorePart(xs[i:n])
		transformedYs.StorePart(ys[i:n])
	}

	return xs, ys
}
