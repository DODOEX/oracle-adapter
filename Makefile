.DEFAULT_GOAL := build
.PHONY: build install docker dockerpush

REPO=Dominator008/asset-price-oracle-adapter
LDFLAGS=-ldflags "-X github.com/Dominator008/asset-price-oracle-adapter/store.Sha=`git rev-parse HEAD`"

build:
	@go build $(LDFLAGS) -o asset-price-oracle-adapter

install:
	@go install $(LDFLAGS)

docker:
	@docker build . -t $(REPO)

dockerpush:
	@docker push $(REPO)