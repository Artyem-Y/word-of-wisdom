# Stage 1: Linting
FROM golangci/golangci-lint:v1.54.1 as linter

WORKDIR /app
COPY . .
COPY .golangci-lint.yml /app/.golangci-lint.yml

# Run the formatter (e.g., gofmt) to format the code
RUN gofmt -s -w .

# Run the golangci-lint tool using the configuration file
RUN golangci-lint run --config=/app/.golangci-lint.yml

# Stage 2: Build
FROM golang:1.20-alpine as builder
ARG BUILD_VERSION=1.0
RUN apk add --no-cache build-base

WORKDIR /build

COPY . .

# Copy the linted code from the linter stage
COPY --from=linter /app /build

RUN go mod download
RUN go test ./... -cover

RUN CGO_ENABLED=0 GOOS=linux go build -o /main main.go

# Stage 3: Final
# Create a minimal production image
FROM alpine:latest

# It's essential to regularly update the packages within the image to include security patches
RUN apk update && apk upgrade

# Reduce image size
RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

# Avoid running code as a root user
RUN adduser -D appuser
USER appuser

EXPOSE 8080
WORKDIR /app

# Create a writable directory for the binary
RUN mkdir -p /app/bin
COPY --from=builder /main /bin/app
COPY quotes.txt /app/quotes.txt

CMD ["/bin/app"]