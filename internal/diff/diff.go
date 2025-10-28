package diff

import (
	"os"
	"patchy/internal"
	"patchy/internal/parser"
	"text/tabwriter"

	"github.com/fatih/color"
)

type Operation struct {
	Type    OpType
	OldLine int
	NewLine int
	OldText string
	NewText string
}

type OpType int

const (
	Equal OpType = iota
	Insert
	Delete
)

func delOp(line int, text string) Operation {
	return Operation{
		Type:    Delete,
		OldLine: line,
		NewLine: -1,
		OldText: text,
		NewText: "",
	}
}

func insOp(line int, text string) Operation {
	return Operation{
		Type:    Insert,
		OldLine: -1,
		NewLine: line,
		OldText: "",
		NewText: text,
	}
}

func eqlOp(oldLine, newLine int, text string) Operation {
	return Operation{
		Type:    Equal,
		OldLine: oldLine,
		NewLine: newLine,
		OldText: text,
		NewText: text,
	}
}

func Diff(file1, file2 string) ([]Operation, error) {
	a, err := parser.ReadFile(file1)
	if err != nil {
		return nil, err
	}
	b, err := parser.ReadFile(file2)
	if err != nil {
		return nil, err
	}

	return myers(a, b), nil
}

func PrintDiff(ops []Operation) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	defer func() {
		_ = writer.Flush()
	}()

	white := color.New(color.FgWhite)
	green := color.New(color.FgGreen)
	red := color.New(color.FgRed)

	for _, op := range ops {
		switch op.Type {
		case Equal:
			_, _ = internal.Fprintln(writer, white.Sprintf("\t%d\t \t%s", op.OldLine, op.NewText))
		case Insert:
			_, _ = internal.Fprintln(writer, green.Sprintf("+\t \t%d\t%s", op.NewLine, op.NewText))
		case Delete:
			_, _ = internal.Fprintln(writer, red.Sprintf("-\t%d\t \t%s", op.OldLine, op.OldText))
		}
	}

}
