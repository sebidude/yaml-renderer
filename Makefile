APPNAME := yaml-renderer
APPSRC := ./cmd/$(APPNAME)

GITCOMMITHASH := $(shell git log --max-count=1 --pretty="format:%h" HEAD)
GITCOMMIT := -X main.gitcommit=$(GITCOMMITHASH)

VERSIONTAG := $(shell git describe --tags --abbrev=0)
VERSION := -X main.appversion=$(VERSIONTAG)

BUILDTIMEVALUE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
BUILDTIME := -X main.buildtime=$(BUILDTIMEVALUE)

LDFLAGS := '-extldflags "-static" -d -s -w $(GITCOMMIT) $(VERSION) $(BUILDTIME)'

all:info clean build

clean:
	rm -rf build

info: 
	@echo - appname:   $(APPNAME)
	@echo - verison:   $(VERSIONTAG)
	@echo - commit:    $(GITCOMMITHASH)
	@echo - buildtime: $(BUILDTIMEVALUE) 

dep:
	@go get -v -d ./...

build-linux: info dep
	@echo Building for linux
	@mkdir -p build/linux
	@CGO_ENABLED=0 \
	GOOS=linux \
	go build -o build/linux/$(APPNAME)-$(VERSIONTAG)-$(GITCOMMITHASH) -a -ldflags $(LDFLAGS) $(APPSRC)
	@cp build/linux/$(APPNAME)-$(VERSIONTAG)-$(GITCOMMITHASH) build/linux/$(APPNAME)


image: 
	@echo Creating docker image
	@docker build -t $(APPNAME):$(VERSIONTAG)-$(GITCOMMITHASH) .

test: build-linux
	TESTVAR=testvarbla ./build/linux/$(APPNAME) -o manifests -s txt -t test/templates -y test/values.yaml
	grep "test1,test2,test3" manifests/file.txt
	grep "foo-obj-name" manifests/file.txt
	grep "testvarbla" manifests/file.txt
	grep 1234 manifests/file.txt
	grep '$$NOTSETVAR' manifests/file.txt


