all: createmod
	go build jenkinsctl.go

createmod: clean
	cd jenkins && go mod init
	cd jenkins && go mod tidy
	go mod init
	go mod tidy

clean:
	rm -f go.sum go.mod jenkinsctl
	cd jenkins && rm -f go.mod go.sum
	go clean --modcache

build:
	rm -f jenkinsctl
	go build jenkinsctl.go
