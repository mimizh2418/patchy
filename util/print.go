package util

import (
	"fmt"
	"io"
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

func Fprint(writer io.Writer, a ...any) (n int, err error) {
	if Quiet {
		return 0, nil
	}
	return fmt.Fprint(writer, a...)
}

func Fprintln(writer io.Writer, a ...any) (n int, err error) {
	if Quiet {
		return 0, nil
	}
	return fmt.Fprintln(writer, a...)
}

func Fprintf(writer io.Writer, format string, a ...any) (n int, err error) {
	if Quiet {
		return 0, nil
	}
	return fmt.Fprintf(writer, format, a...)
}
