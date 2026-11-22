.PHONY: check

# What the release workflow runs before it builds anything.
check:
	@test -z "$$(gofmt -l .)" || { echo "not formatted:"; gofmt -l .; exit 1; }
	go vet ./...
	@echo "All checks passed!"
