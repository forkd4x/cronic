package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"

	"github.com/go-co-op/gocron/v2"
	"github.com/goccy/go-yaml"
	"github.com/goforj/godump"
)

type CronicJob struct {
	Name string
	Desc string
	Cron string
}

func parseFile(filename string) (CronicJob, error) {
	file, err := os.Open(filename)
	if err != nil {
		return CronicJob{}, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close() // nolint

	buffer := make([]byte, 10240)
	_, err = file.Read(buffer)
	if err != nil {
		return CronicJob{}, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	re := regexp.MustCompile(`(?m)^\W+\s+(\w+):\s*([^:\n$]*)\n`)
	matches := re.FindAllStringSubmatch(string(buffer), -1)
	if matches == nil || matches[0][1] != "cronic" {
		return CronicJob{}, nil
	}

	var yamlLines []string
	for _, match := range matches {
		if len(match) != 3 || match[1] == "cronic" {
			continue
		}
		key, value := match[1], match[2]
		// Quote cron expressions containing asterisks if not already quoted
		if key == "cron" && strings.Contains(value, "*") && !strings.HasPrefix(value, "\"") {
			value = "\"" + value + "\""
		}
		yamlLines = append(yamlLines, key+": "+value)
	}
	yamlData := strings.Join(yamlLines, "\n")

	var job CronicJob
	err = yaml.Unmarshal([]byte(yamlData), &job)
	if err != nil {
		return CronicJob{}, fmt.Errorf("failed to parse YAML %s: %w", yamlData, err)
	}
	return job, nil
}

func execute(filename string) error {
	cmd := exec.Command("sh", "-c", "./"+filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	err := os.Chdir("examples")
	if err != nil {
		panic(err)
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	dirEntries, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		job, err := parseFile(dirEntry.Name())
		if err != nil {
			panic(err)
		}
		if job.Name == "" {
			fmt.Println("No cronic yaml in", dirEntry.Name())
			continue
		}
		fmt.Println("Found cronic yaml in", dirEntry.Name())
		godump.Dump(job)
		_, err = s.NewJob(
			gocron.CronJob(job.Cron, true),
			gocron.NewTask(
				func(filename string) {
					err := execute(filename)
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

	s.Start()

	ctx, quit := signal.NotifyContext(context.Background(), os.Interrupt)
	defer quit()
	<-ctx.Done()

	err = s.Shutdown()
	if err != nil {
		panic(err)
	}
}
