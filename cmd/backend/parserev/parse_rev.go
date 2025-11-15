package parserev

import (
	"patchy/refs"
	"patchy/util"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "parse-rev <revspec>",
		Short: "Parses a revspec and finds the commit ID it refers to",
		Long:  `Tries to parse a revspec by checking it against known refs, a relative HEAD reference, or a commit id`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hash, err := refs.ParseRev(args[0])
			if err != nil {
				return err
			}
			util.Println(hash)
			return nil
		},
	}
}
