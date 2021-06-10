FROM golang:alpine

RUN apk update \
    && apk add lvm2

WORKDIR /go/src/lvm-exporter
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o ./lvm-exporter /go/src/lvm-exporter/cmd/main.go

CMD ["./lvm-exporter"]
