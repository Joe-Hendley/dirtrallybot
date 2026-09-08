default: run
.PHONY: test build run generate

# Regenerate the templ views (internal/bot/web/*_templ.go). The generated files
# are committed, so this is only needed after editing a .templ file. Requires
# the templ CLI: go install github.com/a-h/templ/cmd/templ@latest
generate:
	templ generate

test:
	go test ./...

build:
	go build ./...

run:
	go run ./cmd/dirtrallybot
