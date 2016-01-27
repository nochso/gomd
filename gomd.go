package main

import (
	"flag"
	"github.com/toqueteos/webbrowser"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Configuration struct {
	Port         string
	DocumentRoot string
	ShowHelp     bool
}

var args Configuration

func main() {
	//	wd, err := os.Getwd()
	//	if err != nil {
	//		wd = "."
	//	}
	flag.StringVar(&args.DocumentRoot, "d", ".", "Document root")
	flag.StringVar(&args.Port, "p", "12345", "Listening port used by webserver")
	flag.BoolVar(&args.ShowHelp, "h", false, "Show this help")
	flag.Parse()

	if args.ShowHelp {
		flag.PrintDefaults()
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/asset/", http.StripPrefix("/asset/", http.FileServer(http.Dir("asset"))))
	mux.HandleFunc("/", RootHandler)

	go WaitForServer(args.Port)
	err := http.ListenAndServe(":"+args.Port, mux)
	if err != nil {
		log.Fatal("Failed to listen and serve: ", err)
	}
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/html")
	io.WriteString(w, req.URL.String()+"<br>")
	io.WriteString(w, req.RequestURI+"<br>")
	io.WriteString(w, "hello, world!<br>")
	files, _ := filepath.Glob("*.md")
	for _, file := range files {
		io.WriteString(w, "<a href=\"#\">"+file+"</a><br>")
	}

	//	cwd, _ := os.Getwd()
	//	fis, _ := ioutil.ReadDir(cwd)
	//	for _, fi := range fis {
	//		io.WriteString(w, fi.Name()+"<br>")
	//	}
	io.WriteString(w, "<link rel=\"stylesheet\" href=\"asset/css/font-awesome.min.css\">")
	io.WriteString(w, "<link rel=\"stylesheet\" href=\"asset/css/simplemde.min.css\"><script src=\"asset/js/simplemde.min.js\"></script>")
	io.WriteString(w, "<textarea></textarea>")
	io.WriteString(w, "<script>var simplemde = new SimpleMDE({autoDownloadFontAwesome: false, spellChecker: false });</script>")
}

func WaitForServer(port string) {
	log.Print("Waiting for listener on port " + port)
	url := "http://localhost:" + port
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
