//go:build !goexperiment.simd

package affine2d

// TransformXYs transforms separate slices of X and Y coordinates.
func (t *Transform) TransformXYs(xs, ys []float64) ([]float64, []float64) {
	n := min(len(xs), len(ys))
	transformedXs, transformedYs := make([]float64, n), make([]float64, n)
	for i := range n {
		transformedXs[i], transformedYs[i] = t.TransformXY(xs[i], ys[i])
	}
	return transformedXs, transformedYs
}

// TransformXYsInPlace transforms separate slices of X and Y coordinates in place.
func (t *Transform) TransformXYsInPlace(xs, ys []float64) ([]float64, []float64) {
	for i := range min(len(xs), len(ys)) {
		xs[i], ys[i] = t.TransformXY(xs[i], ys[i])
	}
	return xs, ys
}
