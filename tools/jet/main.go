package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/sqlite"
	"github.com/go-jet/jet/v2/generator/template"
	sqlite2 "github.com/go-jet/jet/v2/sqlite"
	"github.com/wolfsblu/recipe-manager/infra/env"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", env.MustGet("DB_PATH"))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	destDir := "infra/sqlite/gen"
	err = os.RemoveAll(destDir)
	if err != nil {
		log.Fatalf("failed to remove destination directory: %v", err)
	}

	err = sqlite.GenerateDB(db, destDir, template.Default(sqlite2.Dialect).
		UseSchema(func(schema metadata.Schema) template.Schema {
			return template.DefaultSchema(schema).
				UseModel(template.DefaultModel().
					UseTable(func(table metadata.Table) template.TableModel {
						return template.DefaultTableModel(table).
							UseField(func(column metadata.Column) template.TableModelField {
								defaultTableModelField := template.DefaultTableModelField(column)
								if column.DataType.Name == "INTEGER" {
									defaultTableModelField.Type = template.NewType(int64(0))
								} else if column.DataType.Name == "REAL" {
									defaultTableModelField.Type = template.NewType(float64(0))
								}
								return defaultTableModelField
							})
					}))
		}))

	if err != nil {
		log.Fatalf("failed to generate jet files: %v", err)
	}

	log.Println("successfully generated jet files")
}
