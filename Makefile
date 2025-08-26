APP=daycount
PKG=./internal/...
COVER_PROFILE=coverage.out

.PHONY: test cover run serve docker-build docker-run lint

test:
	# Exclude cmd (CLI) package from coverage per project policy
	go test -race -count=1 -cover -coverprofile=$(COVER_PROFILE) $(PKG)

cover: test
	go tool cover -func=$(COVER_PROFILE)

run:
	go run ./cmd/daycount --help

serve:
	go run ./cmd/daycount --serve

docker-build:
	docker build -t daycount:latest .

docker-run: docker-build
	docker run --rm -p 8080:8080 daycount:latest

lint:
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		 echo "golangci-lint not found. Install: https://golangci-lint.run/usage/install/"; \
		 exit 2; \
	 fi
	golangci-lint run ./...
