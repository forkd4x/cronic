package main

import (
	"embed"
	"html/template"
	"io"
	"time"

	"github.com/forkd4x/cronic/models"
	"github.com/labstack/echo/v4"
)

//go:embed templates/*.html
var templateFiles embed.FS

type Template struct {
	Templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func TemplateRenderer() *Template {
	return &Template{
		Templates: template.Must(
			template.
				New("").
				Funcs(template.FuncMap{
					"getStatusClass": getStatusClass,
					"formatTime":     formatTime,
					"formatDuration": formatDuration,
					"hasRunning":     hasRunning,
				}).
				ParseFS(templateFiles, "templates/*.html"),
		),
	}
}

func getStatusClass(status string) string {
	switch status {
	case "Pending":
		return "bg-light text-secondary"
	case "Running":
		return "text-primary"
	case "Success":
		return "bg-success bg-opacity-10 text-success"
	case "Warning":
		return "bg-warning bg-opacity-10 text-warning"
	case "Error":
		return "bg-danger bg-opacity-10 text-danger"
	default:
		return ""
	}
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("01/02/2006 15:04:05")
}

func formatDuration(job models.Job) string {
	if job.Status == "Running" && job.LastRun != nil {
		duration := time.Since(*job.LastRun)
		job.Duration = &duration
	}
	if job.Duration == nil {
		return ""
	}
	return time.Unix(0, 0).UTC().Add(*job.Duration).Format("15:04:05")
}

func hasRunning(jobs []models.Job) bool {
	for _, j := range jobs {
		if j.Status == "Running" {
			return true
		}
	}
	return false
}
