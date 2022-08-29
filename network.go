package main

import (
    "fmt"
    "log"
    "net"
    "bufio"
    "io"
    "github.com/mclap/kdbgo"
)

type KdbMsg struct {
    data *kdb.K
    msgtype kdb.ReqType
    e error
}

func runAcceptor(listen *ServerConfig, backend *Backend) {
    log.Println("accepting at", listen.Port, "backend:", backend)
    listener, err := net.Listen("tcp", ":"+fmt.Sprint(listen.Port))
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
        go HandleClientConnectionEx2(conn, handleRequest, backend)
    }
}

func handleRequest(conn net.Conn, data *kdb.K, msgtype kdb.ReqType, e error, backend *Backend) {
    if data != nil {
        //log.Println(data)
        backend.ch <- &KdbMsg{ data:data, msgtype:msgtype, e:e, }
    } else {
        log.Println(e)
    }
}

func HandleClientConnectionEx2(conn net.Conn, handleFunc func(net.Conn, *kdb.K, kdb.ReqType, error, *Backend), backend *Backend ) {
    c := conn.(*net.TCPConn)
    c.SetKeepAlive(true)
    c.SetNoDelay(true)
    var cred = make([]byte, 100)
    n, err := c.Read(cred)
    if err != nil {
        conn.Close()
        return
    }
    c.Write(cred[n-2 : n-1])
    rbuf := bufio.NewReaderSize(conn, 4*1024*1024)
    i := 0
    for {
        d, msgtype, err := kdb.Decode(rbuf)

        if err == io.EOF {
        conn.Close()
        return
        }

        handleFunc(conn, d, msgtype, err, backend)

        if d == nil {
            conn.Close()
            return
        }
    i++
    }
}
