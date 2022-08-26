package main

import (
    "fmt"
    "net"
    "log"
    "github.com/mclap/kdbgo"
)

func main() {
    runAcceptor(5001)
}

func runAcceptor(port int) {
    log.Println("accepting at", port)
    listener, err := net.Listen("tcp", ":"+fmt.Sprint(port))
    if err != nil {
        // Use fatal to exit if the listener fails to start
        log.Fatal(err)
    }
    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            // Print the error using a log.Fatal would exit the server
            log.Println(err)
        }
        // Using a go routine to handle the connection
        go kdb.HandleClientConnectionEx(conn, handleRequest)
    }
}

func handleRequest(conn net.Conn, data *kdb.K, msgtype kdb.ReqType) {
    log.Println(data)
}
