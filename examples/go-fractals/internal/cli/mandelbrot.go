package cli

import (
	"strings"

	"github.com/my-poor-ai-test/fractals/internal/mandelbrot"
	"github.com/spf13/cobra"
)

// newMandelbrotCmd constructs the "mandelbrot" subcommand, which renders an
// ASCII-art Mandelbrot set to stdout via internal/mandelbrot.Render.
func newMandelbrotCmd() *cobra.Command {
	var width int
	var height int
	var iterations int
	var char string

	cmd := &cobra.Command{
		Use:   "mandelbrot",
		Short: "Render a Mandelbrot set fractal",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validation errors below are already clear on their own; dumping
			// the full usage block after them just buries the message. Only
			// silence usage once we know we're past flag parsing and into
			// semantic validation, so parsing errors (e.g. a mistyped flag)
			// still print usage as Cobra normally would.
			if err := requirePositive("width", width); err != nil {
				cmd.SilenceUsage = true
				return err
			}
			if err := requirePositive("height", height); err != nil {
				cmd.SilenceUsage = true
				return err
			}
			if err := requirePositive("iterations", iterations); err != nil {
				cmd.SilenceUsage = true
				return err
			}
			if err := requireEmptyOrSingleASCIIChar("char", char); err != nil {
				cmd.SilenceUsage = true
				return err
			}

			lines := mandelbrot.Render(width, height, iterations, char)
			cmd.Println(strings.Join(lines, "\n"))
			return nil
		},
	}

	cmd.Flags().IntVar(&width, "width", 80, "Output width in characters")
	cmd.Flags().IntVar(&height, "height", 24, "Output height in characters")
	cmd.Flags().IntVar(&iterations, "iterations", 100, "Maximum iterations for escape calculation")
	cmd.Flags().StringVar(&char, "char", "", "Single ASCII character to use instead of the default gradient")

	return cmd
}
