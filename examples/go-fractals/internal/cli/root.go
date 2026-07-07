// Package cli wires up the Cobra command tree for the fractals CLI.
package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCmd constructs the root "fractals" command along with its
// subcommands. Running the root command with no arguments (or with
// --help) prints usage information listing the available subcommands.
//
// Both the sierpinski and mandelbrot subcommands are fully implemented (see
// internal/cli/sierpinski.go and internal/cli/mandelbrot.go).
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "fractals",
		Short: "Generate ASCII fractals",
		Long:  "fractals is a CLI tool that renders ASCII fractals such as Sierpinski triangles and Mandelbrot sets.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(newSierpinskiCmd())
	rootCmd.AddCommand(newMandelbrotCmd())

	return rootCmd
}
