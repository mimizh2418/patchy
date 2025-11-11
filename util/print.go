package util

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

var Quiet bool

func Print(a ...any) {
	if !Quiet {
		fmt.Print(a)
	}
}

func Println(a ...any) {
	if !Quiet {
		fmt.Println(a...)
	}
}

func Printf(format string, a ...any) {
	if !Quiet {
		fmt.Printf(format, a...)
	}
}

func ColorPrint(attribute color.Attribute, a ...any) {
	if !Quiet {
		_, _ = color.New(attribute).Print(a...)
	}
}

func ColorPrintln(attribute color.Attribute, a ...any) {
	if !Quiet {
		_, _ = color.New(attribute).Println(a...)
	}
}

func ColorPrintf(attribute color.Attribute, format string, a ...any) {
	if !Quiet {
		_, _ = color.New(attribute).Printf(format, a...)
	}
}

func Fprint(writer io.Writer, a ...any) {
	if !Quiet {
		_, _ = fmt.Fprint(writer, a...)
	}
}

func Fprintln(writer io.Writer, a ...any) {
	if !Quiet {
		_, _ = fmt.Fprintln(writer, a...)
	}
}

func Fprintf(writer io.Writer, format string, a ...any) {
	if !Quiet {
		_, _ = fmt.Fprintf(writer, format, a...)
	}
}
