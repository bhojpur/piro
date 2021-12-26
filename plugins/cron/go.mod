module github.com/bhojpur/piro/plugins/cron

go 1.17.5

replace github.com/bhojpur/piro => ../..

require (
	github.com/bhojpur/piro v0.0.0-00010101000000-000000000000
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.8.1
)
