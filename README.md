## Cronic Scheduler (WIP)

```sh
go tool templ generate
go run . examples/
```

### Live reload
```sh
go tool templ generate --watch --proxy="http://localhost:1323" --cmd="go run . examples/"
```
Browser will automatically open and auto-refresh http://localhost:7331
