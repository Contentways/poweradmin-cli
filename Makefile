# Copyright (c) 2026 Contentways
# SPDX-License-Identifier: MIT

.PHONY: build clean

build:
	go build -o build/poweradmin .

clean:
	rm -rf build/
