FROM golang:1.13 AS builder
WORKDIR /enlabs
COPY . /enlabs
RUN go get -u -v golang.org/x/lint/golint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.21.0
RUN make all

FROM golang:1.13
RUN mkdir -p /opt/enlabs
WORKDIR /opt/enlabs
COPY --from=builder /enlabs/bin/enlabs .
CMD ./enlabs