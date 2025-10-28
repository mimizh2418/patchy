package cmd

import (
	"patchy/internal"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an empty repository",
	Long:  `Create an empty repository`,
	Run: func(cmd *cobra.Command, args []string) {
		var repoPath string
		var err error
		if len(args) > 0 {
			repoPath, err = internal.Init(args[0])
		} else {
			repoPath, err = internal.Init(".")
		}

		if err != nil {
			internal.Println("Error initializing repository:", err)
			return
		}
		internal.Println("Initialized empty repository in", repoPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
