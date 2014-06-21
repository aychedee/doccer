package main

import "net/http"
import "fmt"
import "os"
import "strconv"


func rootHandler(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "<html><title>Doccer</title></html>")
}

func main() {
        var address string = "127.0.0.1"
        var port int = 9999

        if len(os.Args) > 2 {
            address = os.Args[1]
            given_port, err := strconv.Atoi(os.Args[2])
            if err != nil {
                fmt.Println(err)
                os.Exit(2)
            }
            port = given_port
        }

        http.HandleFunc("/", rootHandler)
        http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), nil)
}
