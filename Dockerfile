FROM golang:1.26-alpine AS builder

WORKDIR /src

# Install certificates so HTTPS downloads work reliably in build stage.
RUN apk add --no-cache ca-certificates

# Cache dependency download separately from source changes.
COPY go.mod ./
RUN go mod download

# Copy source and build a static Linux binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/app ./cmd

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /
COPY --from=builder /out/app /app

USER nonroot:nonroot
ENTRYPOINT ["/app"]
