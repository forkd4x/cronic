package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/goccy/go-yaml"
)

type Job struct {
	File string
	Name string
	Desc string
	Cron string
}

func NewJob(file string) (Job, error) {
	job := Job{
		File: file,
	}
	err := job.ParseFile()
	if err != nil {
		return job, fmt.Errorf("failed to parse file %s: %w", file, err)
	}

	return job, nil

}

func (job *Job) ParseFile() error {
	file, err := os.Open(job.File)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", job.File, err)
	}
	defer file.Close() // nolint

	buffer := make([]byte, 10240)
	_, err = file.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", job.File, err)
	}

	re := regexp.MustCompile(`(?m)^\W+\s+(\w+):\s*([^:\n$]*)\n`)
	matches := re.FindAllStringSubmatch(string(buffer), -1)
	if matches == nil || matches[0][1] != "cronic" {
		return nil
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
	fmt.Println(yamlData)

	err = yaml.Unmarshal([]byte(yamlData), job)
	if err != nil {
		return fmt.Errorf("failed to parse YAML %s: %w", yamlData, err)
	}
	return nil
}

func (job Job) Execute() error {
	cmd := exec.Command("sh", "-c", "./"+job.File)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
