package log

import (
	"patchy/objects"
	"patchy/refs"
	"patchy/util"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var oneLine bool

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "log",
		Short: "Shows commit logs",
		Long:  `Shows commit logs`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var startCommitHash string
			if len(args) < 1 {
				headState, err := refs.ReadHead()
				if err != nil {
					return err
				}
				startCommitHash = headState.Commit
			} else {
				parsedHash, err := refs.ParseRev(args[0])
				if err != nil {
					return err
				}
				startCommitHash = parsedHash
			}

			startCommit, err := objects.ReadCommit(startCommitHash)
			if err != nil {
				return err
			}

			// TODO: Handle multiple parents/children
			// TODO: Markers for branches, tags, HEAD, etc.
			currentCommit := startCommit
			currentCommitHash := startCommitHash
			for currentCommit != nil {
				if oneLine {
					util.Printf("* ")
					util.ColorPrint(color.FgYellow, currentCommitHash[:7])
					util.Printf(" ")
					util.Println(strings.Split(currentCommit.Message, "\n")[0])
				} else {
					util.ColorPrintf(color.FgYellow, "commit %s\n", currentCommitHash)
					util.Println("Author: ", currentCommit.Author)
					util.Println("Date:   ", currentCommit.Time)
					util.Println()
					util.Println("    ", strings.ReplaceAll(currentCommit.Message, "\n", "\n    "))
					util.Println()
				}
				if currentCommit.Parent == nil {
					break
				}
				currentCommitHash = *currentCommit.Parent
				currentCommit, err = objects.ReadCommit(currentCommitHash)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	command.Flags().BoolVar(&oneLine, "oneline", false, "Display each commit on a single line")
	return command
}
