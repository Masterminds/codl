VERSION := $(shell git describe --tags)

build:
	go build -o codl -ldflags "-X main.version ${VERSION}" codl.go

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./codl ${DESTDIR}/usr/local/bin/codl

test: clean
	go test ./parser ./routes ./cmd

clean:
	rm -f ./codl.test
	rm -f ./codl
	rm -f ./routes/app.go

bootstrap:
	cd routes && go test

.PHONY: build test install clean

