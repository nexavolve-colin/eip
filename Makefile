.PHONY: build

GOOS := $$(go env GOOS)
GOARCH := $$(go env GOARCH)

fmt:
	go fmt ./...
	go vet ./...

build:
	@mkdir -p bins
	go build -ldflags "-X main.version=${version} -X main.build=$$(git rev-parse HEAD)" -o bins/eip .
	tar -czvf bins/eip-${version}-${GOOS}-${GOARCH}.tar.gz bins/eip

test: fmt
	go test -race -cover -short
