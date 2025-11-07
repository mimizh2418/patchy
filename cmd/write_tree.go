package cmd

import (
	"patchy/objects"
	"patchy/util"

	"github.com/spf13/cobra"
)

var writeTreeCmd = &cobra.Command{
	Use:   "write-tree <directory>",
	Short: "Recursively writes a tree to the object database and prints its hash",
	Long:  `Recursively writes a tree from a directory to the object database and prints its hash`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash, err := objects.WriteTree(args[0])
		if err != nil {
			return err
		}
		util.Println(hash)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(writeTreeCmd)
}
