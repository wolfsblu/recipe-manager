package main

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target api --clean api/openapi.yml

// Note: go-jet table definitions are manually maintained in infra/sqlite/gen/
// If the schema changes significantly, update infra/sqlite/gen/table/table.go and infra/sqlite/gen/model/models.go
// Or regenerate using: jet -source=sqlite -dsn="file://$(pwd)/tmp/db.sqlite" -path="./infra/sqlite/gen"
// (requires CGO_ENABLED=1 and a C compiler)
