export GO111MODULE := on
export CGO_ENABLED=1

OS_NAME = $(shell uname -s | tr A-Z a-z)

# Go binary path
GOPATH := $(shell go env GOPATH 2>/dev/null || echo $$HOME/go)
GOBIN := $(GOPATH)/bin
PATH := $(GOBIN):$(PATH)
export PATH

# OpenNHP submodule directory
OPENNHP_DIR = third_party/opennhp

all: generate-version-and-build


generate-version-and-build:
	@echo "[StealthDNS] Start building..."
	@$(MAKE) init
	@$(MAKE) build-sdk
	@$(MAKE) build
	@echo "[StealthDNS] Build for platform ${OS_NAME} successfully done!"

init:
	@echo "[StealthDNS] Initializing..."
	@git clean -df release 2>/dev/null || true
	@git submodule update --init --recursive
	@go mod tidy

# Build OpenNHP SDK from submodule
build-sdk:
	@echo "[StealthDNS] Building OpenNHP SDK from submodule..."
ifeq ($(OS_NAME), linux)
	@$(MAKE) build-sdk-linux
else ifeq ($(OS_NAME), darwin)
	@$(MAKE) build-sdk-macos
else
	@echo "[StealthDNS] Skipping SDK build on ${OS_NAME}, use build.bat for Windows"
endif

build-sdk-linux:
	@echo "[StealthDNS] Building Linux SDK (nhp-agent.so)..."
	@cd $(OPENNHP_DIR)/nhp && go mod tidy
	@cd $(OPENNHP_DIR)/endpoints && go mod tidy
	@cd $(OPENNHP_DIR)/endpoints && \
		go build -a -trimpath -buildmode=c-shared -ldflags="-w -s" -v \
		-o ../../../sdk/nhp-agent.so ./agent/main/main.go ./agent/main/export.go
	@echo "[StealthDNS] Linux SDK built successfully!"
	@cd $(OPENNHP_DIR)/nhp && git restore go.mod go.sum 2>/dev/null || git checkout go.mod go.sum 2>/dev/null || true
	@cd $(OPENNHP_DIR)/endpoints && git restore go.mod go.sum 2>/dev/null || git checkout go.mod go.sum 2>/dev/null || true
	@cd $(OPENNHP_DIR) && git reset --hard HEAD 2>/dev/null || true

build-sdk-macos:
	@echo "[StealthDNS] Building macOS SDK (nhp-agent.dylib)..."
	@cd $(OPENNHP_DIR)/nhp && go mod tidy
	@cd $(OPENNHP_DIR)/endpoints && go mod tidy
	@cd $(OPENNHP_DIR)/endpoints && \
		GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 \
		go build -a -trimpath -buildmode=c-shared -ldflags="-w -s" -v \
		-o ../../../sdk/nhp-agent.dylib ./agent/main/main.go ./agent/main/export.go
	@echo "[StealthDNS] macOS SDK built successfully!"
	@cd $(OPENNHP_DIR)/nhp && git restore go.mod go.sum 2>/dev/null || git checkout go.mod go.sum 2>/dev/null || true
	@cd $(OPENNHP_DIR)/endpoints && git restore go.mod go.sum 2>/dev/null || git checkout go.mod go.sum 2>/dev/null || true
	@cd $(OPENNHP_DIR) && git reset --hard HEAD 2>/dev/null || true

build:
	@echo "[StealthDNS] Building package..."
	@mkdir -p ./release/etc/cert ./release/sdk
	go build -trimpath -ldflags="-w -s" -v -o ./release/stealth-dns ./main.go && \
	cp ./etc/*.toml ./release/etc/ && \
	cp ./sdk/nhp-agent.* ./release/sdk/ 2>/dev/null || true && \
	cp ./etc/cert/rootCA.pem ./release/etc/cert/ 2>/dev/null || true
ifeq ($(OS_NAME), darwin)
	install_name_tool -change nhp-agent.dylib ./sdk/nhp-agent.dylib ./release/stealth-dns
endif

# Windows SDK build (call from PowerShell/CMD)
build-sdk-windows:
	@echo "[StealthDNS] Building Windows SDK (nhp-agent.dll)..."
	cd $(OPENNHP_DIR)/nhp && go mod tidy
	cd $(OPENNHP_DIR)/endpoints && go mod tidy
	cd $(OPENNHP_DIR)/endpoints && \
		go build -a -trimpath -buildmode=c-shared -ldflags="-w -s" -v \
		-o ../../sdk/nhp-agent.dll ./agent/main/main.go ./agent/main/export.go
	@echo "[StealthDNS] Windows SDK built successfully!"

# Clean SDK binaries
clean-sdk:
	@echo "[StealthDNS] Cleaning SDK binaries..."
	rm -f sdk/nhp-agent.so sdk/nhp-agent.dylib sdk/nhp-agent.dll sdk/nhp-agent.h

# UI 构建相关目标
ui: ui-init ui-build
	@echo "[StealthDNS UI] Build for platform ${OS_NAME} successfully done!"

ui-init:
	@echo "[StealthDNS UI] Initializing..."
	@cd ui && go mod tidy
	@cd ui/frontend && npm ci || npm install
	@cd ui/frontend && git checkout package-lock.json 2>/dev/null || true

ui-build:
	@echo "[StealthDNS UI] Building UI package..."
ifeq ($(OS_NAME), windows)
	@WAILS_CMD=$$(command -v wails 2>/dev/null || echo "$(GOBIN)/wails"); \
	if [ ! -f "$$WAILS_CMD" ] && [ ! -x "$$WAILS_CMD" ]; then \
		echo "[StealthDNS] Error: wails command not found. Please run 'make check-wails' first."; \
		exit 1; \
	fi; \
	cd ui && PATH="$(GOBIN):$$PATH" $$WAILS_CMD build -platform windows/amd64 -o stealthdns-ui.exe
	cp ./ui/build/bin/stealthdns-ui.exe ./release/
else ifeq ($(OS_NAME), darwin)
	@WAILS_CMD=$$(command -v wails 2>/dev/null); \
	if [ -z "$$WAILS_CMD" ]; then \
		WAILS_CMD="$(GOBIN)/wails"; \
	fi; \
	if [ ! -f "$$WAILS_CMD" ]; then \
		echo "[StealthDNS] Error: wails command not found at $$WAILS_CMD. Please run 'make check-wails' first."; \
		exit 1; \
	fi; \
	cd ui && PATH="$(GOBIN):$$PATH" $$WAILS_CMD build -platform darwin/universal
	cp -r ./ui/build/bin/stealthdns-ui.app ./release/ 2>/dev/null || \
	cp ./ui/build/bin/stealthdns-ui ./release/
	@cd ui/frontend && git checkout wailsjs/ 2>/dev/null || true
else
	@WAILS_CMD=$$(command -v wails 2>/dev/null || echo "$(GOBIN)/wails"); \
	if [ ! -f "$$WAILS_CMD" ] && [ ! -x "$$WAILS_CMD" ]; then \
		echo "[StealthDNS] Error: wails command not found. Please run 'make check-wails' first."; \
		exit 1; \
	fi; \
	cd ui && PATH="$(GOBIN):$$PATH" $$WAILS_CMD build -platform linux/amd64
	cp ./ui/build/bin/stealthdns-ui ./release/
endif

ui-dev:
	@echo "[StealthDNS UI] Starting development mode..."
	cd ui && wails dev

# 完整构建：DNS 服务 + UI
full: all ui
	@echo "[StealthDNS] Full build completed!"

# 清理
clean:
	@echo "[StealthDNS] Cleaning..."
	git clean -df release
	rm -rf ui/build/bin
	rm -rf ui/frontend/dist
	rm -rf ui/frontend/node_modules

.PHONY: all generate-version-and-build init build build-sdk build-sdk-linux build-sdk-macos build-sdk-windows clean-sdk
