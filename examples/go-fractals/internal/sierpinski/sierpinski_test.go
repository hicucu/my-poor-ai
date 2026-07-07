package sierpinski

import (
	"reflect"
	"strings"
	"testing"
)

func TestGenerate_SmallTriangleWithSubdivision(t *testing.T) {
	got := Generate(4, 2, '*')
	want := []string{
		"   *",
		"  * *",
		" *   *",
		"* * * *",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Generate(4, 2, '*') =\n%#v\nwant\n%#v", got, want)
	}
}

func TestGenerate_SizeOneReturnsSingleCharacter(t *testing.T) {
	got := Generate(1, 5, '*')
	want := []string{"*"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Generate(1, 5, '*') = %#v, want %#v", got, want)
	}
}

func TestGenerate_DepthZeroReturnsFilledTriangle(t *testing.T) {
	got := Generate(4, 0, '*')
	want := []string{
		"   *",
		"  ***",
		" *****",
		"*******",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Generate(4, 0, '*') = %#v, want %#v", got, want)
	}
}

// TestGenerate_NonPowerOfTwoSizeHasNoBlankRows is a regression test for a bug
// where fillTriangle used size/2 (plain integer division) to split each
// sub-triangle. For sizes that are not a power of two, that dropped rows
// from the recursion entirely, leaving some rows completely blank (e.g.
// Generate(3, 3, '*') produced an empty string for row 2, and Generate(6, 3,
// '*') / Generate(10, 3, '*') each produced two empty rows). Every row of
// the triangle must contain at least one filled character.
func TestGenerate_NonPowerOfTwoSizeHasNoBlankRows(t *testing.T) {
	cases := []struct {
		size  int
		depth int
	}{
		{3, 3},
		{5, 3},
		{6, 3},
		{10, 3},
		{20, 4},
	}

	for _, c := range cases {
		got := Generate(c.size, c.depth, '*')

		if len(got) != c.size {
			t.Fatalf("Generate(%d, %d, '*') returned %d rows, want %d", c.size, c.depth, len(got), c.size)
		}

		for i, row := range got {
			if strings.Trim(row, " ") == "" {
				t.Errorf("Generate(%d, %d, '*') row %d is blank, want at least one %q", c.size, c.depth, i, '*')
			}
		}
	}
}

// TestGenerate_OddSizeExactShape locks in the exact shape produced for a
// representative odd, non-power-of-two size, so the subdivision logic for
// uneven splits (top sub-triangle absorbs the extra row) doesn't silently
// regress.
func TestGenerate_OddSizeExactShape(t *testing.T) {
	got := Generate(3, 3, '*')
	want := []string{
		"  *",
		" * *",
		"*   *",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Generate(3, 3, '*') =\n%#v\nwant\n%#v", got, want)
	}
}
