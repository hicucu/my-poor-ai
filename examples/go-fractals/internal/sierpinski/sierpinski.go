// Package sierpinski generates ASCII-art Sierpinski triangles using
// recursive midpoint subdivision.
package sierpinski

import "strings"

// Generate returns the lines of an ASCII-art Sierpinski triangle with the
// given base size (in characters), recursion depth, and fill character.
// Each returned string is one row of the triangle; rows are padded on the
// left with spaces so the triangle appears centered, with no trailing
// spaces.
//
// A depth of 0 (or less) produces a fully filled triangle, with no holes
// removed. A size of 1 always produces a single filled character,
// regardless of depth. Any size (not just powers of two) produces a
// triangle with no blank rows; when size does not split evenly, the top
// sub-triangle gets the extra row so the three sub-triangles still cover
// every row exactly once.
func Generate(size, depth int, char rune) []string {
	if size < 1 {
		return []string{}
	}

	width := 2*size - 1
	grid := make([][]rune, size)
	for r := range grid {
		grid[r] = make([]rune, width)
		for c := range grid[r] {
			grid[r][c] = ' '
		}
	}

	fillTriangle(grid, 0, size-1, size, depth, char)

	lines := make([]string, size)
	for r, row := range grid {
		lines[r] = strings.TrimRight(string(row), " ")
	}
	return lines
}

// fillTriangle fills the triangle whose apex sits at (topRow, apexCol) in
// grid and whose base spans the given size (number of rows). If depth has
// been exhausted, or the triangle can no longer be subdivided (size <= 1),
// the whole triangle is filled solid. Otherwise it is split into three
// corner sub-triangles, each half the size, leaving the inverted middle
// triangle empty, and each corner is filled recursively one depth lower.
func fillTriangle(grid [][]rune, topRow, apexCol, size, depth int, char rune) {
	if depth <= 0 || size <= 1 {
		for r := 0; r < size; r++ {
			row := topRow + r
			for c := apexCol - r; c <= apexCol+r; c++ {
				grid[row][c] = char
			}
		}
		return
	}

	// Split size into a top sub-triangle and two bottom sub-triangles that
	// together cover every row of the parent triangle exactly once. When
	// size is odd, the top sub-triangle absorbs the extra row (topSize =
	// ceil(size/2)) so the bottom two (bottomSize = floor(size/2)) still
	// start immediately below it, leaving no uncovered rows.
	topSize := (size + 1) / 2
	bottomSize := size - topSize
	fillTriangle(grid, topRow, apexCol, topSize, depth-1, char)
	fillTriangle(grid, topRow+topSize, apexCol-topSize, bottomSize, depth-1, char)
	fillTriangle(grid, topRow+topSize, apexCol+topSize, bottomSize, depth-1, char)
}
