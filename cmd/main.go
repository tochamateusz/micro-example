package main

import (
	"fmt"

	"github.com/tochamateusz/micro/modules/database"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func main() {
	fx.New(
		database.Module(),
		fx.Decorate(func(*database.Config) *database.Config {
			return &database.Config{
				Host:     "0.0.0.0",
				User:     "root",
				Password: "password",
				DbName:   "simple_bank",
				Port:     5432,
			}
		}),
		fx.Invoke(
			func(db *gorm.DB) {
				fmt.Printf("DB: %+v\n", db)

				var result []map[string]interface{}

				db.Raw("select * from information_schema.columns where table_name='account' and table_schema='public'").Scan(&result)
				for i, rows := range result {

					fmt.Printf("\n\n--------------Row [%d]--------------\n\n", i)
					for fields, v := range rows {
						fmt.Printf("[%+v]:\t%+v\n", fields, v)

					}
				}

			}),
	).Run()
}
