/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

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
