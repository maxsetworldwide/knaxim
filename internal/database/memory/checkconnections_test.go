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
