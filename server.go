// Copyright 2014 Hansel Dunlop. Subject to the GPLv3 license.
package main

import (
    "bufio"
    "crypto/sha1"
    "encoding/json"
    "net/http"
    "net/url"
    "html/template"
    "io/ioutil"
    "strings"
    "strconv"
    "fmt"
    "time"
    "os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func splitUrl(url string) (parts []string) {
    trimmed := strings.Trim(url, "/ ")
    parts = strings.Split(trimmed, "/")
    return
}

func docHistory(name string) (history []string) {
    doc, err := ioutil.ReadFile(fmt.Sprintf("accounts/default/%s", name))
    check(err)

    history = strings.Split(string(doc), "\n")
    return
}

func getBlob(hash string) (contents []byte) {
    contents, err := ioutil.ReadFile(fmt.Sprintf("content/%s", hash))
    check(err)
    return
}

func writeBlob(data string) (stringHash string) {
    bytes := []byte(data)
    hash := sha1.Sum(bytes)
    stringHash = fmt.Sprintf("%x", hash)

    f, err := os.Create(fmt.Sprintf("content/%s", stringHash))
    check(err)
    defer f.Close()

    num, err := f.Write(bytes)
    check(err)
    if num < len(bytes) {
        panic("We didn't write all the bytes")
    }
    f.Sync()
    return
}

type History struct {
    Hash string `json:"hash"`
    Timestamp int64 `json:"ts"`
}

type Document struct {
    Name string `json:"name"`
    Encoded string `json:"encoded"`
    History []History `json:"history"`
    Content string `json:"content"`
}

func makeDoc(name string) (doc Document) {
    f, err := os.Open(fmt.Sprintf("accounts/default/%s", name))
    check(err)
    defer f.Close()

    history := []History{}
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        parts := strings.Split(strings.Trim(scanner.Text(), " \n"), " ")
        ts, err := strconv.ParseInt(parts[1], 10, 64)
        check(err)
        hist := History{parts[0], ts}
        history = append(history, hist)

    }
    escaped := url.QueryEscape(name)
    doc = Document{name, escaped, history, ""}
    return
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("frontends/plain-js/index.html")
    check(err)
    t.Execute(w, "")

}

func blobHandler(w http.ResponseWriter, r *http.Request) {
        parts := splitUrl(r.URL.Path)
        contents := getBlob(parts[1])
        fmt.Fprintf(w, "%s", contents)
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
        docs:= []Document{}
        contents, err := ioutil.ReadDir("accounts/default")
        check(err)
        for _, f := range contents {
            doc := makeDoc(f.Name())
            docs = append(docs, doc)
        }
        b, err := json.Marshal(docs)
        check(err)
        fmt.Fprintf(w, "%s", string(b))
}

func docHandler(w http.ResponseWriter, r *http.Request) {
    parts := splitUrl(r.URL.Path)
    name, err := url.QueryUnescape(parts[1])
    check(err)

    doc := makeDoc(name)

    last := doc.History[len(doc.History)-1].Hash

    doc.Content = string(getBlob(last))
    b, err := json.Marshal(doc)
    check(err)
    fmt.Fprintf(w, "%s", string(b))
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
            err := r.ParseForm()
            check(err)

            name := r.Form["name"][0]
            content := r.Form["content"][0]

            hash := writeBlob(content)

            docF, err := os.OpenFile(fmt.Sprintf("accounts/default/%s", name), os.O_APPEND|os.O_WRONLY, 0600)
            check(err)

            updated := time.Now().Unix()
            _, err = docF.WriteString(fmt.Sprintf("%s %d\n", hash, updated))
            check(err)


            docF.Sync()


            fmt.Fprintf(w, "%s", hash)
        } else {
            fmt.Fprintf(w, "<html><title>Error :: Doccer</title></html>")
        }
}

func newHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            err := r.ParseForm()
            check(err)

            name := r.Form["name"][0]

            hash := writeBlob("")

            noteF, err := os.Create(fmt.Sprintf("accounts/default/%s", name))
            check(err)

            encoded := url.QueryEscape(name)

            defer noteF.Close()

            updated := time.Now().Unix()

            _, err = noteF.WriteString(fmt.Sprintf("%s %d\n", hash, updated))
            check(err)


            noteF.Sync()

            http.Redirect(w, r, fmt.Sprintf("/doc/%s", encoded), 301)

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
        http.HandleFunc("/blob/", blobHandler)
        http.HandleFunc("/new", newHandler)
        http.HandleFunc("/save", saveHandler)
        http.HandleFunc("/doc/", docHandler)
        http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), nil)
}
