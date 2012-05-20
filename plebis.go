package main

import (
    "os"
    "fmt"
    "log"
    "io"
    "bufio"
    "net/http"
    "html/template"
    "encoding/json"
)

const (
    PORT = "8080"
    HTML = "/usr/local/plebis/index.html"
    STORE = "/usr/local/plebis/store.dat"
)

var store = make([]Message, 0, 25)

type Context struct {
    Store []Message
}

func main() {
    load()
    fmt.Println("plebis.net listening on port", PORT)
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":" + PORT, nil))
}

func handler(response http.ResponseWriter, request *http.Request) {

    if request.Method == "GET" {
        doGet(response, request)
    } else if request.Method == "POST" {
        doPost(response, request)
    }
}

func doGet(response http.ResponseWriter, request *http.Request) {

    tmpl, err := template.ParseFiles(HTML)
    if err != nil {
        log.Fatal(err)
    }

    context := Context{Store:store}
    tmpl.Execute(response, context)
}

func doPost(response http.ResponseWriter, request *http.Request) {

    name := request.FormValue("name")
    content := request.FormValue("content")
    date := request.FormValue("date")
    msg := Message{name, content, date}
    store = append(store, msg)

    header := response.Header()
    header.Set("Location", "/")
    response.WriteHeader(http.StatusFound)

    go persist()
}

type Message struct {
    Name string
    Content string
    Date string
}

// Save message store to disk
func persist() {

    var jsonData []byte
    var err error

    outFile, openErr := os.Create(STORE)
    if openErr != nil {
        log.Fatal(openErr)
    }

    for _, msg := range store {
        jsonData, err = json.Marshal(msg)
        if err != nil {
            log.Printf(err.Error())
            continue
        }
        outFile.Write(append(jsonData, '\n'))
    }
}

// Load message store from disk
func load() {

    var line []byte
    var err error
    var msg *Message

    outFile, openErr := os.Open(STORE)
    if openErr != nil {
        if os.IsNotExist(openErr) {     // No data yet, fine
            return
        } else {
            log.Fatal(openErr)
        }
    }
    reader := bufio.NewReader(outFile)

    for {
        line, _, err = reader.ReadLine()
        if err != nil {
            if err == io.EOF {
                break
            } else {
                log.Printf(err.Error())
                continue
            }
        }
        msg = new(Message)
        json.Unmarshal(line, &msg)
        store = append(store, *msg)
    }

}
