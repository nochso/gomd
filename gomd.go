package main

import (
	"fmt"
	"github.com/alecthomas/template"
	"github.com/toqueteos/webbrowser"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var (
	port = kingpin.Flag("port", "Listening port used by webserver").Short('p').Default("1110").Int()
	file = kingpin.Arg("file", "Markdown file").String()
)

type EditorView struct {
	File    string
	Content string
}

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", RootHandler)

	go WaitForServer(port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)
	if err != nil {
		log.Fatal("Failed to listen and serve: ", err)
	}
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.URL)
	if req.Method == "POST" {
		req.ParseForm()
		ioutil.WriteFile(*file, []byte(req.PostForm.Get("content")), 0644)
	}
	t := template.New("base.html")
	t, err := t.ParseFiles("template/base.html")
	if err != nil {
		log.Fatalf("Unable to parse template: ", err)
	}
	w.Header().Set("Content-type", "text/html")
	content, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("Unable to open file " + *file)
	}
	ev := EditorView{File: *file, Content: string(content)}
	err = t.Execute(w, ev)
	if err != nil {
		log.Println(err)
	}
}

func WaitForServer(port *int) {
	log.Printf("Waiting for listener on port %d", *port)
	url := fmt.Sprintf("http://localhost:%d", *port)
	for {
		time.Sleep(time.Millisecond * 50)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			continue
		}
		break
	}
	log.Println("Opening " + url)
	if err := webbrowser.Open(url); err != nil {
		log.Fatal(err)
	}
}
