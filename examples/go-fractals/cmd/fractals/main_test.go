package main

// Black-box integration tests for the compiled fractals binary.
//
// Unlike internal/cli's tests, which execute cobra.Command objects
// in-process and inspect the returned error, these tests build the actual
// binary and run it as a subprocess via os/exec. That is the only way to
// observe the real OS process exit code produced by the os.Exit(1) call in
// main() below, which no other test in the suite exercises.

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// binPath holds the path to the fractals binary built once in TestMain and
// shared by every test in this file, so the binary is compiled a single
// time no matter how many subtests run.
var binPath string

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "fractals-integration-*")
	if err != nil {
		panic(err)
	}

	binPath = filepath.Join(tmpDir, "fractals")
	build := exec.Command("go", "build", "-o", binPath, ".")
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		panic("failed to build fractals binary for integration tests: " + err.Error())
	}

	code := m.Run()
	os.RemoveAll(tmpDir)
	os.Exit(code)
}

// runBinary executes the built fractals binary with the given args and
// returns its combined stdout, stderr, and exit code.
func runBinary(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()

	cmd := exec.Command(binPath, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()
	if err == nil {
		return outBuf.String(), errBuf.String(), 0
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("failed to run fractals binary: %v", err)
	}
	return outBuf.String(), errBuf.String(), exitErr.ExitCode()
}

func TestIntegration_HelpShowsUsage(t *testing.T) {
	stdout, _, exitCode := runBinary(t, "--help")

	if exitCode != 0 {
		t.Errorf("expected exit code 0 for --help, got %d", exitCode)
	}
	if !strings.Contains(stdout, "Usage:") {
		t.Errorf("expected --help output to contain %q, got:\n%s", "Usage:", stdout)
	}
	if !strings.Contains(stdout, "sierpinski") {
		t.Errorf("expected --help output to list %q, got:\n%s", "sierpinski", stdout)
	}
	if !strings.Contains(stdout, "mandelbrot") {
		t.Errorf("expected --help output to list %q, got:\n%s", "mandelbrot", stdout)
	}
}

func TestIntegration_SierpinskiSucceeds(t *testing.T) {
	stdout, stderr, exitCode := runBinary(t, "sierpinski", "--size", "16", "--depth", "3")

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr: %s)", exitCode, stderr)
	}
	if !strings.Contains(stdout, "*") {
		t.Errorf("expected sierpinski output to contain the default fill character %q, got:\n%s", "*", stdout)
	}
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	if len(lines) < 2 {
		t.Errorf("expected sierpinski output to have multiple lines, got %d: %q", len(lines), stdout)
	}
}

func TestIntegration_SierpinskiCustomChar(t *testing.T) {
	stdout, stderr, exitCode := runBinary(t, "sierpinski", "--size", "8", "--depth", "2", "--char", "#")

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr: %s)", exitCode, stderr)
	}
	if !strings.Contains(stdout, "#") {
		t.Errorf("expected sierpinski output to contain custom char %q, got:\n%s", "#", stdout)
	}
}

func TestIntegration_MandelbrotSucceeds(t *testing.T) {
	stdout, stderr, exitCode := runBinary(t, "mandelbrot", "--width", "40", "--height", "12", "--iterations", "20")

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d (stderr: %s)", exitCode, stderr)
	}
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	if len(lines) != 12 {
		t.Errorf("expected mandelbrot output to have 12 lines, got %d:\n%s", len(lines), stdout)
	}
	if strings.TrimSpace(stdout) == "" {
		t.Errorf("expected mandelbrot output to be non-empty")
	}
}

func TestIntegration_SierpinskiInvalidSizeFails(t *testing.T) {
	stdout, stderr, exitCode := runBinary(t, "sierpinski", "--size", "0")

	if exitCode == 0 {
		t.Fatalf("expected non-zero exit code for --size 0, got 0 (stdout: %s)", stdout)
	}
	if !strings.Contains(stderr, "size") {
		t.Errorf("expected stderr to mention the invalid %q flag, got:\n%s", "size", stderr)
	}
}

func TestIntegration_MandelbrotInvalidWidthFails(t *testing.T) {
	stdout, stderr, exitCode := runBinary(t, "mandelbrot", "--width", "-1")

	if exitCode == 0 {
		t.Fatalf("expected non-zero exit code for --width -1, got 0 (stdout: %s)", stdout)
	}
	if !strings.Contains(stderr, "width") {
		t.Errorf("expected stderr to mention the invalid %q flag, got:\n%s", "width", stderr)
	}
}

func TestIntegration_UnknownCommandFails(t *testing.T) {
	_, stderr, exitCode := runBinary(t, "not-a-real-command")

	if exitCode == 0 {
		t.Fatalf("expected non-zero exit code for an unknown command, got 0")
	}
	if !strings.Contains(stderr, "unknown command") {
		t.Errorf("expected stderr to mention %q, got:\n%s", "unknown command", stderr)
	}
}
