package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/toqueteos/webbrowser"
	"github.com/GeertJohan/go.rice"
	"gopkg.in/alecthomas/kingpin.v2"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type InputArgs struct {
	Port *int
	File *string
}

var args = InputArgs{
	Port: kingpin.Flag("port", "Listening port used by webserver").Short('p').Default("1110").Int(),
	File: kingpin.Arg("file", "Markdown file").String(),
}

type EditorView struct {
	File    string
	Content string
}

func main() {
	// Parse command line arguments
	kingpin.Version("0.0.1")
	kingpin.Parse()

	// Prepare (optionally) embedded resources
	templateBox := rice.MustFindBox("template")
	staticHttpBox := rice.MustFindBox("static").HTTPBox()
	staticServer := http.StripPrefix("/static/", http.FileServer(staticHttpBox))

	e := echo.New()

	t := &Template{
		templates: template.Must(template.New("base").Parse(templateBox.MustString("base.html"))),
	}
	e.SetRenderer(t)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Get("/static/*", func(c *echo.Context) error {
		staticServer.ServeHTTP(c.Response().Writer(), c.Request())
		return nil
	})

	edit := e.Group("/edit")
	edit.Get("/*", EditHandler)
	edit.Post("/*", EditHandler)

	go WaitForServer()
	e.Run(fmt.Sprintf("127.0.0.1:%d", *args.Port))
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

func WaitForServer() {
	log.Printf("Waiting for listener on port %d", *args.Port)
	url := fmt.Sprintf("http://localhost:%d/edit/%s", *args.Port, url.QueryEscape(*args.File))
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
