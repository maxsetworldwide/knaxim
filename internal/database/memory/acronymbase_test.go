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

package memory

import "testing"

func TestAcronym(t *testing.T) {
	defer testingComplete.Done()
	DB.Connect(nil)
	ab := DB.Acronym()
	defer DB.Close(nil)
	t.Parallel()

	t.Log("Acronym Put")
	err := ab.Put("t", "test")
	if err != nil {
		t.Fatalf("Unable to put acronym: %s", err)
	}

	t.Log("Acronym Get")
	matches, err := ab.Get("t")
	if err != nil {
		t.Fatalf("Unable to get acronym: %s", err)
	}
	if len(matches) != 1 || matches[0] != "test" {
		t.Fatalf("incorrect matches: %v", matches)
	}
}
