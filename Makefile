.PHONY: all test clean
APPNAME := yaml-renderer
APPSRC := .

GITCOMMITHASH := $(shell git log --max-count=1 --pretty="format:%h" HEAD)
GITCOMMIT := -X main.gitcommit=$(GITCOMMITHASH)

VERSIONTAG := $(shell git describe --tags --abbrev=0)
VERSION := -X main.appversion=$(VERSIONTAG)

BUILDTIMEVALUE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
BUILDTIME := -X main.buildtime=$(BUILDTIMEVALUE)

LDFLAGS := '-extldflags "-static" -d -s -w $(GITCOMMIT) $(VERSION) $(BUILDTIME)'

all:info clean build

clean:
	rm -rf build rendered

info: 
	@echo - appname:   $(APPNAME)
	@echo - version:   $(VERSIONTAG)
	@echo - commit:    $(GITCOMMITHASH)
	@echo - buildtime: $(BUILDTIMEVALUE) 

dep:
	@go get -v -d ./...

build-linux: info dep
	@echo Building for linux
	@mkdir -p build/linux
	CGO_ENABLED=0 \
	GOOS=linux \
	go build -o build/linux/$(APPNAME)-$(VERSIONTAG)-$(GITCOMMITHASH) -a -ldflags $(LDFLAGS) $(APPSRC)
	@cp build/linux/$(APPNAME)-$(VERSIONTAG)-$(GITCOMMITHASH) build/linux/$(APPNAME)


test:
	TESTVAR=testvarbla TESTLIST="testlist1,testlist2" TESTVAR2="var2" go run main.go -t test/templates -y test/values.yaml
	grep "test1,test2,test3" rendered/file.txt
	grep "foo-obj-name" rendered/file.txt
	grep "testvarbla" rendered/file.txt
	grep "foo-testvarbla-var2" rendered/file.txt
	grep "Just a var: testvarbla" rendered/file.txt
	grep "Not a list: testlist1,testlist2" rendered/file.txt
	grep 1234 rendered/file.txt
	grep '$$NOTSETVAR' rendered/file.txt
	grep '$${NOTSETVAR}' rendered/file.txt