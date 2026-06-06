# Copyright (c) 2026 Contentways
# SPDX-License-Identifier: MIT

.PHONY: build clean test coverage

build:
	go build -o build/poweradmin .

clean:
	rm -rf build/

test:
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.txt
	go tool cover -func=coverage.txt
