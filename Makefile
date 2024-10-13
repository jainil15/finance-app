SHELL=/bin/zsh
tailwind:
	@bunx tailwindcss -i ./web/static/css/input.css -o ./web/static/css/output.css 
build: tailwind
	@templ generate
	@go build -o tmp/financeapp cmd/main.go

test:
	@go test -v ./...

run: build
	@./tmp/financeapp

watch:
	@air -c .air.toml
