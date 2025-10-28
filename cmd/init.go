package cmd

import (
	"fmt"
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
			fmt.Println("Error initializing repository:", err)
			return
		}
		fmt.Println("Initialized empty repository in", repoPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
