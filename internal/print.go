package internal

import (
    "fmt"
    "io"
    "patchy/internal/flags"
)

func Print(a ...any) {
    if !flags.Quiet {
        fmt.Print(a)
    }
}

func Println(a ...any) {
    if !flags.Quiet {
        fmt.Println(a...)
    }
}

func Printf(format string, a ...any) {
    if !flags.Quiet {
        fmt.Printf(format, a...)
    }
}

func Fprint(writer io.Writer, a ...any) (n int, err error) {
    if flags.Quiet {
        return 0, nil
    }
    return fmt.Fprint(writer, a...)
}

func Fprintln(writer io.Writer, a ...any) (n int, err error) {
    if flags.Quiet {
        return 0, nil
    }
    return fmt.Fprintln(writer, a...)
}

func Fprintf(writer io.Writer, format string, a ...any) (n int, err error) {
    return fmt.Fprintf(writer, format, a...)
}
