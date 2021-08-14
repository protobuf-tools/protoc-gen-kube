# ----------------------------------------------------------------------------
# global

.DEFAULT_GOAL = build

# hack for replace all whitespace to comma
comma := ,
empty :=
space := $(empty) $(empty)

# ----------------------------------------------------------------------------
# Go

PKG_NAME = $(subst $(GO_PATH)/src/,,$(CURDIR))
APP = $(notdir ${PKG_NAME})

GO_PATH  ?= $(shell go env GOPATH)
GO_OS    ?= $(shell go env GOOS)
GO_ARCH  ?= $(shell go env GOARCH)
CGO_ENABLED ?= 0
GOEXPERIMENT=$(shell go env GOEXPERIMENT 2>/dev/null)

PKG := $(subst $(GO_PATH)/src/,,$(CURDIR))
GO_PACKAGE = $(shell go list ./...)
GO_TEST_PACKAGE = $(shell go list -f='{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' ./...)
GO_LINT_PACKAGE ?= ./...

GO_FLAGS := -trimpath

GO_GCFLAGS=

GO_LDFLAGS=-X=${PKG_NAME}/pkg/version.gitCommit=$(shell git describe --abbrev=12 --always)
TAGS=$(shell git name-rev --tags --name-only)
ifneq (${TAGS},)
GO_LDFLAGS+=-X=${PKG_NAME}/pkg/version.version=${TAGS}
endif
GO_LDFLAGS+=-s -w "-extldflags=-static -static-pie"

GO_BUILDTAGS=
ifeq (${CGO_ENABLED},0)
	GO_BUILDTAGS=osusergo netgo
endif
GO_BUILDTAGS_STATIC=static
GO_INSTALLSUFFIX_STATIC=-installsuffix='netgo'
GO_FLAGS +=-tags='$(subst $(space),$(comma),${GO_BUILDTAGS})'

ifneq (${GO_GCFLAGS},)
	GO_FLAGS+=-gcflags='${GO_GCFLAGS}'
endif
ifneq (${GO_LDFLAGS},)
	GO_FLAGS+=-ldflags='${GO_LDFLAGS}'
endif

TOOLS_DIR := ${CURDIR}/tools
TOOLS_BIN := ${TOOLS_DIR}/bin
TOOLS := $(shell cd ${TOOLS_DIR} && go list -v -x -f '{{ join .Imports " " }}' -tags=tools)

GO_TEST ?= ${TOOLS_BIN}/gotestsum --
GO_TEST_FUNC ?= .
GO_TEST_FLAGS ?= -count=1
GO_COVERAGE_OUT ?= coverage.out
GO_LINT_FLAGS ?=

JOBS := $(shell getconf _NPROCESSORS_CONF)
ifeq ($(CIRCLECI),true)
ifeq (${GO_OS},linux)
	# https://circleci.com/changelog#container-cgroup-limits-now-visible-inside-the-docker-executor
	JOBS := $(shell echo $$(($$(cat /sys/fs/cgroup/cpu/cpu.shares) / 1024)))
	GO_TEST_FLAGS+=-p=${JOBS} -cpu=${JOBS}
endif
endif

# ----------------------------------------------------------------------------
# defines

define target
@printf "+ \\x1b[1;32m$(patsubst ,$@,$(1))\\x1b[0m\\n" >&2
endef

# ----------------------------------------------------------------------------
# targets

##@ build

.PHONY: bin/$(APP)
bin/$(APP):
	$(call target,${TARGET})
	@mkdir -p $(@D)
	GOEXPERIMENT=${GOEXPERIMENT} CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build -v $(strip $(GO_FLAGS)) -o $@ ${PKG_NAME}

.PHONY: build
build: GO_BUILDTAGS+=${GO_BUILDTAGS_STATIC}
build: GO_FLAGS+=${GO_INSTALLSUFFIX_STATIC}
build: bin/$(APP)  ## Builds a executable.


##@ test and coverage

export GOTESTSUM_FORMAT=standard-verbose

.PHONY: test
test: tools/bin/gotestsum
test: CGO_ENABLED=1  # needs race test
test: GO_BUILDTAGS+=${GO_BUILDTAGS_STATIC}
test:  ## Runs package test including race condition.
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) $(GO_TEST) -race $(strip $(GO_TEST_FLAGS)) -run=$(GO_TEST_FUNC) $(GO_TEST_PACKAGE)

.PHONY: coverage
ifneq ($(CIRCLECI),true)
coverage: tools/bin/gotestsum
endif
coverage: CGO_ENABLED=1
coverage: GO_BUILDTAGS+=${GO_BUILDTAGS_STATIC}
coverage:  ## Takes packages test coverage.
	$(call target)
	CGO_ENABLED=$(CGO_ENABLED) $(GO_TEST) $(strip $(GO_TEST_FLAGS)) -covermode=atomic -coverpkg=./... -coverprofile=${GO_COVERAGE_OUT} $(GO_PACKAGE)


##@ fmt, vet and lint

.PHONY: fmt
fmt: tools/goimports tools/gofumpt  ## Run goimports and gofumpt.
	$(call target)
	find . -type f -name '*.go' -not -path './vendor/*' | xargs -P ${JOBS} ${TOOLS_BIN}/goimports -local=${PKG} -w
	find . -type f -name '*.go' -not -path './vendor/*' | xargs -P ${JOBS} ${TOOLS_BIN}/gofumpt -s -extra -w

.PHONY: lint
lint: lint/golangci-lint  ## Run all linters.

.PHONY: lint/golangci-lint
lint/golangci-lint: tools/golangci-lint .golangci.yml  ## Run golangci-lint.
	$(call target)
	${TOOLS_BIN}/golangci-lint -j ${JOBS} run $(strip ${GO_LINT_FLAGS}) ${GO_LINT_PACKAGE}


##@ tools

.PHONY: tools
tools: tools/bin/''  ## Install tools

tools/bin/%: ${TOOLS_DIR}/go.mod ${TOOLS_DIR}/go.sum
	@cd tools; \
	  for t in ${TOOLS}; do \
			if [ -z '$*' ] || [ $$(basename $$t) = '$*' ]; then \
				echo "Install $$t ..." >&2; \
				GOBIN=${TOOLS_BIN} CGO_ENABLED=0 go install -mod=mod ${GO_FLAGS} "$${t}"; \
			fi \
	  done


##@ clean

.PHONY: clean
clean:  ## Cleanups extra files in the package.
	$(call target)
	@$(RM) -r *.out *.test *.txt *.prof trace.txt ${TOOLS_BIN}

.PHONY: distclean
distclean:  ## Cleanups binaries and extra files in the package.
	$(call target)
	@$(RM) -r ./bin *.out *.txt *.test *.prof trace.txt ${TOOLS_BIN}


##@ miscellaneous

.PHONY: AUTHORS
AUTHORS:  ## Creates AUTHORS file.
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@printf "$(shell git log --format="\n%aN <%aE>" | LC_ALL=C.UTF-8 sort -uf)" >> $@

.PHONY: todo
todo:  ## Print the all of (TODO|BUG|XXX|FIXME|NOTE) in packages.
	@rg -t go -C 3 -e '(TODO|BUG|XXX|FIXME|NOTE)(\(.+\):|:)' --follow --hidden --glob='!vendor' --glob='!internal'

.PHONY: nolint
nolint:  ## Print the all of //nolint:... pragma in packages.
	@rg -t go -C 3 -e '//nolint.+' --follow --hidden --glob='!vendor' --glob='!internal'

.PHONY: env/% env
env:  ## Print the value of MAKEFILE_VARIABLE. Use `make env/MAKEFILE_VARIABLE`.
env/%:
	@echo $($*)


##@ help

.PHONY: help
help:  ## Show make target help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[33m<target>\033[0m\n"} /^[a-zA-Z_0-9\/_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
