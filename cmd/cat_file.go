package cmd

import (
	"errors"
	"patchy/internal/objects"
	"patchy/internal/util"

	"github.com/spf13/cobra"
)

var catFileCmd = &cobra.Command{
	Use:   "cat-file <object-hash>",
	Short: "Provides details about an object",
	Long:  `Outputs the contents or details of an object given its hash`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		objType, err := objects.ReadObjectType(args[0])
		if err != nil {
			return err
		}
		switch objType {
		case objects.Blob:
			return objects.PrintBlob(args[0])
		case objects.Tree:
			return objects.PrintTree(args[0])
		case objects.Commit:
			util.Println() // TODO write correct format for this
		default:
			return errors.New("unknown object type")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)
}
