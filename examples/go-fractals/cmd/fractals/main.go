package main

import (
	"os"

	"github.com/my-poor-ai-test/fractals/internal/cli"
)

func main() {
	rootCmd := cli.NewRootCmd()
	// Cobra's cmd.Println/Printf (used by the sierpinski/mandelbrot
	// subcommands to print their rendered output) fall back to stderr when
	// no output writer has been explicitly set. Without these two calls,
	// the fractal art printed via cmd.Println would go to stderr instead
	// of stdout, silently breaking redirection/piping (e.g. `fractals
	// sierpinski > out.txt` would produce an empty file) and violating
	// design.md's "printed to stdout" requirement. Error output (Cobra's
	// own usage/error printing plus PrintErr/PrintErrln) already defaults
	// to stderr; SetErr is set here explicitly for symmetry/clarity.
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
