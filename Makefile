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

GOBIN := $(shell go env GOPATH)/bin
PATH := $(GOBIN):$(PATH)
LDFLAGS=-checklinkname=0

GIT_ROOT=github.com/cybergarage
PRODUCT_NAME=go-mysql
MODULE_ROOT=${GIT_ROOT}/${PRODUCT_NAME}

PKG_NAME=mysql
PKG_VER=$(shell git describe --abbrev=0 --tags)
PKG_COVER=${PKG_NAME}-cover
PKG_SRC_ROOT=${PKG_NAME}
PKG=${MODULE_ROOT}/${PKG_SRC_ROOT}

TEST_SRC_ROOT=${PKG_NAME}test
TEST_PKG=${MODULE_ROOT}/${TEST_SRC_ROOT}

EXAMPLES_ROOT=examples
EXAMPLES_SRC_ROOT=${EXAMPLES_ROOT}
EXAMPLES_DEAMON_BIN=go-mysqld
EXAMPLES_DOCKER_TAG=cybergarage/${EXAMPLES_DEAMON_BIN}:${PKG_VER}
EXAMPLES_PKG_ROOT=${GIT_ROOT}/${PRODUCT_NAME}/${EXAMPLES_ROOT}
EXAMPLES_SRC_ROOT=${EXAMPLES_ROOT}
EXAMPLES_DEAMON_ROOT=${EXAMPLES_PKG_ROOT}/
EXAMPLE_BINARIES=\
	${EXAMPLES_DEAMON_ROOT}/${EXAMPLES_DEAMON_BIN}

BINARIES=${EXAMPLE_BINARIES}

.PHONY: version clean test
.IGNORE: lint

all: test

version:
    @pushd ${PKG_SRC_ROOT} && ./version.gen > version.go && popd
    -git commit ${PKG_SRC_ROOT}/version.go -m "Update version"

format: version
	gofmt -s -w ${PKG_SRC_ROOT} ${TEST_SRC_ROOT} ${EXAMPLES_SRC_ROOT}

vet: format
	go vet ${PKG}

lint: vet
	golangci-lint run ${PKG_SRC_ROOT}/... ${TEST_SRC_ROOT}/... ${EXAMPLES_SRC_ROOT}/...

test: lint
	go test -v -p 1 -timeout 10m -ldflags=${LDFLAGS} -cover -coverpkg=${PKG}/... -coverprofile=${PKG_COVER}.out ${PKG}/... ${TEST_PKG}/...
	go tool cover -html=${PKG_COVER}.out -o ${PKG_COVER}.html

build: test
	go build -v -gcflags=${GCFLAGS} -ldflags=${LDFLAGS} ${BINARIES}

install: build
	go install -v -gcflags=${GCFLAGS} -ldflags=${LDFLAGS} ${BINARIES}

run: install
	$(GOBIN)/${EXAMPLES_DEAMON_BIN}

image: test
	docker image build -t ${EXAMPLES_DOCKER_TAG} .

rund: image
	docker container run -it --rm -p 3307:3307 ${EXAMPLES_DOCKER_TAG}

clean:
	go clean -i ${PKG}
