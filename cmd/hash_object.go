package cmd

import (
	"patchy/objects"
	"patchy/util"

	"github.com/spf13/cobra"
)

var hashObjectCmd = &cobra.Command{
	Use:   "hash-object <file> [--write]",
	Short: "Compute an object hash and optionally create and write an object",
	Long:  `Compute an object hash from a file and optionally compress and write the file into the object database`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hash, blob, err := objects.HashObject(args[0])
		if err != nil {
			return err
		}
		util.Println(hash)
		if !write {
			return nil
		}
		return objects.WriteObject(hash, blob)
	},
}

var write bool

func init() {
	rootCmd.AddCommand(hashObjectCmd)
	hashObjectCmd.Flags().BoolVarP(&write, "write", "w", false, "write object into the object database")
}
