# ---- Stage 1: build the Go binary ----
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /deals-api .

# ---- Stage 2: tiny runtime image ----
FROM alpine:3.20
RUN adduser -D -u 10001 app
USER app
WORKDIR /app
COPY --from=build /deals-api /app/deals-api
ENV GIN_MODE=release PORT=8000
EXPOSE 8000
ENTRYPOINT ["/app/deals-api"]
