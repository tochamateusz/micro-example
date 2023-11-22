package testservice

import (
	"fmt"

	"gorm.io/gorm"
)

func ProvideTestService(db *gorm.DB) {
	fmt.Printf("DB: %+v\n", db)

	var result []map[string]interface{}

	db.Raw("SELECT * FROM pg_indexes where tablename='account' and schemaname='public'").Scan(&result)
	for i, rows := range result {

		fmt.Printf("\n\n--------------Row [%d]--------------\n\n", i)
		for fields, v := range rows {
			fmt.Printf("[%+v]:\t%+v\n", fields, v)

		}
	}

}
