ASSETS_DIR=public/assets

build: deps_save test

deploy: build
	git push heroku master

deps:
	go get github.com/pilu/fresh
	go get

deps_save:
	godep save

test:
	go test

coverage:
	go test -coverprofile=coverage.out ./core
	go tool cover -html=coverage.out
	rm coverage.out
