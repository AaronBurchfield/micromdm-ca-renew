GOOS := "linux"
GOARCH := "amd64"
GOBUILD_LDFLAGS := "-s -w"

all: deps build

deps:
	go get

build-linux:
	GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags=${GOBUILD_LDFLAGS} -o ./build/renew-linux

build:
	go build -ldflags=${GOBUILD_LDFLAGS} -o ./build/renew

clean:
	rm -rf ./build
