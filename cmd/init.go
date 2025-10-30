package cmd

import (
    "patchy/internal"
    "patchy/internal/util"

    "github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
    Use:   "init [<directory>]",
    Short: "Create an empty repository",
    Long:  `Create an empty repository`,
    Args:  cobra.MaximumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        var repoPath string
        var err error
        if len(args) > 0 {
            repoPath, err = internal.Init(args[0])
        } else {
            repoPath, err = internal.Init(".")
        }
        if err != nil {
            return err
        }
        util.Println("Initialized empty repository in", repoPath)
        return nil
    },
}

func init() {
    rootCmd.AddCommand(initCmd)
}
