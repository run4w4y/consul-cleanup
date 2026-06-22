FROM golang:1.22-alpine AS builder

WORKDIR /proj

COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download && go mod verify

COPY . .
RUN mkdir -p dist
RUN go build -v -o ./dist ./...

FROM alpine:3.20 AS runner

COPY --from=builder /proj/dist/consul-cleanup /usr/local/bin/consul-cleanup

ENTRYPOINT [ "consul-cleanup" ]
