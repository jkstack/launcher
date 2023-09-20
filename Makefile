.PHONY: all distclean prepare

PROJ=launcher
VERSION=1.0.0
TIMESTAMP=`date +%s`

MAJOR=`echo $(VERSION)|cut -d'.' -f1`
MINOR=`echo $(VERSION)|cut -d'.' -f2`
PATCH=`echo $(VERSION)|cut -d'.' -f3`

BRANCH=`git rev-parse --abbrev-ref HEAD`
HASH=`git log -n1 --pretty=format:%h`
REVERSION=`git log --oneline|wc -l|tr -d ' '`
BUILD_TIME=`date +'%Y-%m-%d %H:%M:%S'`
LDFLAGS="-X 'main.gitBranch=$(BRANCH)' \
-X 'main.gitHash=$(HASH)' \
-X 'main.gitReversion=$(REVERSION)' \
-X 'main.buildTime=$(BUILD_TIME)' \
-X 'main.version=$(VERSION)'"

all: distclean windows.amd64

version:
	@echo $(VERSION)

windows.amd64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags $(LDFLAGS) \
		-o launcher.exe
	makensis -DPRODUCT_VERSION=$(VERSION) \
		-INPUTCHARSET UTF8 deploy/build.nsi

distclean:
	rm -f launcher.exe deploy/launcher.exe