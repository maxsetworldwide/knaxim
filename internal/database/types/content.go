// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
