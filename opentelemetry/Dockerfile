FROM golang:1.20-alpine

ARG TARGET

ADD . /workspace
WORKDIR /workspace
RUN GOOS=linux GOARCH=amd64 go build -o api $TARGET/main.go
RUN mv /workspace/api /usr/local/bin

CMD ["api"]
