package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/toqueteos/webbrowser"
	"gopkg.in/alecthomas/kingpin.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
	"io"
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

	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}
	e.SetRenderer(t)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static/", "static")

	edit := e.Group("/edit")
	edit.Get("/*", EditHandler)
	edit.Post("/*", EditHandler)

	go WaitForServer(port)
	e.Run(fmt.Sprintf("127.0.0.1:%d", *port))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func EditHandler(c *echo.Context) error {
	filepath := c.P(0)
	if c.Request().Method == "POST" {
		c.Request().ParseForm()
		ioutil.WriteFile(filepath, []byte(c.Request().PostForm.Get("content")), 0644)
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to read requested file")
		log.Fatalf("Unable to open file " + filepath)
	}
	ev := EditorView{File: filepath, Content: string(content)}
	return c.Render(http.StatusOK, "base", ev)
}

func WaitForServer(port *int) {
	log.Printf("Waiting for listener on port %d", *port)
	url := fmt.Sprintf("http://localhost:%d/edit/%s", *port, url.QueryEscape(*file))
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
