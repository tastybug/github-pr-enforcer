# github-pr-enforcer

This service is meant to be notified about pull requests on Github. It'll inspect the labels of the pull request and enforce configurable rules.

TODO: describe more how this is meant to support change failure rate metrics.

### Local Dev Setup and Local Hosting of Webhook Mode

3 shells with the following commands as a local pipeline

1) test and build: `watchexec -e go 'make test build_github_webhook_localarch'`
2) restart local server: `watchexec -w enforcer -r ./enforcer`
3) run ping against local server: `curl -v -d "{\"zen\":\"Wisdom goes here\", \"repository\": {\"name\": \"github-pr-enforcer\", \"id\": 12345}}" localhost:9000/validate-pr` 

### Deploying to minis.fritz.box

`make deploy_devstage`

### Notes:

* setting up a webhook receiver for PR related events can be set under a specific repo using the following URL template: https://github.com/USER/REPO/settings/hooks/new
* general info on how to build webhook receivers: https://docs.github.com/en/developers/webhooks-and-events/webhooks/about-webhooks
