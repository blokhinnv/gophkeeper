package cotesting

import "github.com/spf13/cobra"

func ExecuteCommandC(
	root *cobra.Command,
	args ...string,
) error {
	root.SetArgs(args)
	_, err := root.ExecuteC()
	return err
}
