.PHONY: test cover clean

PKG := $(shell go list ./... | grep -v '/pb' | tr '\n' ' ')
COVERPKG := $(shell go list ./... | grep -v '/pb' | paste -sd, -)
COVER_OUT := coverage.out

test: clean
	go test -v $(PKG)

cover: clean
	go test $(PKG) -coverprofile=$(COVER_OUT) -covermode=atomic -coverpkg=$(COVERPKG)
	go tool cover -func=$(COVER_OUT)

clean:
	rm -f $(COVER_OUT)