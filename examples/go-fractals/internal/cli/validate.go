package cli

import (
	"fmt"
	"unicode"
)

// requirePositive returns an error if value is not greater than 0, naming
// the offending flag and its value so the message is immediately actionable.
func requirePositive(name string, value int) error {
	if value <= 0 {
		return fmt.Errorf("%s must be greater than 0, got %d", name, value)
	}
	return nil
}

// requireNonNegative returns an error if value is negative, naming the
// offending flag and its value so the message is immediately actionable.
func requireNonNegative(name string, value int) error {
	if value < 0 {
		return fmt.Errorf("%s must be greater than or equal to 0, got %d", name, value)
	}
	return nil
}

// requireEmptyOrSingleASCIIChar returns an error unless value is empty or
// exactly one ASCII character.
//
// An empty value is always valid: for mandelbrot's --char flag, empty means
// "use the default gradient" rather than a specific character.
//
// A non-empty value must contain exactly one rune (matching the "exactly
// one character" wording used by sierpinski's --char validation, for a
// consistent error message across commands), and that rune must be ASCII.
// The ASCII restriction is stricter than sierpinski's: sierpinski stores its
// char as a rune end-to-end, so it renders multi-byte characters (e.g. "한"
// or "🔥") correctly. mandelbrot instead renders one byte per output cell
// (see internal/mandelbrot.charForIterations, which indexes char[0]), so a
// multi-byte rune would be silently truncated to its first byte and produce
// garbled output. Rejecting it here trades that asymmetry for a clear error
// instead of silent corruption.
func requireEmptyOrSingleASCIIChar(name, value string) error {
	if value == "" {
		return nil
	}

	runes := []rune(value)
	if len(runes) != 1 {
		return fmt.Errorf("%s must be exactly one character, got %q", name, value)
	}
	if runes[0] > unicode.MaxASCII {
		return fmt.Errorf("%s must be a single ASCII character, got %q", name, value)
	}
	return nil
}
