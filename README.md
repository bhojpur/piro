# Bhojpur Piro - Continuous Integration

A pre-integrated CI/CD system driving many web scale applications and/or services running over the [Bhojpur.NET Platform](https://github.com/bhojpur/platform). On one side, it uses any standard Git repository and drive the _Jobs_ using either Kubernetes/Docker and/or NanoVMs on the other side. It knows no pipelines, just the _Jobs_ and each _Job_ could be a Unikernel application and/or Kubernetes **Pod**. What you do in that _Unikernel_ or _Pod_ is up to you. We do not impose any "declarative pipeline syntax" or some groovy scripting language. Instead, the Bhojpur Piro jobs have run Node, Golang, or Bash scripts in production environments.

Any standard Git repository (e.g. [Bhojpur Seam](https://seam.in.bhojpur.net)) could be integrated with the Bhojpur Piro. Though it is present already in your [Bhojpur.NET Platform](https://github.com/bhojpur/platform) instance, but different variants are built already to suit requirments of specific industry sector.

=======

---

- [Installation](#installation)
  - [GitHub](#github)
  - [Configuration](#configuration)
  - [OAuth](#oauth)
- [Setting up jobs](#setting-up-jobs)
  - [GitHub events](#gitHub-events)
- [Log Cutting](#log-cutting)
  - [GitHub events](#gitHub-events)
- [Command Line Interface](#command-line-interface)
  - [Installation](#installation-1)
  - [Usage](#usage)
- [Annotations](#annotations)
- [Attribution](#attribution)
- [Thank You](#thank-you)

---

## Installation

The easiest way to install the Bhojpur Piro is using its [Helm chart](helm/).
Clone this Git repository, cd into `helm/`, and install using following steps

```
helm dep update
helm upgrade --install piro .
```

### Git-hoster integration

The Bhojpur Piro integrates with standard Git hosting platforms using its plugin system.
Currently, Bhojpur Piro ships with support for GitHub only ([plugins/github-repo](https://repositories.github.com/bhojpur/piro/tree/cw/repo-plugins/plugins/github-repo) and [plugins/github-trigger](https://repositories.github.com/bhojpur/piro/tree/cw/repo-plugins/plugins/github-trigger)).

To add support for other Git hoster, the `github-repo` plugin is a good starting point.

#### GitHub

To use the Bhojpur Piro with the GitHub, you'll need a GitHub app.
To create the app, please [follow the steps here](https://developer.repositories.github.com/apps/building-github-apps/creating-a-github-app/).

When creating the app, please use following values:

| Parameter | Value | Description |
| --------- | ----------- | ------- |
| `User authorization callback URL` | `https://your-piro-installation.com/plugins/github-integration` | The `/plugins/github-integration`
path is important, the domain should match your installation's `config.baseURL` |
| `Webhook URL` | `https://your-piro-installation.com/plugins/github-integration` | The `/plugins/github-integration` path is important,
the domain should match your installation's `config.baseURL` |
| `Permissions` | Contents: Read-Only | |
| | Commit Status: Read & Write | |
| | Issues: Read & Write | |
| | Pull Requests: Read & Write | |
| `Events` | Meta | |
| | Push | |
| | Issue Comments | |

### Configuration

The following table lists the (incomplete set of) configurable parameters of the Bhojpur Piro chart and their default values. The Helm chart's `values.yaml` is the reference for chart's configuration surface.

| Parameter | Description | Default |
| --------- | ----------- | ------- |
| `repositories.github.webhookSecret` | Webhook Secret of your GitHub application. See [GitHub Setup](#github) | `my-webhook-secret` |
| `repositories.github.privateKeyPath` | Path to the private key for your GitHub application. See [GitHub setup](#github) | `secrets/github-app.com` |
| `repositories.github.appID` | AppID of your GitHub application. See [GitHub setup](#github) | `secrets/github-app.com` |
| `repositories.github.installationID` | InstallationID of your GitHub application. Have a look at the _Advanced_ page of your GitHub app to find thi s ID. | `secrets/github-app.com` |
| `config.baseURL` | URL of your Bhojpur Piro installation | `https://piro.bhojpur.net` |
| `config.timeouts.preparation` | Time a job can take to initialize | `10m` |
| `config.timeouts.total` | Total time a job can take | `60m` |
| `image.repository` | Image repository | `bhojpur/piro` |
| `image.tag` | Image tag | `latest` |
| `image.pullPolicy` | Image pull policy | `Always` |
| `replicaCount`  | Number of cert-manager replicas  | `1` |
| `rbac.create` | If `true`, create and use RBAC resources | `true` |
| `resources` | CPU/memory resource requests/limits | |
| `nodeSelector` | Node labels for pod assignment | `{}` |
| `affinity` | Node affinity for pod assignment | `{}` |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example,

```console
helm install --name my-release -f values.yaml .
```

> **Tip**: You can use the default [values.yaml](values.yaml)

### OAuth

The Bhojpur Piro does not support OAuth by itself. However, using [OAuth Proxy](https://github.com/oauth2-proxy/oauth2-proxy) that's easy enough to add. It could leverage [Bhojpur Web](https://github.com/bhojpur/web) for full fledged support for enterprise grade product features.

## Setting up Jobs

The Bhojpur Piro _jobs_ are files in your Git repository where one file represents one Job.
A Bhojpur Piro job file mainly consists of the [PodSpec](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podspec-v1-core)
that will be run. The Bhojpur Piro will add an `/application` mount to your _pod_ where you'll find the checked out repository the job is running on.

For example:

```YAML
pod:
  containers:
  - name: hello-world
    image: alpine:latest
    workingDir: /application
    imagePullPolicy: IfNotPresent
    command:
    - sh 
    - -c
    - |
      echo Hello World
      ls
```

This job would print Hello World and list all files in the root of the repository.

Checkout [Bhojpur Piro's own build job](.piro/build-job.yaml) for a more complete example.

> **Tip**: You can use the Bhojpur Piro CLI to create a new _Job_ using `piro init job`

### GitHub events

The Bhojpur Piro starts _jobs_ based on GitHub push events, if the repository contains a `.piro/config.yaml` file, e.g.

```YAML
defaultJob: ".piro/build-job.yaml"
rules:
- path: ".piro/deploy.yaml"
  matchesAll:
  - or: ["repo.ref ~= refs/tags/"]
  - or: ["trigger !== deleted"]
```

The example above starts `.piro/deploy.yaml` for all tags. For everything else it will start `.piro/build-job.yaml`.

## Log Cutting

The Bhojpur Piro extracts structure from the log output its jobs produce. We call this process log cutting, because the Bhojpur Piro understands logs as a bunch of streams/slices which have to be demultiplexed.

The default cutter in the Bhojpur Piro expects the following syntax:

| Code | Command | Description |
| --------- | ----- | ----------- |
| `[someID\|PHASE] Some description here` | Enter new phase | Enters into a new phase identified by `someID` and described by `Some description here`. All output in this phase that does not explicitely name a slice will use `someID` as slice.
| `[someID] Arbitrary output` | Log to a slice | Logs `Arbitrary output` and marks it as part of the `someID` slice.
| `[someID\|DONE]` | Finish a slice | Marks the `someID` slice as done. No more output is expected from this slice in this phase.
| `[someID\|FAIL] Reason` | Fail a slice | Marks the `someID` slice as failed becuase of `Reason`. No more output is expected from this slice in this phase. Failing a slice does not automatically fail the job.
| `[type\|RESULT] content` | Publish a result | Publishes `content` as result of type `type`

> **Tip**: You can produce this kind of log output using the Bhojpur Piro CLI: `piro log`

## Command Line Interface

The Bhojpur Piro sports a powerful continuous integration capability, which can be used to create, list, start and listen to jobs (e.g. Unikernel applications or services) applied in complex data processing.
=======

### Installation

The Bhojpur Piro CLI is available on the [GitHub release page](https://repositories.github.com/bhojpur/piro/releases), or using this one-liner:

```bash
curl -L bhojpur.net/get-cli.sh | sh
```

### Usage

```
The Bhojpur Piro is a very simple GitHub triggered, Unikernel and/or Kubernetes powered CI system.

Usage:
  piro [command]

Available Commands:
  help        Help about any command
  init        Initializes configuration for Bhojpur Piro
  job         Interacts with currently running or previously run jobs
  log         Prints log-cuttable content
  run         Starts the execution of a job
  version     Prints the version of this binary

Flags:
  -h, --help          help for Bhojpur Piro
      --host string   the Bhojpur Piro host to talk to (defaults to PIRO_HOST env var) (default "localhost:7777")
      --verbose       en/disable verbose logging

Use "piro [command] --help" for more information about a command.
```

## Annotations

Annotations are used by your Bhojpur Piro _job_ to make runtime decisions. The Bhojpur Piro supports passing annotation in three ways:

1. From Pull Request description

You can add annotations in the following form to your Pull Request description and the Bhojpur Piro will pick them up

```sh
/piro someAnnotation
/piro someAnnotation=foobar
- [x] /piro someAnnotation
- [x] /piro someAnnotation=foobar
```

2. From Git commit

The Bhojpur Piro supports same format as above to pass annotations via commit message. The Bhojpur Piro will use the top most commit only.

3. From CLI

```sh
piro run github -a someAnnotation=foobar
```

## Thank You

Thank you to our contributors.
