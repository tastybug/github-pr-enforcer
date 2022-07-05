test:
	cd internal/enforcer; go test
	cd cmd/webhookrcv; go test
build:
	go build cmd/webhookrcv/main.go; mv main enforcer
build_amd64:
	GOOS=linux GOARCH=amd64 go build cmd/webhookrcv/main.go; mv main enforcer_amd64
deploy_hetzner: build_amd64
	ssh 135.181.91.94 'if [[ -f enforcer_amd64 ]]; then mv enforcer_amd64 enforcer_amd64_running; fi'
	scp enforcer_amd64 135.181.91.94:
