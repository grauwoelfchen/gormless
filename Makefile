GIT_REF := $(shell git describe --always)
GIT_TAG := $(shell git describe --contains "$(GIT_REF)" 2>/dev/null | \
  sed "s/.*~.*//")
VERSION ?= $(if $(GIT_TAG),$(GIT_TAG),$(GIT_REF))

DEV_TOOLS = \
  golang.org/x/lint/golint \
  github.com/client9/misspell/cmd/misspell \
  github.com/kisielk/errcheck \
  honnef.co/go/tools/cmd/staticcheck

GOTEST := $(shell if type gotest >/dev/null 2>&1; then echo "gotest -v"; \
  else echo "go test -v"; fi)

MIGRATION_DIR ?= migration
DEFAULT_TAG ?= sqlite

TARGETS = $$(go list ./... | grep -v $(MIGRATION_DIR))

.DEFAULT_GOAL = help

# setup {{{
setup\:tools:
	@echo "Installing/Updating..."
	@for tool in $(DEV_TOOLS); do \
	  echo "  $$tool"; \
	  GO111MODULE=off go get -u $$tool; \
	done
.PHONY: setup\:tools
# }}}

# verify {{{
verify\:fmt:  ## Display file names need to be fixed [alias: fmt]
	@output=`gofmt -l . 2>&1`; \
	if [ "$$output" ]; then \
	echo "Run \`gofmt\` on the following files:"; \
	echo "$$output"; \
	exit 1; \
	fi

fmt: verify\:fmt
.PHONY: verify\:fmt fmt

verify\:format:  ## Print fix suggestions as diff outputs [alias: format]
	@gofmt -d -e `find . -type f -name *.go`

format: verify\:format
.PHONY: verify\:format format

verify\:vet:  ## Print filenames reported by go-vet [alias: vet]
	@GO111MODULE=on go vet -tags="$(DEFAULT_TAG)" $(TARGETS)

vet: verify\:vet
.PHONY: verify\:vet vet

verify\:lint:  ## Verify style using goling [alias: lint]
	@golint -set_exit_status $(TARGETS)

lint: verify\:lint
.PHONY: verify\:lint lint

verify\:error:  ## Verify errors using errcheck [alias: errcheck]
	@GO111MODULE=on errcheck -tags="$(DEFAULT_TAG)" \
		-ignoretests $(TARGETS)

errcheck: verify\:error
.PHONY: verify\:error errcheck

verify\:spell:  ## Verify spell using misspell [alias: misspell]
	@misspell $$(git ls-files)

misspell: verify\:spell
.PHONY: verify\:spell misspell

verify\:code:  ## Verify code using staticcheck [alias: staticcheck]
	@GO111MODULE=on staticcheck -tags="$(DEFAULT_TAG)" $(TARGETS)

staticcheck: verify\:code
.PHONY: verify\:code staticcheck

verify\:all: fmt vet lint errcheck misspell staticcheck   ## Run all verify:xxx checks
.PHONY: verify\:all
# }}}

# test {{{
test:  ## Run unit tests
	@GO111MODULE=on DATABASE_URL=":memory:" $(GOTEST) -tags="$(DEFAULT_TAG)" \
		$(TARGETS)
.PHONY: test
# }}}

# build {{{
__build-%:
	tag="$(subst __build-,,$@)" && \
	GO111MODULE=on go build -tags="$${tag}" -o bin/gormless \
	  -ldflags "-X main.version=$(VERSION)" \
	  gitlab.com/grauwoelfchen/gormless/cmd/gormless

build\:mssql: __build-mssql  ## Build cli application for mssql
.PHONY: build\:mssql

build\:mysql: __build-mysql  ## Build cli application for mysql
.PHONY: build\:mysql

build\:sqlite: __build-sqlite  ## Build cli application for sqlite
.PHONY: build\:sqlite

build\:postgres: __build-postgres  ## Build cli application for postgres
.PHONY: build\:postgres
# }}}

# other utilities {{{
migrations:  ## Build migrations for development
	@for dir in $$(ls migration); do \
	  cd migration/$$dir; \
	  GO111MODULE=on go build -buildmode=plugin; \
	done
.PHONY: migrations

help:  ## Display this message
	@grep -E '^[a-z\:\\\-]+: ' $(firstword $(MAKEFILE_LIST)) | \
	  grep -E '  ## ' | \
	  sed -e 's/\( \| [ :_0-9a-z\-]*\)  /  /g' | \
	  tr -d \\\\ | \
	  awk 'BEGIN {FS = ":  ## "}; {printf "%-16s%s\n", $$1, $$2};'
.PHONY: help

version:  ## Print version (with commit counts)
	@echo "$(VERSION) ($$(git rev-list HEAD --count) commits)"
.PHONY: version
# }}}
