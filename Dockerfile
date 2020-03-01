FROM golang:1.13.5-alpine3.10 as builder

RUN mkdir build
WORKDIR build
COPY go.mod .
COPY go.sum .


RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/server cmd/main.go

FROM alpine:3.10

COPY --from=builder /go/bin/server /app/server
USER 100

ENTRYPOINT ["/app/server"]
