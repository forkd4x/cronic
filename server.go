package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/forkd4x/cronic/models"
	"github.com/labstack/echo/v4"
	"github.com/r3labs/sse/v2"
)

//go:embed static/*
var staticFiles embed.FS

type Server struct {
	Echo      *echo.Echo
	SSE       *sse.Server
	Templates *template.Template
}

func (s Server) RenderTemplate(name string, data any) ([]byte, error) {
	var b bytes.Buffer
	err := s.Templates.ExecuteTemplate(&b, name, data)
	return b.Bytes(), err
}

func NewServer() Server {
	e := echo.New()

	// Serve embedded static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	staticHandler := http.FileServer(http.FS(staticFS))
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", staticHandler)))

	e.Renderer = TemplateRenderer()

	e.GET("/", func(c echo.Context) error {
		jobs, err := models.GetJobs()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.Render(http.StatusOK, "home.html", map[string]any{"Jobs": jobs})
	})

	s := sse.New()
	s.AutoReplay = false
	s.CreateStream("updates")
	e.GET("/sse", func(c echo.Context) error {
		fmt.Println("SSE client connected")
		go func() {
			<-c.Request().Context().Done()
			fmt.Println("SSE client disconnected")
		}()
		s.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	return Server{
		Echo:      e,
		SSE:       s,
		Templates: e.Renderer.(*Template).Templates,
	}
}
