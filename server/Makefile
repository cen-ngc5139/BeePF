GOARCH ?= amd64
GOOS ?= linux

build: 
	cd cmd ; CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_LDFLAGS='-g -lcapstone -static'   go build -tags=netgo,osusergo -gcflags "all=-N -l" -v  -o server


dlv: build
	cd cmd ; dlv --headless --listen=:2345 --api-version=2 exec ./server -- --config ../config.json

run: build
	cd cmd ; ./server --config ../config.json