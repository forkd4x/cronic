package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/a-h/templ"
	"github.com/forkd4x/cronic/models"
	"github.com/forkd4x/cronic/templates"
	"github.com/labstack/echo/v4"
	"github.com/r3labs/sse/v2"
)

//go:embed static/*
var staticFiles embed.FS

type Server struct {
	Echo *echo.Echo
	SSE  *sse.Server
}

func NewServer() Server {
	e := echo.New()

	// Serve embedded static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))))

	e.GET("/", func(c echo.Context) error {
		jobs, err := models.GetJobs()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return Render(c, http.StatusOK, templates.Home(jobs))
	})

	s := sse.New()
	s.AutoReplay = false
	s.CreateStream("updates")

	e.GET("/test", func(c echo.Context) error {
		s.Publish("updates", &sse.Event{
			Data: []byte("this is a test"),
		})
		return nil
	})

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
		Echo: e,
		SSE:  s,
	}
}

// This custom Render replaces Echo's echo.Context.Render() with templ's
// templ.Component.Render().
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
