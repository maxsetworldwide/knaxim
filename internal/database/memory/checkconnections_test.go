package memory

import (
	"sync"
	"testing"
)

var testingComplete = &sync.WaitGroup{}

func init() {
	testingComplete.Add(6)
}

func TestConnections(t *testing.T) {
	t.Parallel()
	testingComplete.Wait()
	if CurrentOpenConnections() != 0 {
		t.Fatalf("Connections not being closed: %d", CurrentOpenConnections())
	}
}
