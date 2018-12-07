.PHONY: build

fmt:
	go fmt ./...

vet:
	go vet ./...

build:
	@mkdir -p bins
	go build -ldflags "-X main.version=${version}-$$(git rev-parse HEAD)" -o bins/eip .

test: fmt vet
	go test -race -cover -short
