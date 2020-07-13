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

package query

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"git.maxset.io/web/knaxim/internal/database/memory"
)

func TestMain(m *testing.M) {
	status := m.Run()
	if status == 0 {
		if c := memory.CurrentOpenConnections(); c != 0 {
			status = 2
			fmt.Printf("Not All Connections to database closed: %d connections", c)
		}
	}
	os.Exit(status)
}

type QueryTest struct {
	Query       string
	QueryParams []interface{}
	Expected    []int
}

func TestQuery(t *testing.T) {
	var qtests = []QueryTest{
		QueryTest{
			Query: `{
        "context": "%s",
        "match": "first"
      }`,
			QueryParams: []interface{}{
				owners[0].GetID().String(),
			},
			Expected: []int{0},
		},
		QueryTest{
			Query: `{
        "context": [{
          "type": "owner",
          "id": "%s"
        },{
          "type": "file",
          "id": "%s"
        }],
        "match": {
          "tagtype": "process",
          "word": "test"
        }
      }`,
			QueryParams: []interface{}{
				owners[1].GetID().String(),
				fileinfo[0].ID.String(),
			},
			Expected: []int{0, 2},
		},
	}
	for i, qt := range qtests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var q Q
			if err := json.NewDecoder(strings.NewReader(fmt.Sprintf(qt.Query, qt.QueryParams...))).Decode(&q); err != nil {
				t.Fatalf("Unable to Decode Query String: %s", err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			for _, c := range q.Context {
				if access, err := c.CheckAccess(owners[0], DB); err != nil {
					t.Fatalf("error checking access: %s", err)
				} else if !access {
					t.Fatal("testuser lacks acess")
				}
			}
			if files, err := q.FindMatching(ctx, DB); err != nil {
				t.Fatalf("Error searching: %s", err)
			} else if len(files) != len(qt.Expected) {
				for _, f := range files {
					for i, fd := range fileinfo {
						if f.Equal(fd.ID) {
							t.Logf("matched: %d", i)
						}
					}
				}
				t.Fatalf("incorrect returns: %v", files)
			}
		})
	}
}
