build: clean-modcache
	go mod init || true
	go mod tidy || true
	go build jenkinsctl.go

clean-modcache:
	rm -f go.sum go.mod
	go clean --modcache

all: build
