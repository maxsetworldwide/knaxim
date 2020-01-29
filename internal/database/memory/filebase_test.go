package memory

import "testing"

func TestFiles(t *testing.T) {
	t.Parallel()
	fb := DB.File(nil)
	defer fb.Close(nil)

}
