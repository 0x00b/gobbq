.PHONY: all build  

include lint.mk

WORK_PATH=$(shell cd .. && pwd)

all:
	$(MAKE) lintAll
	$(MAKE) build 

build:
	go build ./...

test:
	go test ./...

clean:
	go clean ./...
	@echo clean finished
