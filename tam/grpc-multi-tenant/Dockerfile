# Build stage
FROM golang:1.24.1-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY tam/grpc-multi-tenant/. ./tam/grpc-multi-tenant/

# CMD ["go", "run", "tam/grpc-multi-tenant/cmd/app/main.go"]

RUN go mod download
RUN CGO_ENABLED=0 go build -o app ./tam/grpc-multi-tenant/cmd/app/main.go


FROM alpine:latest

RUN apk update

WORKDIR /app/

COPY --from=builder /app/app .

CMD ["/app/app"]