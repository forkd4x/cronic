package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

//go:embed static/*
var staticFiles embed.FS

func NewServer(cronic *Cronic) *echo.Echo {
	e := echo.New()

	// Serve embedded static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(err)
	}
	e.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))))

	e.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, Home(cronic))
	})

	return e
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
