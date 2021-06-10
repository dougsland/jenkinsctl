build: clean
	go mod init || true
	go mod tidy || true
	go build jenkinsctl.go

clean:
	rm -f go.sum go.mod jenkinsctl
	go clean --modcache

all: build
