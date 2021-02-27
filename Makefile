release:
ifeq ($(v), )
	@echo "[z] Usage: make release v=[version]"
else
	@echo ""
	@echo "[z] Preparing Zen realease v"$(v)
	@echo ""
	@make compile v="$(v)"
	@export V=$(v) && sh -c "'$(CURDIR)/scripts/github_release.sh'"
endif

compile:
	@echo ""
	@echo "[z] Compiling Zen for all platforms..."
	@echo ""
	@for GOOS in darwin linux windows freebsd; do \
		for GOARCH in 386 amd64; do \
			go build -v -o bin/zen$(v)-$$GOOS-$$GOARCH ; \
		done ; \
	done