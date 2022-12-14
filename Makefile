#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --tags)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
LEDGER_ENABLED ?= true
DOCKER := $(shell which docker)
BUILDDIR ?= $(CURDIR)/build
export GO111MODULE = on

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

# process linker flags
ldflags = -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"
ldflags += -w -s

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) .

build:
	go build $(BUILD_FLAGS) -o bin/defiant-swap .

check-swagger:
	which swagger || (GO111MODULE=on go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	GO111MODULE=on ${HOME}/go/bin/swagger generate spec -o ./swagger.yaml --scan-models

serve-swagger: check-swagger
	${HOME}/go/bin/swagger serve -F=swagger swagger.yaml