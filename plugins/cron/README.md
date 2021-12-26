This plugin starts jobs based on time.
For example:
```
tasks:
- spec: "@every 10m"
  repo: github.com/shashi-rai/test-repo:piro
```

See https://godoc.org/github.com/robfig/cron for more details about the time specification.
Have a look at the `Config` struct in `main.go` w.r.t the configuration format.
