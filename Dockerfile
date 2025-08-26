# syntax=docker/dockerfile:1

# ---- build stage ----
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod .
RUN go mod download
COPY . .
# Build static binary
RUN CGO_ENABLED=0 go build -ldflags='-s -w' -o /out/daycount ./cmd/daycount

# ---- runtime stage ----
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/daycount /daycount
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/daycount", "--serve"]
