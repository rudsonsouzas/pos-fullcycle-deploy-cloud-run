FROM golang:1.22.6-alpine AS build

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY . .

RUN go mod download
RUN go build -o api-server ./cmd/main.go

# Test stage
FROM builder AS test
RUN go test -v ./...

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/api-server .

EXPOSE 8080

CMD ["./api-server"]