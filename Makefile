# Copyright (C) 2020 The go-mysql Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL := bash

PREFIX?=$(shell pwd)

GIT_ROOT=
PRODUCT_NAME=go-mysql
PACKAGE_NAME=mysql

MODULE_ROOT=${PACKAGE_NAME}
MODULE_PACKAGE_ROOT=${GIT_ROOT}${PRODUCT_NAME}/${MODULE_ROOT}
MODULE_PACKAGES=\
	${MODULE_PACKAGE_ROOT} \
	${MODULE_PACKAGE_ROOT}/log \
	${MODULE_PACKAGE_ROOT}/query

EXAMPLES_ROOT=examples
EXAMPLES_PACKAGE_ROOT=${GIT_ROOT}${PRODUCT_NAME}/${EXAMPLES_ROOT}
EXAMPLES_DEAMON_BIN=go-mysqld
EXAMPLES_DEAMON_ROOT=${EXAMPLES_PACKAGE_ROOT}/${EXAMPLES_DEAMON_BIN}
EXAMPLES_PACKAGES=\
	${EXAMPLES_DEAMON_ROOT}/server \
	${EXAMPLES_DEAMON_ROOT}/server/storage
EXAMPLE_BINARIES=\
	${EXAMPLES_DEAMON_ROOT}

TEST_ROOT=test
TEST_PACKAGES=\
	${TEST_ROOT} \
	${TEST_ROOT}/client 

ALL_ROOTS=\
	${MODULE_ROOT} \
	${EXAMPLES_ROOT} \
	${TEST_ROOT}

ALL_PACKAGES=\
	${MODULE_PACKAGES} \
	${EXAMPLES_PACKAGES} \
	${TEST_PACKAGES}

BINARIES=${EXAMPLE_BINARIES}

.PHONY: clean check_style

all: test

format:
	gofmt -w ${ALL_ROOTS}

vet: format
	go vet ${ALL_PACKAGES}

build: vet
	go build -v ${MODULE_PACKAGES}

test: vet
	go test -v -cover -p=1 ${ALL_PACKAGES}

install: build
	go install -v -gcflags=${GCFLAGS} ${BINARIES}

clean:
	go clean -i ${ALL_PACKAGES}

check_style:
	(! find . -name '*.go' | xargs gofmt -s -d | grep '^')
	(! find . -name '*.go' | xargs goimports -d | grep '^')
	golangci-lint run --timeout=5m
