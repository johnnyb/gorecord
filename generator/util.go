package generator

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/johnnyb/gorecord/inflect"
	"io"
	"os"
)

type FileFunc func(io.Writer)

func WithFileAppend(fname string, appender FileFunc) error {
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	if f == nil {
		return errors.New("File did not open")
	}
	defer f.Close()

	fh := bufio.NewWriter(f)
	defer fh.Flush()

	appender(fh)
	return nil
}

// WriteFile is a convenience function for writing whole files at once
func WriteFile(fname string, data string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	fh := bufio.NewWriter(f)
	defer fh.Flush()

	fmt.Fprint(fh, data)

	return nil
}

// ColumnToStructMember takes a column name and converts it to the name of the struct value.
func ColumnToStructMember(val string) string {
	return "Raw" + inflect.Camelize(val)
}
