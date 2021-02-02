package genutils

import (
	"bytes"
	"fmt"
)

func newContentBlock(indent int) Block {
	return &contentBlock{indent: indent}
}

// contentBlock is a content block used for writing full blocks.
type contentBlock struct {
	buf           bytes.Buffer
	indent        int
	beforeContent []func() error
}

// P is a function similar to the Writer.P but it writes the content with specific indent.
func (c *contentBlock) P(args ...interface{}) {
	for i := 0; i < c.indent; i++ {
		fmt.Fprint(&c.buf, "\t")
	}
	for _, x := range args {
		fmt.Fprint(&c.buf, x)
	}
	fmt.Fprintln(&c.buf)
}

// Content gets the content of given block.
func (c *contentBlock) Content() ([]byte, error) {
	var err error
	for _, hook := range c.beforeContent {
		if err = hook(); err != nil {
			return nil, err
		}
	}
	return c.buf.Bytes(), nil
}
