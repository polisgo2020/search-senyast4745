FROM golang:1.13 AS builder
WORKDIR /go/src/github.com/polisgo2020/senyast4745/
COPY . .
RUN go get -v && go build -o app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=0 /go/src/github.com/polisgo2020/senyast4745/app .
RUN mkdir /output
ENV IND_FILE final.csv
CMD ./app search -index /output/$IND_FILE