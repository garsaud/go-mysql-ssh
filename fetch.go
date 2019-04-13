package db

import (
    "database/sql"
    "github.com/go-sql-driver/mysql"
    "net"
)

func fetch(uri string, query string, callback func(row *sql.Rows)) error {
    db, err := sql.Open("mysql", uri)
    if err != nil { return err }
    defer db.Close()

    rows, err := db.Query(query)
    if err != nil { return err }

    for rows.Next() {
        callback(rows)
    }

    return nil
}
