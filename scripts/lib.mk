# require variable REPO_ROOT

## https://www.gnu.org/software/make/manual/html_node/Parallel-Disable.html
.NOTPARALLEL:

# man git-clean
.PHONY: git-reset
git-reset:
	git reset --hard
	git clean -fd

.PHONY: env
env:
	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Module \e[0m \e[0;90m<<<\e[0m"
	@echo "$(MODULE)"
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Go env \e[0m \e[0;90m<<<\e[0m"
	go env
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Packages \e[0m \e[0;90m<<<\e[0m"
	$(GO_PACKAGES)
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Folders \e[0m \e[0;90m<<<\e[0m"
	$(GO_FOLDERS)
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Files \e[0m \e[0;90m<<<\e[0m"
	$(GO_FILES)
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Tools \e[0m \e[0;90m<<<\e[0m"
	@echo '$(TOOLS_BIN)'
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Path \e[0m \e[0;90m<<<\e[0m"
	@echo "$${PATH}" | tr ':' '\n'
	@echo ""

	@echo -e "\e[0;90m>>>\e[0m \e[0;94m Shell \e[0m \e[0;90m<<<\e[0m"
	@echo "SHELL         :$${SHELL}"
	@echo "BASH          :$${BASH}"
	@echo "BASH_VERSION  :$${BASH_VERSION}"
	@echo "BASH_VERSINFO :$${BASH_VERSINFO}"
	@echo ""

.PHONY: checks
checks: vet staticcheck golangci-lint

.PHONY: ci-fmt
ci-fmt: golangci-lint-fmt
	$(REPO_ROOT)/scripts/git-check-dirty

.PHONY: ci-mod
ci-mod: mod
	$(REPO_ROOT)/scripts/git-check-dirty

.PHONY: ci-sh
ci-sh: shfmt shellcheck
	@$(REPO_ROOT)/scripts/git-check-dirty

.PHONY: ci-dependabot
ci-dependabot:
	@$(REPO_ROOT)/scripts/dependabot-file
	@$(REPO_ROOT)/scripts/git-check-dirty
