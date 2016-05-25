uhura-api
=========


### Install

With [golang](https://golang.org/) installed and configured, run:


    $ mkdir -p $GOPATH/src/github.com/uhuraapp
    $ cd $GOPATH/src/github.com/uhuraapp
    $ go get
    $ go get github.com/pilu/fresh
    $ // setup database
    $ createdb uhura -h 127.0.0.1 -U postgres -p 5432

### Up Server
    $ export DATABASE_URL=postgres://postgres@127.0.0.1:5432/uhura?sslmode=disable
    $ export REDIS_URL=redis://127.0.0.1:6379
    $ export PORT=3000
    $ fresh  // API is on http://localhost:3000/api

