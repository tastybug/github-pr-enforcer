test:
	cd internal/enforcer; go test
	cd cmd/webhookrcv; go test
build_mac:
	go build cmd/webhookrcv/main.go; mv main enforcer_darwin
build_amd64:
	GOOS=linux GOARCH=amd64 go build cmd/webhookrcv/main.go; mv main enforcer_amd64
deploy_hetzner: build_amd64
	ssh 65.108.84.79 'mv enforcer_amd64 enforcer_amd64_running'
	scp enforcer_amd64 65.108.84.79: