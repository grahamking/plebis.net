/*
Copyright 2012 Graham King <graham@gkg.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

http://www.gnu.org/licenses/agpl.html
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	PORT  = "8081"
	HTML  = "/usr/local/plebis/index.html"
	STORE = "/usr/local/plebis/store.dat"
)

var store = make([]Message, 0, 25)
var spam = 0

type Context struct {
	Store []Message
	Spam  int
}

func main() {
	load()
	fmt.Println("plebis.net listening on port", PORT)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
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

	context := Context{Store: store, Spam: spam}
	tmpl.Execute(response, context)
}

func doPost(response http.ResponseWriter, request *http.Request) {

	name := request.FormValue("name")
	content := request.FormValue("content")
	date := request.FormValue("date")
	msg := Message{name, content, date}
	if msg.IsSpam() {
		spam++
	} else {
		store = append([]Message{msg}, store[:]...)
	}

	header := response.Header()
	header.Set("Location", "/")
	response.WriteHeader(http.StatusFound)

	go persist()
}

type Message struct {
	Name    string
	Content string
	Date    string
}

// Does this message look like spam
func (self Message) IsSpam() bool {

	if len(self.Date) < 5 {
		return true
	}

	if strings.Count(self.Content, "http://") > 4 {
		return true
	}

	// Date must include the year 20??
	if ! strings.Contains(self.Date, "20") {
		return true
	}

	return false
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
		if os.IsNotExist(openErr) { // No data yet, fine
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
