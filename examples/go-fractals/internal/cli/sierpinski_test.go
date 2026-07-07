package cli

import (
	"strings"
	"testing"

	"github.com/my-poor-ai-test/fractals/internal/sierpinski"
)

func TestSierpinskiCmd_DefaultFlagsMatchesGenerate(t *testing.T) {
	output, err := executeCommand("sierpinski")
	if err != nil {
		t.Fatalf("unexpected error running sierpinski: %v", err)
	}

	want := strings.Join(sierpinski.Generate(32, 5, '*'), "\n") + "\n"
	if output != want {
		t.Errorf("executeCommand(\"sierpinski\") =\n%q\nwant\n%q", output, want)
	}
}

func TestSierpinskiCmd_CustomSizeAndDepthMatchesGenerate(t *testing.T) {
	output, err := executeCommand("sierpinski", "--size", "4", "--depth", "2")
	if err != nil {
		t.Fatalf("unexpected error running sierpinski --size 4 --depth 2: %v", err)
	}

	want := strings.Join(sierpinski.Generate(4, 2, '*'), "\n") + "\n"
	if output != want {
		t.Errorf("executeCommand(\"sierpinski\", \"--size\", \"4\", \"--depth\", \"2\") =\n%q\nwant\n%q", output, want)
	}
}

func TestSierpinskiCmd_CustomCharMatchesGenerate(t *testing.T) {
	output, err := executeCommand("sierpinski", "--size", "4", "--depth", "2", "--char", "#")
	if err != nil {
		t.Fatalf("unexpected error running sierpinski with --char '#': %v", err)
	}

	want := strings.Join(sierpinski.Generate(4, 2, '#'), "\n") + "\n"
	if output != want {
		t.Errorf("executeCommand with --char '#' =\n%q\nwant\n%q", output, want)
	}
}

func TestSierpinskiCmd_CustomCharAppearsAndDefaultCharDoesNotAppear(t *testing.T) {
	output, err := executeCommand("sierpinski", "--size", "8", "--depth", "3", "--char", "#")
	if err != nil {
		t.Fatalf("unexpected error running sierpinski with --char '#': %v", err)
	}

	body := strings.TrimRight(output, "\n")
	if !strings.Contains(body, "#") {
		t.Errorf("expected output to contain custom char '#', got:\n%s", body)
	}
	if strings.Contains(body, "*") {
		t.Errorf("expected output to not contain the default char '*' when --char '#' is given, got:\n%s", body)
	}
}

func TestSierpinskiCmd_SmallerSizeProducesFewerRows(t *testing.T) {
	output, err := executeCommand("sierpinski", "--size", "16", "--depth", "3")
	if err != nil {
		t.Fatalf("unexpected error running sierpinski --size 16 --depth 3: %v", err)
	}

	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	if len(lines) != 16 {
		t.Errorf("expected 16 rows of output, got %d:\n%s", len(lines), output)
	}
}

func TestSierpinskiCmd_ZeroSizeReturnsError(t *testing.T) {
	_, err := executeCommand("sierpinski", "--size", "0")
	if err == nil {
		t.Fatal("expected an error for --size 0, got nil")
	}
	if !strings.Contains(err.Error(), "size") {
		t.Errorf("expected error message to mention %q, got: %v", "size", err)
	}
}

func TestSierpinskiCmd_NegativeSizeReturnsError(t *testing.T) {
	_, err := executeCommand("sierpinski", "--size", "-5")
	if err == nil {
		t.Fatal("expected an error for --size -5, got nil")
	}
	if !strings.Contains(err.Error(), "size") {
		t.Errorf("expected error message to mention %q, got: %v", "size", err)
	}
}

func TestSierpinskiCmd_NegativeDepthReturnsError(t *testing.T) {
	_, err := executeCommand("sierpinski", "--depth", "-1")
	if err == nil {
		t.Fatal("expected an error for --depth -1, got nil")
	}
	if !strings.Contains(err.Error(), "depth") {
		t.Errorf("expected error message to mention %q, got: %v", "depth", err)
	}
}

func TestSierpinskiCmd_EmptyCharReturnsError(t *testing.T) {
	_, err := executeCommand("sierpinski", "--char", "")
	if err == nil {
		t.Fatal("expected an error for --char '', got nil")
	}
	if !strings.Contains(err.Error(), "char") {
		t.Errorf("expected error message to mention %q, got: %v", "char", err)
	}
}

func TestSierpinskiCmd_MultiCharReturnsError(t *testing.T) {
	_, err := executeCommand("sierpinski", "--char", "ab")
	if err == nil {
		t.Fatal("expected an error for --char 'ab', got nil")
	}
	if !strings.Contains(err.Error(), "char") {
		t.Errorf("expected error message to mention %q, got: %v", "char", err)
	}
}

func TestSierpinskiCmd_ValidZeroDepthSucceeds(t *testing.T) {
	_, err := executeCommand("sierpinski", "--depth", "0")
	if err != nil {
		t.Fatalf("expected --depth 0 to be valid, got error: %v", err)
	}
}

func TestSierpinskiCmd_HelpDocumentsFlags(t *testing.T) {
	output, err := executeCommand("sierpinski", "--help")
	if err != nil {
		t.Fatalf("unexpected error running sierpinski --help: %v", err)
	}

	for _, flag := range []string{"--size", "--depth", "--char"} {
		if !strings.Contains(output, flag) {
			t.Errorf("expected help output to document flag %q, got:\n%s", flag, output)
		}
	}
}
