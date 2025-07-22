package models

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/lnquy/cron"
	"gorm.io/gorm"
)

type Job struct {
	ID        uint `gorm:"primaryKey"`
	File      string
	Name      string
	Desc      string
	Cron      string
	Cmd       string
	Status    string
	LastRun   *time.Time
	Duration  *time.Duration
	NextRun   *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
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

	re := regexp.MustCompile(`(?m)^\W+(?:\s(cronic)|\s\s+(\w+)):\s*([^:\n]*)\n`)
	matches := re.FindAllStringSubmatch(string(buffer), -1)
	if matches == nil || matches[0][1] != "cronic" {
		return nil
	}

	var yamlLines []string
	for _, match := range matches {
		if len(match) != 4 || match[1] == "cronic" {
			continue
		}
		key, value := match[2], match[3]
		// Quote cron expressions containing asterisks if not already quoted
		if key == "cron" && strings.Contains(value, "*") && !strings.HasPrefix(value, "\"") {
			value = "\"" + value + "\""
		}
		yamlLines = append(yamlLines, key+": "+value)
	}
	yamlData := strings.Join(yamlLines, "\n")

	err = yaml.Unmarshal([]byte(yamlData), job)
	if err != nil {
		return fmt.Errorf("failed to parse YAML %s: %w", yamlData, err)
	}
	if job.Cmd == "" {
		job.Cmd = "./" + job.File
	}
	return nil
}

func (job Job) CronHuman() string {
	exprDesc, err := cron.NewDescriptor()
	if err != nil {
		return job.Cron
	}
	desc, err := exprDesc.ToDescription(job.Cron, cron.Locale_en)
	if err != nil {
		return job.Cron
	}
	return desc
}

func GetJobs() ([]Job, error) {
	var jobs []Job
	orderBy := `
		case when status='Running' then last_run else null end desc,
		ifnull(next_run, '2099-01-01')
	`
	result := DB.Order(orderBy).Find(&jobs)
	return jobs, result.Error
}
