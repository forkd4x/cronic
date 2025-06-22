package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/forkd4x/cronic/models"
	"github.com/go-co-op/gocron/v2"
	"github.com/goforj/godump"
	"github.com/labstack/echo/v4"
)

type Cronic struct {
	Context   context.Context
	Jobs      []models.Job
	Scheduler gocron.Scheduler
	Server    *echo.Echo
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
	cronic.Scheduler, err = gocron.NewScheduler()
	if err != nil {
		return cronic, fmt.Errorf("failed to initialize scheduler: %w", err)
	}
	cronic.Server = NewServer(cronic)
	return cronic, nil
}

func (cronic *Cronic) Start() {
	cronic.Scheduler.Start()
	go func() {
		err := cronic.Server.Start(":1323")
		if err != nil && err != http.ErrServerClosed {
			cronic.Server.Logger.Fatal("shutting down the server")
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
	err1 := cronic.Server.Shutdown(ctx)
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
			File: dirEntry.Name(),
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
		cronic.Jobs = append(cronic.Jobs, job)
		_, err = cronic.Scheduler.NewJob(
			gocron.CronJob(job.Cron, true),
			gocron.NewTask(
				func(filename string) {
					err := job.Execute()
					if err != nil {
						panic(err)
					}
				},
				dirEntry.Name(),
			),
		)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
