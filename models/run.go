package models

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Run struct {
	ID        uint `gorm:"primaryKey"`
	JobID     uint
	Status    string // TODO: Enum?
	StartedAt time.Time
	EndedAt   *time.Time
}

func (job Job) Run() error {
	run := Run{
		JobID:     job.ID,
		Status:    "Running",
		StartedAt: time.Now(),
	}
	if err := DB.Create(&run).Error; err != nil {
		return fmt.Errorf("failed to insert JobRun: %w", err)
	}
	cmd := exec.Command("sh", "-c", "./"+job.File)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run job: %w", err)
	}
	// TODO: Check return code and/or stderr
	run.Status = "Success"
	end := time.Now()
	run.EndedAt = &end
	if err := DB.Save(&run).Error; err != nil {
		return fmt.Errorf("failed to update JobRun: %w", err)
	}
	return nil
}
