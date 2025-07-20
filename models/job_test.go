package models

import (
	"testing"

	"github.com/goforj/godump"
)

func TestParseFile(t *testing.T) {
	wants := []Job{
		{
			File: "example1.go",
			Name: "Example Go Job",
			Desc: "Say hello every 4 seconds",
			Cron: "*/4 * * * * *",
		},
		{
			File: "example2.py",
			Name: "Example Python Job",
			Desc: "Say hello every 6 seconds",
			Cron: "*/6 * * * * *",
		},
		{
			File: "example3.sh",
			Name: "Example Shell Job",
			Desc: "Say hello every 10 seconds",
			Cron: "*/10 * * * * *",
		},
		{
			File: "example4.Dockerfile",
			Name: "Example Dockerfile Job",
			Desc: "Say hello every 14 seconds",
			Cron: "*/14 * * * * *",
		},
		{
			File: "ignore.go",
			Name: "",
			Desc: "",
			Cron: "",
		},
	}

	for _, want := range wants {
		t.Run(want.File, func(t *testing.T) {
			job := Job{File: "../examples/" + want.File}
			err := job.ParseFile()
			if err != nil {
				t.Errorf("Job.ParseFile(%s) error: %v", want.File, err)
			}
			job.File = want.File
			if job != want {
				t.Errorf(
					"Job.ParseFile(%s) fail:\nGot: %v\nWant:%v",
					want.File, godump.DumpStr(job), godump.DumpStr(want),
				)
			}
		})
	}
}
