package simple_bank

import (
	"os"
	"testing"
)

var testQueries *Queries

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
