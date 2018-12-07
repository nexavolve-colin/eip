.PHONY: build

build:
	@mkdir -p bins
	go build -ldflags "-X main.version=$$(git rev-parse HEAD)" -o bins/eip .

test:
	go test
