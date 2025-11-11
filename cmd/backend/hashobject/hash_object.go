package hashobject

import (
    "patchy/objects"
    "patchy/util"

    "github.com/spf13/cobra"
)

var write bool

func NewCommand() *cobra.Command {
    command := &cobra.Command{
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
    command.Flags().BoolVarP(&write, "write", "w", false, "write object into the object database")
    return command
}
