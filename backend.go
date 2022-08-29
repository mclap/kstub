package main

import (
    "log"
)

func runBackend(name string) {
    switch(backends[name].config.Driver) {
    case "csv": runBackendCSV(backends[name])
    case "sqlite": runBackendSqlite(backends[name])
    default:
        log.Fatal("Unknown backend driver:", backends[name].config.Driver)
    }
}
