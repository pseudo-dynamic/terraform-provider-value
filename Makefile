TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=github.com
NAMESPACE=pseudo-dynamic
NAME=value
VERSION=0.1.0
BINARY=terraform-provider-${NAME}
BINARY_VERSION_CORE=${BINARY}_v${VERSION}

ifdef APPDATA
	OS_ARCH=windows_amd64
	PLUGIN_CACHE=${APPDATA}\terraform.d\plugins
	BINARY_VERSION=${BINARY_VERSION_CORE}.exe
else
	OS_ARCH=darwin_amd64
	PLUGIN_CACHE=~/.terraform.d/plugins
	BINARY_VERSION=${BINARY_VERSION_CORE}
endif

default: install

build:
	go build -o ${BINARY_VERSION}

generate:
	go generate

# release: generate
# 	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
# 	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
# 	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
# 	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
# 	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
# 	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
# 	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
# 	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
# 	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
# 	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
# 	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
# 	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build
	mkdir -p ${PLUGIN_CACHE}/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY_VERSION} ${PLUGIN_CACHE}/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/

# test:
# 	go test -i $(TEST) || exit 1
# 	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4
