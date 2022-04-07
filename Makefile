# Makefile

PKG = $(shell cat go.mod | head -n 1 | cut -d ' ' -f 2)
CMD ?= $(shell find cmd -type d ! -name cmd | xargs -I {} printf "%s " $(PKG)/{})

all: test install

test:
	go test ./...
.PHONY: test

install:
	go install $(CMD)
.PHONY: install
