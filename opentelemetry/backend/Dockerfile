FROM golang:1.20-alpine

ADD . /workspace
WORKDIR /workspace
RUN GOOS=linux GOARCH=amd64 go build -o api .
RUN mv /workspace/api /usr/local/bin

CMD ["api"]
