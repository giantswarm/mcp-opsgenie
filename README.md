# General Go template repository

This is a general template repository containing some basic files every GitHub repo owned by Giant Swarm should have.

Note also these more specific repositories:

- [template-app](https://github.com/giantswarm/template-app)
- [gitops-template](https://github.com/giantswarm/gitops-template)
- [python-app-template](https://github.com/giantswarm/python-app-template)

## Creating a new repository

Please do not use the `Use this template` function in the GitHub web UI.

Check out the according [handbook article](https://handbook.giantswarm.io/docs/dev-and-releng/repository/go/) for better instructions.

### Some suggestions for your README

After you have created your new repository, you may want to add some of these badges to the top of your README.

- **CircleCI:** After enabling builds for this repo via [this link](https://circleci.com/setup-project/gh/giantswarm/REPOSITORY_NAME), you can find badge code on [this page](https://app.circleci.com/settings/project/github/giantswarm/REPOSITORY_NAME/status-badges).

- **Go reference:** use [this helper](https://pkg.go.dev/badge/) to create the markdown code.

- **Go report card:** enter the module name on the [front page](https://goreportcard.com/) and hit "Generate report". Then use this markdown code for your badge: `[![Go report card](https://goreportcard.com/badge/github.com/giantswarm/REPOSITORY_NAME)](https://goreportcard.com/report/github.com/giantswarm/REPOSITORY_NAME)`

- **Sourcegraph "used by N projects" badge**: for public Go repos only: `[![Sourcegraph](https://sourcegraph.com/github.com/giantswarm/REPOSITORY_NAME/-/badge.svg)](https://sourcegraph.com/github.com/giantswarm/REPOSITORY_NAME)`
