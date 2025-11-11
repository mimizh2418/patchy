package catfile

import (
	"errors"
	"patchy/objects"
	"patchy/objects/objecttype"
	"patchy/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
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
			case objecttype.Blob:
				util.ColorPrintf(color.FgCyan, "[blob %s]\n", args[0])
				return objects.PrintBlob(args[0])
			case objecttype.Tree:
				util.ColorPrintf(color.FgCyan, "[tree %s]\n", args[0])
				return objects.PrintTree(args[0])
			case objecttype.Commit:
				util.ColorPrintf(color.FgCyan, "[commit %s]\n", args[0])
				return objects.PrintCommit(args[0])
			default:
				return errors.New("unknown object type")
			}
		},
	}
}
