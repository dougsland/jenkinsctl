build: clean
	cd jenkins && go mod init
	cd jenkins && go mod tidy
	go mod init
	go mod tidy
	go build jenkinsctl.go

clean:
	rm -f go.sum go.mod jenkinsctl
	cd jenkins && rm -f go.mod go.sum
	go clean --modcache

all: build
