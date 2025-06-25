package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/forkd4x/cronic/models"
	"github.com/go-co-op/gocron/v2"
	"github.com/goforj/godump"
	"github.com/labstack/echo/v4"
)

type Cronic struct {
	Context   context.Context
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
	err = models.Init()
	if err != nil {
		return cronic, err
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
			result := models.DB.Save(&job)
			if result.Error != nil {
				return fmt.Errorf("failed updating job: %w", result.Error)
			}
		}

		_, err = cronic.Scheduler.NewJob(
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
		)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
