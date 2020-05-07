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
