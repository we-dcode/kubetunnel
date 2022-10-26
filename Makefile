# Output File Location
ROOT_DIR=$(shell git rev-parse --show-toplevel)
DIR=$(ROOT_DIR)/build
BINARY=kubetunnel
$(shell mkdir -p ${DIR})

APP_VERSION="0.2.7" # TODO: get this one from env var
OPERATOR_VERSION="0.0.6"
# Go build flags
LDFLAGS=-ldflags "-X main.Version=${APP_VERSION} -X main.OperatorVersion=${OPERATOR_VERSION}"

default-cli:
	go build ${LDFLAGS} -o ${DIR}/${BINARY} ${ROOT_DIR}/cmd/cli

# Compile CLI - Windows x64
windows-cli:
	mkdir -p ${DIR}/windows
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${DIR}/windows/${BINARY}.exe ${ROOT_DIR}/cmd/cli
	#$(shell zip ${DIR}/${BINARY}-win.zip ${DIR}/windows/*)
# Compile CLI - Linux x64
linux-cli:
	mkdir -p ${DIR}/linux
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${DIR}/linux/${BINARY} ${ROOT_DIR}/cmd/cli
	#$(shell zip ${DIR}/${BINARY}-linux.zip ${DIR}/linux/*)

# Compile CLI - Darwin x64
darwin-cli:
	mkdir -p ${DIR}/mac
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${DIR}/mac/${BINARY} ${ROOT_DIR}/cmd/cli
	zip ${DIR}/${BINARY}-mac.zip ${DIR}/mac/*

all-cli: darwin-cli linux-cli windows-cli

clean:
	rm -rf ${DIR}*