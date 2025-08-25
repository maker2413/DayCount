APP=daycount
PKG=./...
COVER_PROFILE=coverage.out

.PHONY: test cover run

test:
	go test -race -count=1 -cover -coverprofile=$(COVER_PROFILE) $(PKG)

cover: test
	go tool cover -func=$(COVER_PROFILE)

run:
	go run ./cmd/daycount --help
