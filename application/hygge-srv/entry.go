package hygge_srv

import (
    "github.com/spf13/cobra"
)

func Execute() {}

var cmdRoot = &cobra.Command{
    Use:   "hygge-srv",
    Short: "A brief description of your application",
    RunE:  execRoot,
}

func execRoot(cmd *cobra.Command, args []string) error {
    cmd.SilenceUsage = true
}
