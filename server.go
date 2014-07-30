package main

import (
    "bufio"
    "bytes"
    "crypto/sha1"
    "encoding/json"
    "net/http"
    "net/url"
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
    Encoded string `json:"encoded"`
    Hash string `json:"hash"`
    Created int64 `json:"created"`
}

type Content struct {
    Content string
    Name string
}


func rootHandler(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("index.html")
    check(err)
    content := Content{"You can start by editing this...", "None"}
    t.Execute(w, content)

}

func getBlob(hash string) (contents []byte) {
    contents, err := ioutil.ReadFile(fmt.Sprintf("content/%s", hash))
    check(err)
    return
}

func blobHandler(w http.ResponseWriter, r *http.Request) {
        contentHash := strings.FieldsFunc(r.URL.Path, forwardSlash)
        contents := getBlob(contentHash[1])
        fmt.Fprintf(w, "%s", contents)
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
        documents := []Document{}
        contents, err := ioutil.ReadDir("accounts/default")
        check(err)
        for _, f := range contents {
            fi, err := os.Open(fmt.Sprintf("accounts/default/%s", f.Name()))
            check(err)
            defer fi.Close()
            scanner := bufio.NewScanner(fi)
            for scanner.Scan() {
                fields := strings.Fields(scanner.Text())
                ts, err := strconv.ParseInt(fields[1], 10, 64)
                check(err)
                escaped := url.QueryEscape(f.Name())

                document := Document{f.Name(), escaped, fields[0], ts}
                documents = append(documents, document)
            }
        }
        b, err := json.Marshal(documents)
        check(err)
        fmt.Fprintf(w, "%s", string(b))
}

func docHandler(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("index.html")
    check(err)

    urlParts := strings.FieldsFunc(r.URL.Path, forwardSlash)
    name, err := url.QueryUnescape(urlParts[1])

    doc, err := ioutil.ReadFile(fmt.Sprintf("accounts/default/%s", name))
    check(err)

    edits := strings.Split(string(doc), "\n")

    last := edits[len(edits)-2]
    fields := strings.Fields(last)

    latest := getBlob(fields[0])
    content := Content{string(latest), name}
    t.Execute(w, content)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
            err := r.ParseForm()
            check(err)

            name := r.Form["name"][0]
            content := r.Form["content"][0]
            bytes := []byte(content)
            hash := sha1.Sum(bytes)
            f, err := os.Create(fmt.Sprintf("content/%x", hash))
            check(err)

            defer f.Close()

            num, err := f.WriteString(content)
            check(err)

            log.Println(num)
            f.Sync()

            docF, err := os.OpenFile(fmt.Sprintf("accounts/default/%s", name), os.O_APPEND|os.O_WRONLY, 0600)
            check(err)

            updated := time.Now().Unix()
            written, err := docF.WriteString(fmt.Sprintf("%x %d\n", hash, updated))
            check(err)

            log.Println(written)

            docF.Sync()


            fmt.Fprintf(w, "%x", hash)
        } else {
            fmt.Fprintf(w, "<html><title>Error :: Doccer</title></html>")
        }
}

func newHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            err := r.ParseForm()
            check(err)

            name := r.Form["name"][0]
            var buffer bytes.Buffer

            buffer.WriteString(name + "\n")
            for i :=0; i < len(name); i++ {
                buffer.WriteString("=")
            }
            buffer.WriteString("\n")
            data := []byte(buffer.String())

            hash := sha1.Sum(data)

            globF, err := os.Create(fmt.Sprintf("content/%x", hash))
            check(err)

            defer globF.Close()

            num, err := globF.WriteString(buffer.String())
            check(err)

            log.Println(num)

            globF.Sync()

            noteF, err := os.Create(fmt.Sprintf("accounts/default/%s", name))
            check(err)

            encoded := url.QueryEscape(name)

            defer noteF.Close()

            updated := time.Now().Unix()

            written, err := noteF.WriteString(fmt.Sprintf("%x %d\n", hash, updated))
            check(err)

            log.Println(written)

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
