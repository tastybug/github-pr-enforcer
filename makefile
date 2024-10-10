test:
	cd internal/enforcer/service; go test
	cd cmd/webhook-github; go test
build:
	go build cmd/webhook-github/main.go && mv main github_webhook_rcv_bin
	go build cmd/cmdline-github/main.go && mv main github_cmdline_bin
build_github_webhook_amd64:
	GOOS=linux GOARCH=amd64 go build cmd/webhook-github/main.go; mv main github_webhook_rcv_amd64_bin
build_github_cmdline_amd64:
	GOOS=linux GOARCH=amd64 go build cmd/cmdline-github/main.go; mv main github_cmdline_amd64_bin
deploy_devstage: build_github_webhook_amd64
	ssh minis.fritz.box 'if [[ -f github_webhook_rcv_amd64_bin ]]; then mv github_webhook_rcv_amd64_bin enforcer_amd64_running; fi'
	scp github_webhook_rcv_amd64_bin minis.fritz.box:
