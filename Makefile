export GO111MODULE := on
export CGO_ENABLED=1

OS_NAME = $(shell uname -s | tr A-Z a-z)

all: generate-version-and-build


generate-version-and-build:
	@echo "[StealthDNS] Start building..."
	@$(MAKE) init
	@$(MAKE) build
	@echo "[StealthDNS] Build for platform ${OS_NAME} successfully done!"

init:
	@echo "[StealthDNS] Initializing..."
	git clean -df release
	go mod tidy

build:
	@echo "[StealthDNS] Building package..."
	go build -trimpath -ldflags="-w -s" -v -o ./release/stealth-dns ./main.go && \
	cp ./etc/*.toml ./release/etc/ && \
	cp ./sdk/nhp-agent.* ./release/sdk/
ifeq ($(OS_NAME), darwin)
	install_name_tool -change nhp-agent.dylib ./sdk/nhp-agent.dylib ./release/stealth-dns
endif

.PHONY: all generate-version-and-build init build
