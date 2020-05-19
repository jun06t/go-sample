FROM golang:1.14-alpine

ARG REPO="github.com/jun06t/go-sample/cfs"

RUN apk add --update --no-cache git
RUN mkdir -p /go/src/${REPO}
COPY . /go/src/${REPO}

WORKDIR /go/src/${REPO}

RUN go build -o concprime ./main.go && \
  mv concprime /usr/local/bin/

CMD ["time", "concprime", "-n", "10000"]

