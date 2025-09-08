package db

import (
    "database/sql"
    "embed"
    "sort"
    "strings"

    _ "github.com/go-sql-driver/mysql"
)

func Connect(url string) (*sql.DB, error) {
    return sql.Open("mysql", url)
}

//go:embed migrations/*.sql
var migrationFS embed.FS

func Migrate(db *sql.DB) error {
    entries, err := migrationFS.ReadDir("migrations")
    if err != nil { return nil } // klasör boşsa sessiz geç
    names := make([]string, 0, len(entries))
    for _, e := range entries { names = append(names, e.Name()) }
    sort.Strings(names)
    for _, name := range names {
        b, err := migrationFS.ReadFile("migrations/" + name)
        if err != nil { return err }
        stmt := strings.ReplaceAll(string(b), "\r", "")
        if _, err := db.Exec(stmt); err != nil { return err } // DSN'de multiStatements=true ise tek seferde birden fazla SQL çalışır
    }
    return nil
}
