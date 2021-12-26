module github.com/bhojpur/piro/plugins/github-repo

go 1.17.5

replace github.com/bhojpur/piro => ../..

require (
	github.com/bradleyfalzon/ghinstallation v1.1.1
	github.com/bhojpur/piro v0.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.5.6
	github.com/google/go-github/v31 v31.0.0
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.36.0
	k8s.io/api v0.21.1
)
