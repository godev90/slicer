.PHONY: test cover clean proto

PKG := $(shell go list ./... | grep -v '/pb' | tr '\n' ' ')
COVERPKG := $(shell go list ./... | grep -v '/pb' | paste -sd, -)
COVER_OUT := coverage.out

PROTO_DIR := pb
PB_DIR    := ./

PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

proto:
	@echo ">> Generating gRPC stubs"
	protoc \
	  --go_out=paths=source_relative:$(PB_DIR) \
	  --go-grpc_out=paths=source_relative:$(PB_DIR) \
	  $(PROTO_FILES)


test: clean
	go test -v $(PKG)

cover: clean
	go test $(PKG) -coverprofile=$(COVER_OUT) -covermode=atomic -coverpkg=$(COVERPKG)
	go tool cover -func=$(COVER_OUT)

clean:
	rm -f $(COVER_OUT)