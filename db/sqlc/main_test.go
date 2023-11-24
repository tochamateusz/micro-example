package simple_bank

import (
	"fmt"
	"log"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@0.0.0.0:5432/simple_bank?sslmode=disable"
)

func TestMain(m *testing.M) {

	dns := fmt.Sprintf("host=0.0.0.0 user=root password=password dbname=simple_bank port=5432 sslmode=disable")

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
  
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}

	testQueries = New(db.ConnPool)

	os.Exit(m.Run())
}
