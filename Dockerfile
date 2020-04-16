FROM golang:1.13 AS builder
WORKDIR /go/src/github.com/polisgo2020/senyast4745/
COPY . .
RUN go get -v && CGO_ENABLED=0 go build -o app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
ENV IND_FILE final.csv
RUN mkdir /output
COPY --from=0 /go/src/github.com/polisgo2020/senyast4745/app .
CMD ./app search -index /output/$IND_FILE