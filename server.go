package main

import "net/http"
import "fmt"


func rootHandler(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World")
}

func main() {
        http.HandleFunc("/", rootHandler)
        http.ListenAndServe(":9999", nil)
}
