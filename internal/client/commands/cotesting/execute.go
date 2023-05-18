// Package cotesting contains a functions for testing cobra commands.
package cotesting

import "github.com/spf13/cobra"

// ExecuteCommandC takes a root *cobra.Command and a variable number of arguments as strings,
// sets the arguments on the root command, and executes it. Returns an error if the execution fails.
func ExecuteCommandC(
	root *cobra.Command,
	args ...string,
) error {
	root.SetArgs(args)
	_, err := root.ExecuteC()
	return err
}
