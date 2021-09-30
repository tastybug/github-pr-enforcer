# github-pr-enforcer

This service is meant to be notified about pull requests on Github. It'll inspect the labels of the pull request and enforce configurable rules.

TODO: describe more how this is meant to support change failure rate metrics.

Notes:

* setting up a webhook receiver for PR related events can be set under a specific repo using the following URL template: https://github.com/USER/REPO/settings/hooks/new
* general info on how to build webhook receivers: https://docs.github.com/en/developers/webhooks-and-events/webhooks/about-webhooks