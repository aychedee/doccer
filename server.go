package main

import "bytes"
import "net/http"
import "fmt"
import "os"
import "log"
import "strconv"


func rootHandler(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "<html><title>Doccer</title></html>")
}

func newDocHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            err := r.ParseForm()
            if err != nil {
                return
            }

            log.Println(r.Form["name"][0])
            content := r.Form["name"][0]
            var buffer bytes.Buffer

            buffer.WriteString(content + "\n")
            for i :=0; i < len(content); i++ {
                buffer.WriteString("=")
            }
            log.Println(buffer.String())
            log.Println(content)
            http.Redirect(w, r, "/doc/8d6ac39986ccd929d7cc1825efb0faa841a46e0a", 301)

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
