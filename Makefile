SHELL := /usr/bin/env bash

REPO_ROOT = $(CURDIR)

#include $(REPO_ROOT)/scripts/go.mk
include $(REPO_ROOT)/tools/tools.mk
include $(REPO_ROOT)/scripts/lib.mk
