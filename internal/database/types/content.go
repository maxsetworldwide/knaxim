package types

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// ContentLine is a line of content within a file
type ContentLine struct {
	ID StoreID `bson:"id"`
	//PageNum  int              `bson:"pagenum"`
	Position int      `bson:"position"`
	Content  []string `bson:"content"`
}

// NewContentReader streams content lines together into a reader
func NewContentReader(lines []ContentLine) (result io.Reader, err error) {
	defer func() {
		if r := recover(); r != nil {
			result = nil
			switch v := r.(type) {
			case error:
				err = v
			case string:
				err = errors.New(v)
			default:
				err = fmt.Errorf("Building Content Reader %v", v)
			}
		}
	}()
	out := make([]ContentLine, len(lines))
	copy(out, lines)
	for i := range out {
		for target := out[i].Position; target != i; target = out[i].Position {
			if target == out[target].Position {
				panic("double position")
			}
			out[i], out[target] = out[target], out[i]
		}
	}
	linereaders := make([]io.Reader, 0, len(out))
	for _, line := range out {
		linereaders = append(linereaders, strings.NewReader(strings.Join(line.Content, ", ")))
	}
	return io.MultiReader(linereaders...), nil
}