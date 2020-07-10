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
	"io"
	"strings"
	"testing"
)

func TestContent(t *testing.T) {
	t.Parallel()
	sid := StoreID{
		Hash:  10,
		Stamp: 10,
	}
	lines := []ContentLine{
		ContentLine{
			ID:       sid,
			Position: 0,
			Content:  []string{"a"},
		},
		ContentLine{
			ID:       sid,
			Position: 2,
			Content:  []string{"c"},
		},
		ContentLine{
			ID:       sid,
			Position: 1,
			Content:  []string{"b"},
		},
	}

	rdr, err := NewContentReader(lines)
	if err != nil {
		t.Fatalf("Failed to create Content Reader: %s", err)
	}

	sb := new(strings.Builder)
	if _, err := io.Copy(sb, rdr); err != nil {
		t.Fatalf("Unable to Read: %s", err)
	}

	if s := sb.String(); s != "abc" {
		t.Fatalf("Incorrect resulting string: %s", s)
	}
}
