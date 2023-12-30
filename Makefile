SHELL := /usr/bin/env bash

.SHELLFLAGS = -e -o pipefail

# 현재 디렉토리 저장
CURRENT_DIR := $(shell pwd)

# 릴리스 함수
release:
	@echo "Releasing $(TAG)"
	@cd $(CURRENT_DIR) && \
	go mod vendor && \
	rm -rf $(CURRENT_DIR)/dist $(CURRENT_DIR)/gopath && \
	export GOPATH=$(CURRENT_DIR)/gopath && \
	if [ -z "$(TAG)" ]; then \
		echo "not found tag name"; \
		exit 1; \
	fi && \
	git tag -a $(TAG) -m "Add $(TAG)" && \
	git push origin $(TAG) && \
	goreleaser release --clean

# 기본적으로 실행할 타겟
all: release
