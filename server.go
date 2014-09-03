// Copyright 2014 Hansel Dunlop. Subject to the GPLv3 license.
package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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

func getBlob(hash string) (contents []byte, err error) {
	contents, e := ioutil.ReadFile(fmt.Sprintf("content/%s", hash))
	if e != nil {
		err = errors.New("No such blob")
	}
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
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"ts"`
}

type Document struct {
	Name    string    `json:"name"`
	Encoded string    `json:"encoded"`
	History []History `json:"history"`
	Content string    `json:"content"`
}

func makeDoc(name string) (doc Document) {
	// TODO unit test for the case where  this file doesn't exist
	f, err := os.Open(fmt.Sprintf("accounts/default/%s", name))
	check(err)
	defer f.Close()

	history := []History{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Split(strings.Trim(scanner.Text(), " \n"), " ")
		ts, err := strconv.ParseInt(parts[1], 10, 64)
		check(err)
		date := time.Unix(0, ts)
		hist := History{parts[0], date.UTC()}
		history = append(history, hist)

	}
	escaped := url.QueryEscape(name)
	doc = Document{name, escaped, history, ""}
	return
}

func parseCommit(commit string) (commitMap map[string]string) {
	lines := strings.Split(strings.Trim(string(commit), "\n "), "\n")
	commitMap = make(map[string]string)
	for i := 0; i < len(lines); i++ {
		parts := strings.SplitN(lines[i], " ", 2)
		commitMap[parts[0]] = parts[1]
	}
	return
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontends/plain-js/index.html")
	check(err)
	t.Execute(w, "")

}

func blobHandler(w http.ResponseWriter, r *http.Request) {
	parts := splitUrl(r.URL.Path)
	contents, err := getBlob(parts[1])
	check(err)
	fmt.Fprintf(w, "%s", contents)
}

func docsHandler(w http.ResponseWriter, r *http.Request) {
	docs := []Document{}
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
	var commit string
	// Get commit from url if availble or use the latest
	if len(parts) == 3 {
		commit = parts[2]
	} else {
		commit = doc.History[len(doc.History)-1].Hash
	}
	cContent, err := getBlob(commit)
	check(err)
	c := parseCommit(string(cContent))
	cont, err := getBlob(c["content"])
	check(err)
	doc.Content = string(cont)

	b, err := json.Marshal(doc)
	check(err)
	fmt.Fprintf(w, "%s", string(b))
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		check(err)

		name := r.Form["name"][0]
		content := ""
		if val, ok := r.Form["content"]; ok {
			content = val[0]
		}
		parent := ""
		if val, ok := r.Form["parent"]; ok {
			parent = val[0]
		}
		hash := writeBlob(content)

		cHash := writeBlob(
			fmt.Sprintf(
				"committer %s\nparent %s\ncontent %s\n", "", parent, hash))

		fileName := fmt.Sprintf("accounts/default/%s", name)

		var f *os.File
		if _, err := os.Stat(fileName); err != nil {
			if os.IsNotExist(err) {
				f, err = os.Create(fileName)
				check(err)
			} else {
				check(err)
			}
		} else {
			f, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
			check(err)
		}

		updated := time.Now().UnixNano()
		_, err = f.WriteString(fmt.Sprintf("%s %d\n", cHash, updated))
		check(err)

		f.Sync()

		doc := makeDoc(name)
		last := doc.History[len(doc.History)-1].Hash
		blob, err := getBlob(last)
		check(err)
		doc.Content = string(blob)

		b, err := json.Marshal(doc)
		check(err)
		fmt.Fprintf(w, "%s", string(b))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "<html><title>Error :: Doccer</title></html>")
	}
}

func main() {
	var address string = "127.0.0.1"
	var port int = 4121
    if _, err := os.Stat("content"); err != nil {
        if os.IsNotExist(err) {
            err := os.Mkdir("content", 0744)
            check(err)
        }
    }

	if len(os.Args) > 2 {
		address = os.Args[1]
		givenPort, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		port = givenPort
	}

	staticHandler := http.FileServer(http.Dir("static_path"))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/d/", rootHandler)
	http.HandleFunc("/docs/", docsHandler)
	http.HandleFunc("/blob/", blobHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/doc/", docHandler)
	http.Handle("/static/", staticHandler)

	http.ListenAndServe(fmt.Sprintf("%s:%d", address, port),
		nil)
}
