package models

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Run struct {
	ID        uint `gorm:"primaryKey"`
	JobID     uint
	Status    string // TODO: Enum?
	StartedAt time.Time
	EndedAt   *time.Time
}

func (job *Job) Run() error {
	run := Run{
		JobID:     job.ID,
		Status:    "Running",
		StartedAt: time.Now(),
	}
	if err := DB.Create(&run).Error; err != nil {
		return fmt.Errorf("failed to insert Run: %w", err)
	}
	var cmd *exec.Cmd
	if filepath.Ext(job.File) == ".Dockerfile" {
		tag := strings.TrimSuffix(filepath.Base(job.File), filepath.Ext(job.File))
		cmd = exec.Command("sh", "-c", "docker build -f "+job.File+" -t "+tag+" . && docker run "+tag)
	} else {
		cmd = exec.Command("sh", "-c", "./"+job.File)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start run: %w", err)
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			DB.Create(&Log{
				JobID:   job.ID,
				RunID:   run.ID,
				Type:    "stdout",
				Message: line,
			})
		}
	}()

	run.Status = "Success"
	if err := cmd.Wait(); err != nil {
		DB.Create(&Log{
			JobID:   job.ID,
			RunID:   run.ID,
			Type:    "stderr",
			Message: err.Error(),
		})
		run.Status = "Error"
	}
	end := time.Now()
	run.EndedAt = &end
	if err := DB.Save(&run).Error; err != nil {
		return fmt.Errorf("failed to update JobRun: %w", err)
	}
	job.Status = run.Status
	return nil
}
