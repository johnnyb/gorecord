package generator

import (
	"bufio"
	"fmt"
	"os"
)

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
