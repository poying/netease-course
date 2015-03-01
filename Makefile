PROJNAME = necourse
SERVICEPATH = github.com/poying

SOURCES = $(wildcard *.go)
GOPATH = $(shell pwd)
PROJPATH = src/$(SERVICEPATH)/$(PROJNAME)
BINFILE = ./$(PROJNAME)

define build
	echo build $(PROJNAME)-$(1)-$(2); \
	GO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build -o "bin/$(PROJNAME)-$(1)-$(2)" $(SOURCES);
endef

RUN_ARGS=$(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

all: test clean

test: init $(BINFILE)
	@mkdir -p test-data
	@$(BINFILE) test-data >/dev/null 2>&1 &
	@go test -v ./redisfs/...
	@$(BINFILE) --unmount test-data
	@rm -r test-data

run:
	@go run $(SOURCES) $(RUN_ARGS)

cross:
	@$(call build,linux,amd64)
	@$(call build,linux,386)
	@$(call build,linux,arm)
	@$(call build,darwin,amd64)

$(BINFILE):
	@go build -o redis-mount

install:
	@go install

init: get-deps $(PROJPATH)

get-deps:
	@go get github.com/tj/go-spin
	@go get github.com/codegangsta/cli

$(PROJPATH):
	@mkdir -p $(dir $@)
	@ln -s $(GOPATH) $@

clean:
	-@rm -rf bin src pkg $(BINFILE)

.PHONY: $(BINFILE)
