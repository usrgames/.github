.PHONY: check media

# What the release workflow runs before it builds anything.
check:
	@test -z "$$(gofmt -l .)" || { echo "not formatted:"; gofmt -l .; exit 1; }
	go vet ./...
	@echo "All checks passed!"

# The recording, as a game makes it for the site: hello.gif and hello.png in the folder MEDIA names.
MEDIA ?= dist/media
media:
	go run ./cmd/media -out $(MEDIA)
