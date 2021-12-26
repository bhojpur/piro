# Bhojpur Piro - GitHub Integration
This plugin provides GitHub integration using a GitHub app. It can trigger builds when a branch is pushed,
pull annotations from a PR, take commands from a PR comment and update the commit status.

## Installation
First you must create a GitHub app with the following permissions:
- Deployments: Read & Write
- Issues: Read & Write
- Metadata: Read-only
- Pull Requests: Read & Write
- Commit Status: Read & Write

subscribing to the following events:
- Meta
- Issue Comment
- Push
- Pull Request

Once you have created this application, please install it on the repositories you intent to use Bhojpur Piro with.

Then add the following to your Bhojpur Piro config file:
```YAML
plugins:
  - name: "github-integration"
    type:
    - integration
    config:
      baseURL: https://your-piro-installation-url.com
      webhookSecret: choose-a-sensible-secret-here
      privateKeyPath: path-to-your/app-private-key.pem
      appID: 00000              # appID of your GitHub app
      installationID: 0000000   # installation ID of your GitHub app installation
      pullRequestComments:
        enabled: true
        updateComment: true
        requiresOrg: []
        requiresWriteAccess: true
```

## PR Commands
This integration plugin listens for comments on PRs to trigger operations in Bhojpur Piro.

```YAML
# start a Bhojpur Piro job for this PR
/piro run
```

## Commit Checks
For all jobs that carry the `updateGitHubStatus` annotation, the Bhojpur Piro attempts to add a commit check on the repository
pointed to in that annotation e.g. if the job ran with `updateGitHubStatus=bhojpur/piro`, upon completion of that job, this
plugin would add a check indiciating job success or failure.
By default, all jobs started using this integration plugin (push events or comments) will carry this annotation.

In addition to the job success annotation, jobs can add additional checks to commits. Any result posted to the `github`
or any channel starting with `github-check-` will become a check on the commit. For example:

- ```
  piro log result -d "dev installation" -c github url http://foobar.com
  ``` 
  adds the following check (`Details` points to `http://foobar.com`)

- ```
  piro log result -d "dev installation" -c github-check-tests conclusion success
  ```
  would add a successful check named `continuous-integration/piro/result-tests`, whereas 
  ```
  piro log result -d "dev installation" -c github-check-tests conclusion failure
  ```
  would add a failed check named `continuous-integration/piro/result-tests`.

  Valid values for `conclusion` results in this case are listed in the [GitHub API docs](https://docs.github.com/en/rest/reference/checks#update-a-check-run).
