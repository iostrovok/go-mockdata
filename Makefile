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

test: tests-onefunction tests-mockdata

tests-onefunction:
	@echo "Run 'tests-onefunction' race test for ./onefunction/..."
	cd $(LOCDIR)/onefunction/ && $(DIR) $(GODEBUG) go test -cover -race ./...

tests-mockdata:
	@echo "Run 'tests-mockdata' race test for ./..."
	cd $(LOCDIR)/ && $(DIR) $(GODEBUG) go test -cover -race ./...


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
	GO111MODULE=on ./bin/mockgen -package mmock github.com/iostrovok/go-mockdata/test-code IOne > ./test-code/mmock/ione_mock.go

