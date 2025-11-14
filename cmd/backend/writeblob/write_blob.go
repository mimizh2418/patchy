package writeblob

import (
	"patchy/objects"
	"patchy/util"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "write-blob <file>",
		Short: "Write a blob object from a file and print its hash",
		Long:  `Create and write a blob object to the object database from a file and print its hash`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hash, err := objects.WriteBlob(args[0])
			if err != nil {
				return err
			}
			util.Println(hash)
			return nil
		},
	}
	return command
}
