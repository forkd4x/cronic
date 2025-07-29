package models

import (
	"os"
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
			Cmd:  "./example1.go",
		},
		{
			File: "example2.py",
			Name: "Example Python Job",
			Desc: "Say hello every 6 seconds",
			Cron: "*/6 * * * * *",
			Cmd:  "./example2.py",
		},
		{
			File: "example3.sh",
			Name: "Example Shell Job",
			Desc: "Say hello every 10 seconds",
			Cron: "*/10 * * * * *",
			Cmd:  "./example3.sh",
		},
		{
			File: "example4.Dockerfile",
			Name: "Example Dockerfile Job",
			Desc: "Say hello every 14 seconds",
			Cron: "*/14 * * * * *",
			Cmd:  "docker build -f $f -t ${f%.*} . && docker run --rm ${f%.*}",
		},
		{
			File: "example5.py",
			Name: "Example Python2 Job",
			Desc: "Say hello every 22 seconds",
			Cron: "*/22 * * * * *",
			Cmd:  "./example5.py",
		},
		{
			File: "example6.php",
			Name: "Example PHP5 Job",
			Desc: "Say hello every 26 seconds",
			Cron: "*/26 * * * * *",
			Cmd:  "docker run --rm -v .:/app php:5.6-cli php /app/$f",
		},
		{
			File: "ignore.go",
			Name: "",
			Desc: "",
			Cron: "",
			Cmd:  "",
		},
	}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %q", err)
	}
	err = os.Chdir("../examples")
	if err != nil {
		t.Fatalf("Error changing to examples directory: %q", err)
	}
	defer func() {
		err := os.Chdir(dir)
		if err != nil {
			t.Fatalf("Error changing back to %q: %q", dir, err)
		}
	}()

	for _, want := range wants {
		t.Run(want.File, func(t *testing.T) {
			job := Job{File: want.File}
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
