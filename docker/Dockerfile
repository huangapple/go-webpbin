FROM golang

WORKDIR /

RUN go get -v github.com/nickalie/go-binwrapper && \
    go get -v github.com/stretchr/testify/assert && \
    go get -v golang.org/x/image/webp

RUN mkdir -p $GOPATH/src/github.com/huangapple/go-webpbin
COPY . $GOPATH/src/github.com/huangapple/go-webpbin
WORKDIR $GOPATH/src/github.com/huangapple/go-webpbin
RUN go test -v ./...