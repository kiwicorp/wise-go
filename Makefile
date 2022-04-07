# Makefile

PKG = $(shell cat go.mod | head -n 1 | cut -d ' ' -f 2)

all: test

test:
	go test ./...
.PHONY: test
