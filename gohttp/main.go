package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {

    //Initialize the server to process requests
    s := Server{  }

    // start server
    fmt.Printf("Starting server on :8080\n")
    fmt.Println("--------------------------------------------------")

    handler := http.NewServeMux()
    handler.HandleFunc("/gohttp", s.gohttp)

    http := &http.Server{
        Addr:      ":8080",
        Handler:   handler,
    }

    log.Fatal(http.ListenAndServe())

}
