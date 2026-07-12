#!/bin/bash


# 預設目標
.PHONY: help
help:






.PHONY: wire
wire:
	cd internal/container/ && wire

.PHONY: protoc
protoc:
	@protoc \
	-I ./proto \
	--go_out=paths=source_relative:./pb \
	--go-grpc_out=paths=source_relative:./pb \
	$(shell find ./proto -name "*.proto")