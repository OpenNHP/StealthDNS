export GO111MODULE := on

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
	cp ./etc/*.toml ./release/etc/
.PHONY: all generate-version-and-build init build
