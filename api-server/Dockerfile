FROM golang:1.22.6-alpine AS build

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api-server ./cmd/main.go

# Test stage
FROM build AS test
RUN go test -v ./...

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/api-server .

EXPOSE 8080

CMD ["./api-server"]