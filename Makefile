BINARY_NAME := ips
GO := go
LDFLAGS := -ldflags "-s -w"

PLATFORMS := \
	darwin/amd64 \
	darwin/arm64 \
	windows/386 \
	windows/amd64 \
	windows/arm64 \
	linux/386 \
	linux/arm \
	linux/amd64 \
	linux/arm64

UPX_PLATFORMS := \
	darwin/amd64 \
	windows/386 \
	windows/amd64 \
	linux/386 \
	linux/amd64

.PHONY: all clean lint tidy test build crossbuild upx

all: clean lint tidy test build

clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/

lint:
	@echo "🕵️‍♂️ Running linters..."
	golangci-lint run ./...

tidy:
	@echo "🧼 Tidying up dependencies..."
	$(GO) mod tidy

test:
	@echo "🧪 Running tests..."
	$(GO) test ./... -cover

build:
	@echo "🔨 Building for current platform..."
	CGO_ENABLED=0 $(GO) build -trimpath $(LDFLAGS) -o bin/$(BINARY_NAME) main.go

crossbuild: clean
	@echo "🌍 Building for multiple platforms..."
	for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d/ -f1); \
		arch=$$(echo $$platform | cut -d/ -f2); \
		float=$$(echo $$platform | cut -d/ -f3); \
		output_name=bin/ips_$${os}_$${arch}; \
		[ "$$float" != "" ] && output_name=$$output_name_$$float; \
		echo "🔨 Building for $$os/$$arch..."; \
		echo "🔨 Building for $$output_name..."; \
		GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 GOARM=$$float $(GO) build -trimpath $(LDFLAGS) -o $$output_name main.go ; \
		if [ "$(ENABLE_UPX)" = "1" ] && echo "$(UPX_PLATFORMS)" | grep -q "$$os/$$arch"; then \
			echo "⚙️ Compressing binary $$output_name..." && upx --best $$output_name; \
		fi; \
	done