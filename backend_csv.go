package main

import (
    "log"
)

func runBackendCSV(backend *Backend) {
    for {
        msg, ok := <-backend.ch
        if !ok {
            log.Println(msg, ok, "<-- loop broke!")
            break // exit break loop
        } else {
            log.Println("Received msg:", *msg, *msg.data)
        }
    }
}
