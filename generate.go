package main

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target api --clean api/openapi.yml

//go:generate go run  github.com/go-jet/jet/v2/cmd/jet@latest -source=sqlite -dsn="file://tmp/db.sqlite" -path="infra/sqlite/gen"
