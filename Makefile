# Read version from git tag
VERSION := $(shell git describe --tags | sed -e 's/^v//')
# Release artifacts are stored in ARTIFACTDIR
ARTIFACTDIR = artifacts
# Release packages
LINUX_PKG = when-$(VERSION).linux-amd64.tar.gz
WIN_PKG = when-$(VERSION).windows-amd64.zip

# Build release packages into the artifacts dir and calculate SHA-256 hashes.
.PHONY: release
release:
	@GOOS=linux GOARCH=amd64 go build -o "$(ARTIFACTDIR)/when"
	@cd "$(ARTIFACTDIR)" && tar --remove-files -czf $(LINUX_PKG) when
	@GOOS=windows GOARCH=amd64 go build -o "$(ARTIFACTDIR)/when.exe"
	@cd "$(ARTIFACTDIR)" && zip --quiet --move $(WIN_PKG) when.exe
	@cd "$(ARTIFACTDIR)" && sha256sum $(LINUX_PKG) $(WIN_PKG) > sha256sums.txt

# Remove artifacts directory.
.PHONY: clean
clean:
	@rm -rf "$(ARTIFACTDIR)"
