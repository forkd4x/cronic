package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-co-op/gocron/v2"
	"github.com/goforj/godump"
)

type Cronic struct {
	Context   context.Context
	Scheduler gocron.Scheduler
}

func NewCronic(root string) (Cronic, error) {
	cronic := Cronic{
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
	return cronic, nil
}

func (cronic Cronic) Start() {
	cronic.Scheduler.Start()
	var quit context.CancelFunc
	cronic.Context, quit = signal.NotifyContext(cronic.Context, os.Interrupt)
	defer quit()
	<-cronic.Context.Done()
}

func (cronic Cronic) Shutdown() error {
	err := cronic.Scheduler.Shutdown()
	return err
}

func (cronic Cronic) LoadJobs() error {
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
		job := CronicJob{
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
