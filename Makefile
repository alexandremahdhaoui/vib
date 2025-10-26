.PHONY: install test-install

# Get the latest git tag to use as the version for the installer.
VERSION ?= $(shell git describe --tags --abbrev=0)

install:
	@echo "Installing vib version $(VERSION)..."
	@./cmd/vib-installer/vib-installer.sh $(VERSION)

test-install:
	@echo "Running installer test..."
	@./cmd/vib-installer/test.sh
