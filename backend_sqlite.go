package main

import (
    "log"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func runBackendSqlite(backend *Backend) {
    db, err := sql.Open("sqlite3", backend.config.Output)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    for {
        msg, ok := <-backend.ch
        if ok == false {
            log.Println(msg, ok, "<-- loop broke!")
            break // exit break loop
        } else {
            log.Println("Received msg:", *msg, *msg.data)
        }
    }
}
