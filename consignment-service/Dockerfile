FROM golang:1.9.0 as builder

WORKDIR /go/src/docker/consignment-service

COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep init -v
RUN dep ensure -v
RUN CGO_ENABLED=1 GOOS=linux
RUN go build -a -installsuffix cgo .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN mkdir /app
WORKDIR /app

COPY --from=builder /go/src/docker/consignment-service .

CMD ["./consignment-service"]