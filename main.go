package main

import (
    "flag"
    "log"
)

type Backend struct {
    config *BackendConfig
    outputs map[string]string
    context interface{}
    ch chan *KdbMsg
}

var backends = make(map[string]*Backend)

func main() {
    flag.Parse()

    for index, backend := range config.Backend {
        backends[backend.Name] = &Backend{
            config: config.Backend[index],
            outputs: make(map[string]string),
            context: nil,
            ch: make(chan *KdbMsg),
            }
        log.Println("register backend:", backend.Name, backends[backend.Name])
        go runBackend(backend.Name)
    }

    for _, listen := range config.Listen {
        log.Println("starting listener", listen)
        go runAcceptor(listen, backends[listen.Backend])
    }

    // FIXME: sleep forever
    select{}
}
