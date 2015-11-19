FROM golang:1.5

RUN mkdir -p $GOPATH/src/bitbucket.org/dukex/uhura-api
WORKDIR $GOPATH/src/bitbucket.org/dukex/uhura-api

COPY . $GOPATH/src/bitbucket.org/dukex/uhura-api

RUN go get github.com/tools/godep
RUN godep restore

RUN go get github.com/pilu/fresh

RUN go get

RUN rm -Rf Godeps

EXPOSE 3000

CMD ["/go/bin/uhura-api"]

