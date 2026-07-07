package mandelbrot

import "testing"

// TestRender_DimensionsMatchRequest verifies the output has exactly the
// requested number of rows, and each row has exactly the requested number
// of columns, regardless of the character set in use.
func TestRender_DimensionsMatchRequest(t *testing.T) {
	const width, height = 40, 12

	lines := Render(width, height, 50, "")

	if len(lines) != height {
		t.Fatalf("Render(%d, %d, ...) returned %d lines, want %d", width, height, len(lines), height)
	}
	for i, line := range lines {
		if len(line) != width {
			t.Errorf("line %d has length %d, want %d", i, len(line), width)
		}
	}
}

// TestEscapeIterations_OriginNeverEscapes checks the known point c = 0+0i,
// which is the center of the main cardioid and never escapes: z stays at 0
// forever, so escapeIterations must run for the full maxIter budget.
func TestEscapeIterations_OriginNeverEscapes(t *testing.T) {
	const maxIter = 50

	got := escapeIterations(complex(0, 0), maxIter)

	if got != maxIter {
		t.Errorf("escapeIterations(0+0i, %d) = %d, want %d (should never escape)", maxIter, got, maxIter)
	}
}

// TestEscapeIterations_FarPointEscapesImmediately checks the known point
// c = 2+0i, which lies outside the escape radius after only a couple of
// iterations (|z| exceeds 2 by the second step), so it should report a very
// low iteration count.
func TestEscapeIterations_FarPointEscapesImmediately(t *testing.T) {
	const maxIter = 50

	got := escapeIterations(complex(2, 0), maxIter)

	if got >= maxIter {
		t.Fatalf("escapeIterations(2+0i, %d) = %d, want a small number well below maxIter", maxIter, got)
	}
	if got != 2 {
		t.Errorf("escapeIterations(2+0i, %d) = %d, want 2 (z=0 -> 2 -> 6, |6|>2 on the 3rd step)", maxIter, got)
	}
}

// TestCharForIterations_InsideSetUsesLastGradientChar confirms that a point
// which never escapes (iterations == maxIter) is rendered with the darkest,
// last character of the gradient.
func TestCharForIterations_InsideSetUsesLastGradientChar(t *testing.T) {
	const maxIter = 50
	gradient := " .:-=+*#%@"
	want := gradient[len(gradient)-1]

	got := charForIterations(maxIter, maxIter, gradient, "")

	if got != want {
		t.Errorf("charForIterations(%d, %d, gradient, \"\") = %q, want %q", maxIter, maxIter, got, want)
	}
}

// TestCharForIterations_QuickEscapeUsesFirstGradientChar confirms that a
// point escaping almost immediately (iterations far below maxIter) is
// rendered with the lightest, first character of the gradient.
func TestCharForIterations_QuickEscapeUsesFirstGradientChar(t *testing.T) {
	const maxIter = 50
	gradient := " .:-=+*#%@"
	want := gradient[0]

	got := charForIterations(2, maxIter, gradient, "")

	if got != want {
		t.Errorf("charForIterations(2, %d, gradient, \"\") = %q, want %q", maxIter, got, want)
	}
}

// TestCharForIterations_CustomCharUsesCharForInsideSet confirms that when a
// non-empty custom char is supplied, a point that never escapes (inside the
// set) is rendered with that char's first byte rather than a gradient
// character.
func TestCharForIterations_CustomCharUsesCharForInsideSet(t *testing.T) {
	const maxIter = 50
	gradient := " .:-=+*#%@"

	got := charForIterations(maxIter, maxIter, gradient, ".")

	if got != '.' {
		t.Errorf("charForIterations(%d, %d, gradient, \".\") = %q, want %q", maxIter, maxIter, got, '.')
	}
}

// TestCharForIterations_CustomCharUsesSpaceForOutsideSet confirms that when
// a non-empty custom char is supplied, a point that escapes before maxIter
// (outside the set) is rendered as a space, not the custom char or any
// gradient character.
func TestCharForIterations_CustomCharUsesSpaceForOutsideSet(t *testing.T) {
	const maxIter = 50
	gradient := " .:-=+*#%@"

	got := charForIterations(2, maxIter, gradient, ".")

	if got != ' ' {
		t.Errorf("charForIterations(2, %d, gradient, \".\") = %q, want %q (space)", maxIter, got, ' ')
	}
}

// TestRender_KnownPointInsideSetUsesMaxIterationChar exercises Render itself
// (not the internal helpers directly) to confirm that a known point inside
// the set is drawn with the darkest, last gradient character.
//
// planePoint maps pixel (col, row) in a width x height grid to
// real = -2.5 + (col/(width-1))*3.5, imag = -1.0 + (row/(height-1))*2.0.
// Choosing width=8, height=3 makes col=5, row=1 land exactly on c = 0+0i
// (col/(width-1) = 5/7, so real = -2.5 + (5/7)*3.5 = 0.0; row/(height-1) =
// 1/2, so imag = -1.0 + 0.5*2.0 = 0.0). c = 0+0i is the center of the main
// cardioid and never escapes, so it must render as the last gradient
// character.
func TestRender_KnownPointInsideSetUsesMaxIterationChar(t *testing.T) {
	const width, height, maxIter = 8, 3, 50
	const col, row = 5, 1
	gradient := " .:-=+*#%@"
	want := gradient[len(gradient)-1]

	lines := Render(width, height, maxIter, "")

	got := lines[row][col]
	if got != want {
		t.Errorf("Render(%d, %d, %d, \"\")[%d][%d] = %q, want %q (point 0+0i should never escape)",
			width, height, maxIter, row, col, got, want)
	}
}

// TestRender_KnownPointOutsideSetUsesMinIterationChar exercises Render
// itself to confirm that a known point outside the set escapes quickly and
// is drawn with the lightest, first gradient character.
//
// The plan's example point, 2+0i, lies outside the rendered region (real
// axis only spans -2.5 to 1.0), so it cannot be reproduced through Render's
// pixel grid. Instead this uses the grid's bottom-right corner
// (col=width-1, row=height-1), which planePoint maps to real = realMax =
// 1.0, imag = imagMax = 1.0, i.e. c = 1+1i. That point escapes after only 2
// iterations (0 -> 1+1i -> 2i+(1+1i) = 1+3i, |1+3i|^2 = 10 > 4), well below
// maxIter, so it must render as the first gradient character.
func TestRender_KnownPointOutsideSetUsesMinIterationChar(t *testing.T) {
	const width, height, maxIter = 8, 3, 50
	const col, row = width - 1, height - 1
	gradient := " .:-=+*#%@"
	want := gradient[0]

	lines := Render(width, height, maxIter, "")

	got := lines[row][col]
	if got != want {
		t.Errorf("Render(%d, %d, %d, \"\")[%d][%d] = %q, want %q (point 1+1i should escape quickly)",
			width, height, maxIter, row, col, got, want)
	}
}

// TestRender_CustomCharAppearsAtKnownInsidePoint reuses the known
// inside-the-set point from TestRender_KnownPointInsideSetUsesMaxIterationChar
// (col=5, row=1 in an 8x3 grid maps to c = 0+0i, which never escapes) to
// confirm that, with a custom char supplied, Render actually draws that char
// at the point -- not a gradient character.
func TestRender_CustomCharAppearsAtKnownInsidePoint(t *testing.T) {
	const width, height, maxIter = 8, 3, 50
	const col, row = 5, 1

	lines := Render(width, height, maxIter, ".")

	got := lines[row][col]
	if got != '.' {
		t.Errorf("Render(%d, %d, %d, \".\")[%d][%d] = %q, want %q (custom char at known inside-set point)",
			width, height, maxIter, row, col, got, '.')
	}
}

// TestRender_CustomCharUsesSpaceAtKnownOutsidePoint reuses the known
// outside-the-set point from TestRender_KnownPointOutsideSetUsesMinIterationChar
// (bottom-right corner, mapping to c = 1+1i, which escapes quickly) to
// confirm that, with a custom char supplied, Render draws a space there
// rather than the custom char or a gradient character.
func TestRender_CustomCharUsesSpaceAtKnownOutsidePoint(t *testing.T) {
	const width, height, maxIter = 8, 3, 50
	const col, row = width - 1, height - 1

	lines := Render(width, height, maxIter, ".")

	got := lines[row][col]
	if got != ' ' {
		t.Errorf("Render(%d, %d, %d, \".\")[%d][%d] = %q, want %q (space at known outside-set point)",
			width, height, maxIter, row, col, got, ' ')
	}
}

// TestRender_CustomCharOutputOnlyUsesCharAndSpace confirms that, with a
// custom char supplied, every rune in every output row is either that char
// or a space -- i.e. the gradient is not used at all, satisfying the
// "single character instead of gradient" behavior end-to-end.
func TestRender_CustomCharOutputOnlyUsesCharAndSpace(t *testing.T) {
	const width, height, maxIter = 40, 12, 50

	lines := Render(width, height, maxIter, ".")

	for r, line := range lines {
		for c, ch := range line {
			if ch != '.' && ch != ' ' {
				t.Fatalf("Render(%d, %d, %d, \".\")[%d][%d] = %q, want '.' or ' ' only (found gradient character)",
					width, height, maxIter, r, c, ch)
			}
		}
	}
}
