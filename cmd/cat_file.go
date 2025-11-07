package cmd

import (
	"errors"
	objects2 "patchy/objects"
	"patchy/objects/objecttype"

	"github.com/spf13/cobra"
)

var catFileCmd = &cobra.Command{
	Use:   "cat-file <object-hash>",
	Short: "Provides details about an object",
	Long:  `Outputs the contents or details of an object given its hash`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		objType, err := objects2.ReadObjectType(args[0])
		if err != nil {
			return err
		}
		switch objType {
		case objecttype.Blob:
			return objects2.PrintBlob(args[0])
		case objecttype.Tree:
			return objects2.PrintTree(args[0])
		case objecttype.Commit:
			return objects2.PrintCommit(args[0])
		default:
			return errors.New("unknown object type")
		}
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)
}
