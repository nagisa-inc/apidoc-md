package apidocmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type buffer struct {
	buf *bytes.Buffer
}

func NewBuffer(b []byte) *buffer {
	return &buffer{
		buf: bytes.NewBuffer(b),
	}
}

func (b buffer) String() string {
	return b.buf.String()
}

func (b *buffer) Writeln(format string, a ...interface{}) (int, error) {
	return b.buf.WriteString(fmt.Sprintf(format, a...) + "  \n\n")
}

func (b buffer) writeToFile(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.WriteString(file, b.buf.String()); err != nil {
		return err
	}

	return nil
}
