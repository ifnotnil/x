SHELL := /usr/bin/env bash

REPO_ROOT = $(shell cd .. && pwd)

include $(REPO_ROOT)/scripts/go.mk
include $(REPO_ROOT)/tools/tools.mk
include $(REPO_ROOT)/scripts/lib.mk

.PHONY: nginx
nginx:
	@docker compose --file ./compress/testdata/nginx/compose.yml up --force-recreate --detach

.PHONY: nginx-down
nginx-down:
	@docker compose --file ./compress/testdata/nginx/compose.yml down --volumes --remove-orphans --rmi local

.PHONY: test-integration
test-integration: TAGS=integration
test-integration: nginx test
