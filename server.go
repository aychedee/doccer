package main

import (
    "bufio"
    "bytes"
    "crypto/sha1"
    "encoding/json"
    "net/http"
    "html/template"
    "io/ioutil"
    "strings"
    "strconv"
    "fmt"
    "log"
    "time"
    "os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func forwardSlash(r rune) bool {
    return r == '/'
}

type Document struct {
    Name string `json:"name"`
    Hash string `json:"hash"`
    Created int64 `json:"created"`
}

type Content struct {
    Content string
}


func rootHandler(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("index.html")
    check(err)
    content := Content{"You can start by editing this..."}
    t.Execute(w, content)

}

func docHandler(w http.ResponseWriter, r *http.Request) {
        contentHash := strings.FieldsFunc(r.URL.Path, forwardSlash)
        contents, err := ioutil.ReadFile(fmt.Sprintf("content/%s.md", contentHash[1]))
        check(err)
        t, err := template.ParseFiles("index.html")
        check(err)
        content := Content{string(contents)}
        t.Execute(w, content)
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
        documents := []Document{}
        contents, err := ioutil.ReadDir("accounts/default")
        check(err)
        for _, f := range contents {
            fi, err := os.Open(fmt.Sprintf("accounts/default/%s", f.Name()))
            check(err)
            scanner := bufio.NewScanner(fi)
            for scanner.Scan() {
                fields := strings.Fields(scanner.Text())
                ts, err := strconv.ParseInt(fields[1], 10, 64)
                check(err)

                document := Document{f.Name(), fields[0], ts}
                documents = append(documents, document)
            }
        }
        b, err := json.Marshal(documents)
        check(err)
        fmt.Fprintf(w, string(b))
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

            hash := sha1.Sum(data)

            globF, err := os.Create(fmt.Sprintf("content/%x.md", hash))
            check(err)

            defer globF.Close()

            num, err := globF.WriteString(buffer.String())
            check(err)

            log.Println(num)

            globF.Sync()

            noteF, err := os.Create(fmt.Sprintf("accounts/default/%s", content))
            check(err)

            defer noteF.Close()

            updated := time.Now().Unix()

            written, err := noteF.WriteString(fmt.Sprintf("%x %d", hash, updated))
            check(err)

            log.Println(written)

            noteF.Sync()

            http.Redirect(w, r, fmt.Sprintf("/doc/%x", hash), 301)

        } else {
            fmt.Fprintf(w, "<html><title>Error :: Doccer</title></html>")
        }
}

func main() {
        var address string = "127.0.0.1"
        var port int = 9999

        if len(os.Args) > 2 {
            address = os.Args[1]
            givenPort, err := strconv.Atoi(os.Args[2])
            if err != nil {
                fmt.Println(err)
                os.Exit(2)
            }
            port = givenPort
        }

        http.HandleFunc("/", rootHandler)
        http.HandleFunc("/docs/", docsHandler)
        http.HandleFunc("/doc/", docHandler)
        http.HandleFunc("/doc/new", newDocHandler)
        http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), nil)
}
