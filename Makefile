release:
ifeq ($(v), )
	@echo "[M] Usage: make release v=[version]"
else
	@echo ""
	@echo "[M] Preparing Zen realease v"$(v)
	@echo ""
	@make compile v="$(v)"
	@export V=$(v) && sh -c "'$(CURDIR)/scripts/github_release.sh'"
endif

compile:
ifeq ($(v), )
	@echo "[M] Usage: make release v=[version]"
else
	@echo ""
	@echo "[M] Compiling Zen for all platforms..."
	@echo ""
	@for GOOS in darwin linux windows freebsd; do \
		for GOARCH in 386 amd64; do \
			go build -v -o bin/zen$(v)-$$GOOS-$$GOARCH cmd/zen/main.go; \
		done ; \
	done
endif