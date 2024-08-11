package main

import (
    "log"
    "time"
    "strings"
    "github.com/mclap/kdbgo"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func runBackendSqlite(backend *Backend) {
    log.Printf("Queries. insert=%q, select=%q", config.Queries.Insert, config.Queries.Select)
    db, err := sql.Open("sqlite3", backend.config.Output)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    for {
        msg, ok := <-backend.ch
        if !ok {
            log.Println(msg, ok, "<-- loop broke!")
            break // exit break loop
        }

        Len := msg.data.Len()
        log.Printf("Received type %d, len %d, tt %d, msg %q", msg.msgtype, Len, msg.data.Type, msg.data)

        switch Len {
        case 3:
            Req := msg.data.Index(0).(*kdb.K).Data.(string)
            Target := msg.data.Index(1).(*kdb.K).Data.(string)
            Params := msg.data.Index(2).(*kdb.K)
            processInsert(db, Req, Target, Params)
        }
    }
}

func processInsert(db *sql.DB, Req string, Target string, Params *kdb.K) {
    log.Printf("Req:%s, Target:%s, Params:%d/%q", Req, Target, Params.Type, Params)

    if Params.Type == kdb.XT {
        t := Params.Data.(kdb.Table)
        log.Printf("Columns: %q, Rows: %q", t.Columns, t.Data)
        var buf strings.Builder
        buf.WriteString("INSERT INTO ")
        buf.WriteString(Target)
        // columns
        buf.WriteString(" (")
        buf.WriteString(strings.Join(t.Columns, ","))
        buf.WriteString(")")
        // placeholders
        buf.WriteString(" VALUES (")
        for pos := range t.Columns {
            if pos > 0 {
                buf.WriteString(",")
            }
            buf.WriteString("?")
        }
        buf.WriteString(") ")

        // Prepare query
        query := buf.String()
        log.Printf("query: %s", query)
        stmt, err := db.Prepare(query)
        if err != nil {
            log.Printf("error: %q", err)
            return
        }
        defer stmt.Close()
        log.Printf("stmt = %q", stmt)

        // Bulk start
        tx, err := db.Begin()
        log.Printf("tx = %q", tx)
        if err != nil {
            log.Printf("begin: %q", err)
            return;
        }
        defer tx.Rollback()

        for row:=0; row < t.Data[0].Len(); row++ {
            args := make([]interface{}, len(t.Columns))
            for ci := range t.Columns {
                elm := t.Data[ci].Index(row)
                switch elm := elm.(type) {
                case kdb.Time: args[ci] = (time.Time)(elm)
                default:
                    args[ci] = elm
                }
            }
            res, err := stmt.Exec(args...)
            log.Printf("exec = %q / %q", res, err)
        }

        stmt.Close()

        err = tx.Commit()
        log.Printf("commit = %q", err)
    }
}
