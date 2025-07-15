package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/forkd4x/cronic/models"
	"github.com/go-co-op/gocron/v2"
	"github.com/goforj/godump"
	"github.com/google/uuid"
	"github.com/r3labs/sse/v2"
)

type Cronic struct {
	Context   context.Context
	Scheduler gocron.Scheduler
	Server    Server
	mu        sync.Mutex
	debounce  sync.Map
}

func NewCronic(root string) (*Cronic, error) {
	cronic := &Cronic{
		Context: context.Background(),
	}
	var err error
	if root != "." {
		err = os.Chdir(root)
		if err != nil {
			return cronic, fmt.Errorf("failed to chdir to %s: %w", root, err)
		}
	}
	err = models.Init()
	if err != nil {
		return cronic, err
	}
	cronic.Scheduler, err = gocron.NewScheduler()
	if err != nil {
		return cronic, fmt.Errorf("failed to initialize scheduler: %w", err)
	}
	cronic.Server = NewServer()
	return cronic, nil
}

func (cronic *Cronic) Start() {
	cronic.Scheduler.Start()
	go func() {
		err := cronic.Server.Echo.Start(":1323")
		if err != nil && err != http.ErrServerClosed {
			cronic.Server.Echo.Logger.Fatal("shutting down the server")
		}

	}()
	var quit context.CancelFunc
	cronic.Context, quit = signal.NotifyContext(cronic.Context, os.Interrupt)
	defer quit()
	<-cronic.Context.Done()
}

func (cronic *Cronic) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err1 := cronic.Server.Echo.Shutdown(ctx)
	err2 := cronic.Scheduler.Shutdown()
	return errors.Join(err1, err2)
}

func (cronic *Cronic) LoadJobs() error {
	dirEntries, err := os.ReadDir(".")
	if err != nil {
		cwd, err := os.Getwd()
		return fmt.Errorf("error reading directory %s: %w", cwd, err)
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			// TODO: Recursively load jobs in subdirectories
			continue
		}
		job := models.Job{
			File:   dirEntry.Name(),
			Status: "Pending",
		}
		err := job.ParseFile()
		if err != nil {
			return fmt.Errorf("failed to parse file %s: %w", job.File, err)
		}
		if job.Name == "" {
			fmt.Println("No cronic yaml in", job.File)
			continue
		}
		fmt.Println("Found cronic yaml in", job.File)
		godump.Dump(job)

		var dbJob models.Job
		where := "file = ? AND name = ?"
		result := models.DB.
			Where(where, job.File, job.Name).
			Order("updated_at desc").
			Limit(1).
			Find(&dbJob)
		if result.RowsAffected == 0 {
			result = models.DB.
				Where(strings.Replace(where, " AND ", " OR ", 1), job.File, job.Name).
				Order("updated_at desc").
				Limit(1).
				Find(&dbJob)
		}
		if result.RowsAffected == 0 {
			result := models.DB.Create(&job)
			if result.Error != nil {
				return fmt.Errorf("failed to insert job: %w", result.Error)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed querying database for job: %w", result.Error)
		} else {
			job.ID = dbJob.ID
			job.CreatedAt = dbJob.CreatedAt
			job.LastRun = dbJob.LastRun
			job.Duration = dbJob.Duration
			job.NextRun = dbJob.NextRun
			result := models.DB.Save(&job)
			if result.Error != nil {
				return fmt.Errorf("failed updating job: %w", result.Error)
			}
		}

		// TODO: Mark unfound jobs as deleted

		var scheduledJob gocron.Job
		scheduledJob, err = cronic.Scheduler.NewJob(
			gocron.CronJob(job.Cron, true),
			gocron.NewTask(
				func(filename string) {
					err := job.Run()
					if err != nil {
						panic(err)
					}
				},
				dirEntry.Name(),
			),
			gocron.WithEventListeners(
				gocron.BeforeJobRuns(
					func(jobID uuid.UUID, jobName string) {
						cronic.mu.Lock()
						defer cronic.mu.Unlock()
						now := time.Now()
						job.Status = "Running"
						job.LastRun = &now
						nextRuns, err := scheduledJob.NextRuns(2)
						if err != nil {
							job.NextRun = nil
						} else {
							job.NextRun = &nextRuns[1]
						}
						models.DB.Save(&job)

						jobs, err := models.GetJobs()
						if err != nil {
							panic(err)
						}
						html, err := cronic.Server.RenderTemplate("jobs.html", jobs)
						if err != nil {
							panic(err)
						}
						// fmt.Println("```")
						// fmt.Println(string(html))
						// fmt.Println("```")
						cronic.Publish(&sse.Event{
							Event: fmt.Append(nil, "table"),
							Data:  html,
						})

						// Update duration of running tasks
						go func() {
							for {
								time.Sleep(time.Second)
								html, err := cronic.Server.RenderTemplate("job.html", job)
								if err != nil {
									panic(err)
								}
								// fmt.Println("```")
								// fmt.Println(string(html))
								// fmt.Println("```")
								cronic.Publish(&sse.Event{
									Event: fmt.Append(nil, job.ID),
									Data:  html,
								})
								if job.Status != "Running" {
									break
								}
							}
						}()
					},
				),
				gocron.AfterJobRuns(
					func(jobID uuid.UUID, jobName string) {
						duration := time.Since(*job.LastRun)
						cronic.mu.Lock()
						defer cronic.mu.Unlock()
						job.Duration = &duration
						if nextRun, err := scheduledJob.NextRun(); err == nil {
							job.NextRun = &nextRun
						}
						models.DB.Save(&job)

						jobs, err := models.GetJobs()
						if err != nil {
							panic(err)
						}
						html, err := cronic.Server.RenderTemplate("jobs.html", jobs)
						if err != nil {
							panic(err)
						}
						// fmt.Println("```")
						// fmt.Println(string(html))
						// fmt.Println("```")
						cronic.Publish(&sse.Event{
							Event: fmt.Append(nil, "table"),
							Data:  html,
						})
					},
				),
			),
		)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func (cronic *Cronic) Publish(event *sse.Event) {
	key := string(event.Event)
	if timer, ok := cronic.debounce.Load(key); ok {
		timer.(*time.Timer).Stop()
	}
	// Remove all newline and carriage return characters to avoid multi-line SSE data issues,
	// then trim surrounding spaces so HTMX SSE swap matches correctly.
	event.Data = bytes.Map(func(r rune) rune {
		if r == '\n' || r == '\r' {
			return -1
		}
		return r
	}, event.Data)
	event.Data = bytes.TrimSpace(event.Data)

	cronic.debounce.Store(key, time.AfterFunc(50*time.Millisecond, func() {
		cronic.Server.SSE.Publish("updates", event)
	}))
}
