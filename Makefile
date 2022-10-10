# Output File Location
ROOT_DIR=$(shell git rev-parse --show-toplevel)
DIR=$(ROOT_DIR)/build
SLUG=kubetunnel
$(shell mkdir -p ${DIR})

# Go build flags
LDFLAGS=-ldflags '-X main.build=${BUILD} -buildid='

default-cli:
	go build ${LDFLAGS} -o ${DIR} ${ROOT_DIR}/cmd/cli

# Compile CLI - Windows x64
windows-cli:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${DIR}/${SLUG}-Windows-x64.exe ${ROOT_DIR}/cmd/cli

# Compile CLI - Linux x64
linux-cli:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${DIR}/${SLUG}-Linux-x64 ${ROOT_DIR}/cmd/cli

# Compile CLI - Darwin x64
darwin-cli:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${DIR}/${SLUG}-Darwin-x64 ${ROOT_DIR}/cmd/cli

all-cli: darwin-cli linux-cli windows-cli

clean:
	rm -rf ${DIR}*