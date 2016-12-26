package main

import (
	"front/api"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"log"

	"front/socket"

	"github.com/labstack/echo"
)

const BIND = ":8080"

const (
	origin = "http://localhost:8080"
	url    = "ws://localhost:8081"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// serve template
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		pwd, _ := os.Getwd()

		t.templ = template.Must(template.ParseFiles(filepath.Join(pwd, "src/app/room/templates",
			t.filename)))
	})
	t.templ.Execute(w, r)
}

func (t *templateHandler) Render(c echo.Context) error {
	t.once.Do(func() {
		pwd, _ := os.Getwd()

		t.templ = template.Must(template.ParseFiles(filepath.Join(pwd, "src/app/room/templates",
			t.filename)))
	})
	return t.templ.Execute(c.Response().Writer, c.Request())
}

func main() {

	// connect backend server
	err := socket.ConnectBack(url)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	t := &templateHandler{filename: "chat.html"}
	e.GET("/", t.Render)

	e.GET("/get", api.GetSocket)
	e.GET("/ping", api.Ping)
	e.Logger.Fatal(e.Start(BIND))
}

func top(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}
