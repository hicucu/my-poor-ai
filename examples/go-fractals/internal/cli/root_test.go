package cli

import (
	"bytes"
	"strings"
	"testing"
)

// executeCommand runs the root command with the given args and returns
// combined stdout/stderr output.
func executeCommand(args ...string) (string, error) {
	cmd := NewRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return buf.String(), err
}

func TestRootCmd_HelpListsSubcommands(t *testing.T) {
	output, err := executeCommand("--help")
	if err != nil {
		t.Fatalf("unexpected error running --help: %v", err)
	}

	availableSection := extractSection(t, output, "Available Commands:", "Flags:")

	if !strings.Contains(availableSection, "sierpinski") {
		t.Errorf("expected %q to be listed as an available command, got:\n%s", "sierpinski", output)
	}
	if !strings.Contains(availableSection, "mandelbrot") {
		t.Errorf("expected %q to be listed as an available command, got:\n%s", "mandelbrot", output)
	}
}

// extractSection returns the substring of output between startMarker
// (inclusive) and the next occurrence of endMarker, failing the test if
// either marker is missing.
func extractSection(t *testing.T, output, startMarker, endMarker string) string {
	t.Helper()

	startIdx := strings.Index(output, startMarker)
	if startIdx == -1 {
		t.Fatalf("expected output to contain %q, got:\n%s", startMarker, output)
	}
	rest := output[startIdx:]

	endIdx := strings.Index(rest, endMarker)
	if endIdx == -1 {
		t.Fatalf("expected output to contain %q after %q, got:\n%s", endMarker, startMarker, output)
	}
	return rest[:endIdx]
}

func TestRootCmd_NoArgsShowsHelp(t *testing.T) {
	output, err := executeCommand()
	if err != nil {
		t.Fatalf("unexpected error running with no args: %v", err)
	}

	if !strings.Contains(output, "Usage:") {
		t.Errorf("expected help output containing %q when run with no args, got:\n%s", "Usage:", output)
	}
	if !strings.Contains(output, "sierpinski") {
		t.Errorf("expected no-args output to list %q, got:\n%s", "sierpinski", output)
	}
	if !strings.Contains(output, "mandelbrot") {
		t.Errorf("expected no-args output to list %q, got:\n%s", "mandelbrot", output)
	}
}
