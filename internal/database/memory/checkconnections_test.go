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

import (
	"sync"
	"testing"
)

var testingComplete = &sync.WaitGroup{}

func init() {
	testingComplete.Add(7)
}

func TestConnections(t *testing.T) {
	t.Parallel()
	testingComplete.Wait()
	t.Log("Checking Connections")
	if CurrentOpenConnections() != 0 {
		t.Fatalf("Connections not being closed: %d", CurrentOpenConnections())
	}
}
