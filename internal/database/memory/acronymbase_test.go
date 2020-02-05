package memory

import "testing"

func TestAcronym(t *testing.T) {
	t.Parallel()
	defer testingComplete.Done()
	ab := DB.Acronym(nil)
	defer ab.Close(nil)

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
