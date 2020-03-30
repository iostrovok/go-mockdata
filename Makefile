CURDIR := $(shell pwd)
GOBIN := $(CURDIR)/bin/
ENV:=GOBIN=$(GOBIN)
DIR:=FILE_DIR=$(CURDIR)/testfiles TEST_SOURCE_PATH=$(CURDIR)
GODEBUG:=GODEBUG=gocacheverify=1
LOCDIR:=$(PWD)

##
## List of commands:
##

## default:
all: mod deps fmt lint test

all-deps: mod deps

deps:
	@echo "======================================================================"
	@echo 'MAKE: deps...'
	@mkdir -p $(GOBIN)
	@$(ENV) go get -u golang.org/x/lint/golint
	@$(ENV) GO111MODULE=on go get github.com/golang/mock/mockgen@latest
	@$(ENV) GO111MODULE=on go get github.com/dizzyfool/genna@latest
	@$(ENV) GO111MODULE=on go get github.com/go-pg/migrations/v7@latest

test: tests-imports tests-pkparser tests-onefunction tests-receivers tests-mockdata

tests-pkparser:
	@echo "Run 'tests-pkparser' race test for ./pkparser/..."
	cd $(LOCDIR)/pkparser/ && $(DIR) $(GODEBUG) go test -cover -race ./...

tests-onefunction:
	@echo "Run 'tests-onefunction' race test for ./onefunction/..."
	cd $(LOCDIR)/onefunction/ && $(DIR) $(GODEBUG) go test -cover -race ./...

tests-mockdata:
	@echo "Run 'tests-mockdata' race test for ./..."
	cd $(LOCDIR)/ && $(DIR) $(GODEBUG) go test -cover -race ./...

tests-receivers:
	@echo "Run 'tests-receivers' race test for ./receivers/..."
	cd $(LOCDIR)/receivers/ && $(DIR) $(GODEBUG) go test -cover -race ./...

tests-imports:
	@echo "Run 'tests-imports' race test for ./imports/..."
	cd $(LOCDIR)/imports/ && $(DIR) $(GODEBUG) go test -cover -race ./...

tests-tool:
	@echo "Run 'tests-tool' race test for ./tool/..."
	cd $(LOCDIR)/tool/ && $(DIR) $(GODEBUG) go test -cover -race ./...

mod:
	@echo "======================================================================"
	@echo "Run MOD"
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod verify
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod tidy
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod vendor
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod download
	GO111MODULE=on GONOSUMDB="*" GOPROXY=direct go mod verify


clean-cache:
	@echo "clean-cache started..."
	go clean -cache
	go clean -testcache
	@echo "clean-cache complete!"
	@echo "clean-cache complete!"


mock-gen:
	GO111MODULE=on ./bin/mockgen -package mmock  -destination ./test-code/mmock/ione_mock.go -source ./test-code/main.go IOne

