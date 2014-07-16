package main

import (
    "bytes"
    "crypto/sha1"
    "net/http"
    "fmt"
    "log"
    "os"
    "strconv"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "<html><title>Doccer</title></html>")
}


func newDocHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            err := r.ParseForm()
            check(err)

            content := r.Form["name"][0]
            var buffer bytes.Buffer

            buffer.WriteString(content + "\n")
            for i :=0; i < len(content); i++ {
                buffer.WriteString("=")
            }
            buffer.WriteString("\n")
            data := []byte(buffer.String())

            f, err := os.Create(fmt.Sprintf("content/%x.md", sha1.Sum(data)))
            check(err)

            defer f.Close()

            num, err := f.WriteString(buffer.String())
            check(err)

            log.Println(num)

            f.Sync()

            http.Redirect(w, r, fmt.Sprintf("/doc/%x", sha1.Sum(data)), 301)

        } else {
            fmt.Fprintf(w, "<html><title>Error :: Doccer</title></html>")
        }
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
        http.HandleFunc("/doc/new", newDocHandler)
        http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), nil)
}
