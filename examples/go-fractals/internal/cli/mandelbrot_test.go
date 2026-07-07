package cli

import (
	"strings"
	"testing"

	"github.com/my-poor-ai-test/fractals/internal/mandelbrot"
)

func TestMandelbrotCmd_DefaultFlagsMatchesRender(t *testing.T) {
	output, err := executeCommand("mandelbrot")
	if err != nil {
		t.Fatalf("unexpected error running mandelbrot: %v", err)
	}

	want := strings.Join(mandelbrot.Render(80, 24, 100, ""), "\n") + "\n"
	if output != want {
		t.Errorf("executeCommand(\"mandelbrot\") =\n%q\nwant\n%q", output, want)
	}
}

func TestMandelbrotCmd_CustomWidthHeightMatchesRender(t *testing.T) {
	output, err := executeCommand("mandelbrot", "--width", "40", "--height", "12")
	if err != nil {
		t.Fatalf("unexpected error running mandelbrot --width 40 --height 12: %v", err)
	}

	want := strings.Join(mandelbrot.Render(40, 12, 100, ""), "\n") + "\n"
	if output != want {
		t.Errorf("executeCommand(\"mandelbrot\", \"--width\", \"40\", \"--height\", \"12\") =\n%q\nwant\n%q", output, want)
	}
}

func TestMandelbrotCmd_CustomIterationsAndCharMatchesRender(t *testing.T) {
	output, err := executeCommand("mandelbrot", "--width", "20", "--height", "10", "--iterations", "25", "--char", "#")
	if err != nil {
		t.Fatalf("unexpected error running mandelbrot with custom iterations/char: %v", err)
	}

	want := strings.Join(mandelbrot.Render(20, 10, 25, "#"), "\n") + "\n"
	if output != want {
		t.Errorf("executeCommand with custom iterations/char =\n%q\nwant\n%q", output, want)
	}
}

func TestMandelbrotCmd_CustomCharUsesSingleCharacterNotGradient(t *testing.T) {
	output, err := executeCommand("mandelbrot", "--width", "40", "--height", "12", "--char", ".")
	if err != nil {
		t.Fatalf("unexpected error running mandelbrot --char '.': %v", err)
	}

	body := strings.TrimRight(output, "\n")
	if !strings.Contains(body, ".") {
		t.Errorf("expected output to contain custom char '.', got:\n%s", body)
	}

	for _, ch := range body {
		if ch != '.' && ch != ' ' && ch != '\n' {
			t.Errorf("expected output to use only '.' and ' ' (no gradient characters) with --char '.', found %q in:\n%s", ch, body)
			break
		}
	}
}

func TestMandelbrotCmd_SmallerSizeProducesFewerRows(t *testing.T) {
	output, err := executeCommand("mandelbrot", "--width", "40", "--height", "12")
	if err != nil {
		t.Fatalf("unexpected error running mandelbrot --width 40 --height 12: %v", err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) != 12 {
		t.Errorf("expected 12 rows of output, got %d:\n%s", len(lines), output)
	}
}

func TestMandelbrotCmd_NegativeWidthReturnsError(t *testing.T) {
	_, err := executeCommand("mandelbrot", "--width", "-1")
	if err == nil {
		t.Fatal("expected an error for --width -1, got nil")
	}
	if !strings.Contains(err.Error(), "width") {
		t.Errorf("expected error message to mention %q, got: %v", "width", err)
	}
}

func TestMandelbrotCmd_ZeroWidthReturnsError(t *testing.T) {
	_, err := executeCommand("mandelbrot", "--width", "0")
	if err == nil {
		t.Fatal("expected an error for --width 0, got nil")
	}
	if !strings.Contains(err.Error(), "width") {
		t.Errorf("expected error message to mention %q, got: %v", "width", err)
	}
}

func TestMandelbrotCmd_ZeroHeightReturnsError(t *testing.T) {
	_, err := executeCommand("mandelbrot", "--height", "0")
	if err == nil {
		t.Fatal("expected an error for --height 0, got nil")
	}
	if !strings.Contains(err.Error(), "height") {
		t.Errorf("expected error message to mention %q, got: %v", "height", err)
	}
}

func TestMandelbrotCmd_NegativeHeightReturnsError(t *testing.T) {
	_, err := executeCommand("mandelbrot", "--height", "-3")
	if err == nil {
		t.Fatal("expected an error for --height -3, got nil")
	}
	if !strings.Contains(err.Error(), "height") {
		t.Errorf("expected error message to mention %q, got: %v", "height", err)
	}
}

func TestMandelbrotCmd_ZeroIterationsReturnsError(t *testing.T) {
	_, err := executeCommand("mandelbrot", "--iterations", "0")
	if err == nil {
		t.Fatal("expected an error for --iterations 0, got nil")
	}
	if !strings.Contains(err.Error(), "iterations") {
		t.Errorf("expected error message to mention %q, got: %v", "iterations", err)
	}
}

func TestMandelbrotCmd_NegativeIterationsReturnsError(t *testing.T) {
	_, err := executeCommand("mandelbrot", "--iterations", "-10")
	if err == nil {
		t.Fatal("expected an error for --iterations -10, got nil")
	}
	if !strings.Contains(err.Error(), "iterations") {
		t.Errorf("expected error message to mention %q, got: %v", "iterations", err)
	}
}

func TestMandelbrotCmd_EmptyCharIsValid(t *testing.T) {
	_, err := executeCommand("mandelbrot", "--width", "10", "--height", "5", "--char", "")
	if err != nil {
		t.Fatalf("expected empty --char to be valid (use gradient), got error: %v", err)
	}
}

func TestMandelbrotCmd_MultiCharReturnsError(t *testing.T) {
	_, err := executeCommand("mandelbrot", "--width", "10", "--height", "5", "--char", "ab")
	if err == nil {
		t.Fatal("expected an error for --char 'ab', got nil")
	}
	if !strings.Contains(err.Error(), "char must be exactly one character") {
		t.Errorf("expected error message to mention %q, got: %v", "char must be exactly one character", err)
	}
	if !strings.Contains(err.Error(), `"ab"`) {
		t.Errorf("expected error message to quote the offending value %q, got: %v", "ab", err)
	}
}

func TestMandelbrotCmd_MultiByteSingleRuneCharReturnsError(t *testing.T) {
	// "한" is a single rune but three UTF-8 bytes. mandelbrot renders one
	// byte per output cell (internal/mandelbrot.charForIterations takes
	// char[0]), so a multi-byte rune would be silently truncated to its
	// first byte and produce garbled output. Reject it with a clear error
	// instead of corrupting the rendered output.
	_, err := executeCommand("mandelbrot", "--width", "10", "--height", "5", "--char", "한")
	if err == nil {
		t.Fatal("expected an error for --char '한' (multi-byte rune), got nil")
	}
	if !strings.Contains(err.Error(), "ASCII") {
		t.Errorf("expected error message to mention %q, got: %v", "ASCII", err)
	}
}

func TestMandelbrotCmd_SingleASCIICharIsValid(t *testing.T) {
	_, err := executeCommand("mandelbrot", "--width", "10", "--height", "5", "--char", "#")
	if err != nil {
		t.Fatalf("expected single ASCII --char '#' to be valid, got error: %v", err)
	}
}

func TestMandelbrotCmd_HelpDocumentsFlags(t *testing.T) {
	output, err := executeCommand("mandelbrot", "--help")
	if err != nil {
		t.Fatalf("unexpected error running mandelbrot --help: %v", err)
	}

	for _, flag := range []string{"--width", "--height", "--iterations", "--char"} {
		if !strings.Contains(output, flag) {
			t.Errorf("expected help output to document flag %q, got:\n%s", flag, output)
		}
	}
}
