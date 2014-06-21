package main

import "net/http"
import "fmt"
import "os"


func rootHandler(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "<html><title>Doccer</title></html>")
}

func main() {
        address := os.Args[1]
        port := os.Args[2]

        http.HandleFunc("/", rootHandler)
        http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), nil)
}
