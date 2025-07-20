## Cronic Scheduler (WIP)

```sh
go test ./... -v
go run . examples/
```
Open http://localhost:1323

### Live reload
```sh
go install github.com/air-verse/air@latest

air \
    --tmp_dir "/tmp" \
    --build.bin "/tmp/cronic.air" \
    --build.cmd "go build -o /tmp/cronic.air" \
    examples/
```
