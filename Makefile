GO ?= go

# Heavy targets run through the shared hosts' cpu-policy tool (load gate,
# machine-wide lock, width cap). Where the tool is not installed this expands
# to nothing and the targets run as before.
CPU_POLICY_RUN := $(shell command -v cpu-policy >/dev/null 2>&1 && echo cpu-policy run --)

# The README hero tooling is its own module (examples/hero), so its banner
# dependency never reaches iqlink's go.mod.
.PHONY: hero hero-svg
hero: ## Play the README hero in the terminal
	env GOWORK=off $(GO) -C examples/hero run .

hero-svg: ## Regenerate docs/hero.svg
	env GOWORK=off $(GO) -C examples/hero generate .

.PHONY: hero-check hero-check-inner
hero-check: ## Check real fixtures and generated hero freshness
	$(CPU_POLICY_RUN) $(MAKE) hero-check-inner

hero-check-inner:
	env GOWORK=off $(GO) -C examples/hero test ./...
