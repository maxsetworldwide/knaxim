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

package tag

import (
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	message := "the cheese"
	tags, err := ExtractContentTags(strings.NewReader(message))
	if err != nil {
		t.Fatalf("unable to extract content tags: %s", err.Error())
	}
	if len(tags) != 2 {
		t.Fatalf("incorrect result: %v", tags)
	}
}

func TestName(t *testing.T) {
	name := "the_File.txt"
	tags, err := BuildNameTags(name)
	if err != nil {
		t.Fatalf("unable to build name tags: %s", err.Error())
	}
	if len(tags) != 4 && tags[0].Word == name {
		t.Fatalf("incorrect result: %v", tags)
	}
}
