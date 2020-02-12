package memory

import (
	"flag"
	"os"
	"testing"
)

var DB Database

func TestMain(m *testing.M) {
	flag.Parse()
	DB.Init(nil, true)
	fill(&DB)
	os.Exit(m.Run())
}

func fill(db *Database) {
	fillowners(db)
	fillstores(db)
}
