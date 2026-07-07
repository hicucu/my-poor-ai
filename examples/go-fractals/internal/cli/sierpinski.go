package cli

import (
	"fmt"
	"strings"

	"github.com/my-poor-ai-test/fractals/internal/sierpinski"
	"github.com/spf13/cobra"
)

// newSierpinskiCmd constructs the "sierpinski" subcommand, which renders an
// ASCII-art Sierpinski triangle to stdout via internal/sierpinski.Generate.
func newSierpinskiCmd() *cobra.Command {
	var size int
	var depth int
	var char string

	cmd := &cobra.Command{
		Use:   "sierpinski",
		Short: "Render a Sierpinski triangle fractal",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validation errors below are already clear on their own; dumping
			// the full usage block after them just buries the message. Only
			// silence usage once we know we're past flag parsing and into
			// semantic validation, so parsing errors (e.g. a mistyped flag)
			// still print usage as Cobra normally would.
			if err := requirePositive("size", size); err != nil {
				cmd.SilenceUsage = true
				return err
			}
			if err := requireNonNegative("depth", depth); err != nil {
				cmd.SilenceUsage = true
				return err
			}
			runes := []rune(char)
			if len(runes) != 1 {
				cmd.SilenceUsage = true
				return fmt.Errorf("char must be exactly one character, got %q", char)
			}

			lines := sierpinski.Generate(size, depth, runes[0])
			cmd.Println(strings.Join(lines, "\n"))
			return nil
		},
	}

	cmd.Flags().IntVar(&size, "size", 32, "Width of the triangle base in characters")
	cmd.Flags().IntVar(&depth, "depth", 5, "Recursion depth")
	cmd.Flags().StringVar(&char, "char", "*", "Character to use for filled points")

	return cmd
}
