// Package mandelbrot renders the Mandelbrot set as ASCII art by mapping a
// rectangular region of the complex plane onto a grid of characters, one per
// output pixel, based on how quickly each point escapes under iteration.
package mandelbrot

// The rendered region of the complex plane: real axis from -2.5 to 1.0,
// imaginary axis from -1.0 to 1.0. This range comfortably frames the whole
// Mandelbrot set.
const (
	realMin = -2.5
	realMax = 1.0
	imagMin = -1.0
	imagMax = 1.0
)

// defaultGradient maps iteration counts to characters, from lightest
// (quick escape) to darkest (never escapes, i.e. inside the set).
const defaultGradient = " .:-=+*#%@"

// escapeRadiusSquared is the squared magnitude threshold beyond which a
// point is considered to have escaped to infinity. The standard escape
// radius is 2, so the squared threshold is 4.
const escapeRadiusSquared = 4

// Render draws the Mandelbrot set as ASCII art, returning one string per
// row. Each row has exactly width characters, and there are exactly height
// rows.
//
// Each output pixel (col, row) is mapped to a point c in the complex plane,
// and escapeIterations counts how many iterations of z = z*z + c it takes
// for z to escape (or reports maxIter if it never does within the budget).
// That count is then mapped to a character: by default, using the gradient
// " .:-=+*#%@" from lightest (quick escape) to darkest (inside the set); if
// char is a non-empty string, points inside the set are drawn with char and
// all other points are drawn as spaces.
func Render(width, height, maxIter int, char string) []string {
	lines := make([]string, height)
	for row := 0; row < height; row++ {
		pixels := make([]byte, width)
		for col := 0; col < width; col++ {
			c := planePoint(col, row, width, height)
			iterations := escapeIterations(c, maxIter)
			pixels[col] = charForIterations(iterations, maxIter, defaultGradient, char)
		}
		lines[row] = string(pixels)
	}
	return lines
}

// planePoint maps an output pixel (col, row) within a width x height grid to
// the corresponding point in the rendered region of the complex plane,
// linearly interpolating real across [realMin, realMax] and imaginary across
// [imagMin, imagMax]. A grid dimension of 1 maps to the minimum of its
// corresponding range.
func planePoint(col, row, width, height int) complex128 {
	var realFrac, imagFrac float64
	if width > 1 {
		realFrac = float64(col) / float64(width-1)
	}
	if height > 1 {
		imagFrac = float64(row) / float64(height-1)
	}

	real := realMin + realFrac*(realMax-realMin)
	imag := imagMin + imagFrac*(imagMax-imagMin)
	return complex(real, imag)
}

// escapeIterations runs the Mandelbrot recurrence z = z*z + c, starting from
// z = 0, and returns the number of iterations before |z| exceeds the escape
// radius (2). If z has not escaped after maxIter iterations, it is treated
// as belonging to the set (or close enough to it), and maxIter is returned.
func escapeIterations(c complex128, maxIter int) int {
	var z complex128
	for i := 0; i < maxIter; i++ {
		if magnitudeSquared(z) > escapeRadiusSquared {
			return i
		}
		z = z*z + c
	}
	return maxIter
}

// magnitudeSquared returns |z|^2, avoiding the square root needed for |z|
// since comparing squared magnitudes against a squared threshold is
// equivalent and cheaper.
func magnitudeSquared(z complex128) float64 {
	return real(z)*real(z) + imag(z)*imag(z)
}

// charForIterations maps an escape-iteration count to the character used to
// render that pixel.
//
// If char is non-empty, it takes priority over the gradient: points that
// never escaped (iterations == maxIter, i.e. inside the set) are drawn with
// char's first byte, and every other point is drawn as a space. This is a
// deliberately simple inside/outside rule; finer-grained character-set
// configuration is expected to build on this later.
//
// Otherwise, iterations is scaled linearly onto the gradient string, so that
// 0 (immediate escape) maps to the first, lightest character and maxIter
// (never escapes) maps to the last, darkest character.
func charForIterations(iterations, maxIter int, gradient, char string) byte {
	if char != "" {
		if iterations >= maxIter {
			return char[0]
		}
		return ' '
	}

	if maxIter <= 0 {
		return gradient[0]
	}

	index := iterations * (len(gradient) - 1) / maxIter
	return gradient[index]
}
