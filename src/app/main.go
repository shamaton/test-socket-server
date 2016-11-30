package main

import (
	"app/api"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/labstack/echo"
)

const BIND = ":8080"

// templは1つのテンプレートを表します
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTPはHTTPリクエストを処理します
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
	return t.templ.Execute(c.Response().Writer(), c.Request())
}

func main() {
	e := echo.New()

	/*
		e.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})
	*/
	t := &templateHandler{filename: "chat.html"}
	e.GET("/", t.Render)

	e.GET("/get_and_create", api.GetSocketAndCreateRoom)
	e.GET("/get", api.GetSocket)
	e.Logger.Fatal(e.Start(BIND))
}

func top(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}
