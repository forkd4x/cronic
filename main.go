package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"unicode"

	"github.com/go-co-op/gocron/v2"
	"github.com/goccy/go-yaml"
	"github.com/goforj/godump"
)

type CronicJob struct {
	Name string
	Desc string
	Cron string
}

func parseFileHeader(data string) (CronicJob, error) {
	lines := strings.Split(string(data), "\n")

	// Find the "cronic:" line
	startIndex := -1
	for i, line := range lines {
		if strings.Contains(line, "cronic:") {
			startIndex = i + 1
			break
		}
	}
	if startIndex == -1 || startIndex >= len(lines) {
		return CronicJob{}, fmt.Errorf("cronic yaml not found")
	}

	// Determine the prefix characters from the first line after "cronic:"
	var prefix string
	for _, char := range lines[startIndex] {
		if !unicode.IsLetter(char) {
			prefix += string(char)
		} else {
			break
		}
	}

	// Remove the prefix from all subsequent lines and build YAML
	var yamlLines []string
	for i := startIndex; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, prefix) {
			yamlContent := line[len(prefix):]
			if strings.HasPrefix(strings.TrimSpace(yamlContent), "cron:") {
				parts := strings.SplitN(yamlContent, ":", 2)
				if len(parts) == 2 {
					cronValue := strings.TrimSpace(parts[1])
					if strings.Contains(cronValue, "*") && !strings.HasPrefix(cronValue, "\"") {
						yamlContent = parts[0] + ": \"" + cronValue + "\""
					}
				}
			}
			yamlLines = append(yamlLines, yamlContent)
		} else {
			break
		}
	}
	yamlData := strings.Join(yamlLines, "\n")

	var job CronicJob
	err := yaml.Unmarshal([]byte(yamlData), &job)
	if err != nil {
		return CronicJob{}, fmt.Errorf("failed to parse YAML: %w", err)
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
		file, err := os.Open(dirEntry.Name())
		if err != nil {
			panic(err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}()

		// Read the first 10 kB of the file looking for `cronic:` yaml
		buffer := make([]byte, 10240)
		n, err := file.Read(buffer)
		if err != nil {
			panic(err)
		}
		header := string(buffer[:n])
		if strings.Contains(header, "cronic:") {
			fmt.Println("Found cronic yaml in", dirEntry.Name())
			yaml, err := parseFileHeader(header)
			if err != nil {
				panic(err)
			}
			godump.Dump(yaml)
			_, err = s.NewJob(
				gocron.CronJob(yaml.Cron, true),
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
