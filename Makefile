.DEFAULT_GOAL := build
.PHONY: build install docker dockerpush

REPO=DODOEX/oracle-adapter
LDFLAGS=-ldflags "-X github.com/DODOEX/oracle-adapter/store.Sha=`git rev-parse HEAD`"

build:
	@go build $(LDFLAGS) -o oracle-adapter

install:
	@go install $(LDFLAGS)

docker:
	@docker build . -t $(REPO)

dockerpush:
	@docker push $(REPO)