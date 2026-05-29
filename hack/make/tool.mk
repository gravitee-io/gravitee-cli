##@ 🛠️  Tools

GCTL_TOOL_MK_LOADED := 1

ROOT_DIR ?= $(shell git -C $(CURDIR) rev-parse --show-toplevel)
LOCALBIN ?= $(ROOT_DIR)/bin

ADDLICENSE ?= $(LOCALBIN)/addlicense
GORELEASER ?= $(LOCALBIN)/goreleaser

$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

.PHONY: install-tools
install-tools: $(LOCALBIN) ## Install every tool declared under hack/tools/ into bin/
	@cd $(ROOT_DIR)/hack/tools && \
	for item in $$(find . -mindepth 1 -type d); do \
		( cd $$item && \
		  TOOL=$$(grep -e '^tool ' go.mod | sed -e 's/tool //') && \
		  echo "Installing $$TOOL" && \
		  GOBIN=$(LOCALBIN) go install $$TOOL ); \
	done

.PHONY: clean-tools
clean-tools: ## Remove every installed tool binary from bin/
	@rm -rf $(LOCALBIN)

.PHONY: reinstall-tools
reinstall-tools: clean-tools install-tools ## Wipe bin/ and reinstall every tool (use after bumping versions in hack/tools/*/go.mod)
