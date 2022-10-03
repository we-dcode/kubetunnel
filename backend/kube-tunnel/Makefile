# Output File Location
DIR=build
BINARY_NAME=kubetunnel
$(shell mkdir -p ${DIR})

# Go build flags
LDFLAGS=-ldflags '-X main.build=${BUILD} -buildid='

default:
	go build ${LDFLAGS} -o ${DIR} cmd/kubetunnel.go

# Compile Server - Windows x64
windows:
	export GOOS=windows;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/${BINARY_NAME}-Windows-x64.exe cmd/kubetunnel.go

# The SEED must be the exact same that was used when compiling the agent
# Garble version 0.5.2 or later must be installed and accessible in the PATH environment variable
windows-garble:
	export GOGARBLE=${GOGARBLE};export GOOS=windows GOARCH=amd64;garble -tiny -literals -seed ${SEED} build ${LDFLAGS} -o ${DIR}/${BINARY_NAME}-Windows-x64.exe cmd/kubetunnel.go

# Compile Server - Linux x64
linux:
	export GOOS=linux;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/kubetunnel-Linux-x64 main.go

# The SEED must be the exact same that was used when compiling the agent
# Garble version 0.5.2 or later must be installed and accessible in the PATH environment variable
linux-garble:
	export GOGARBLE=${GOGARBLE};export GOOS=linux GOARCH=amd64;garble -tiny -literals -seed ${SEED} build ${LDFLAGS} -o ${DIR}/${BINARY_NAME}-Linux-x64 cmd/kubetunnel.go

# Compile Server - Darwin x64
darwin:
	export GOOS=darwin;export GOARCH=amd64;go build ${LDFLAGS} -o ${DIR}/merlinServer-Darwin-x64 cmd/kubetunnel

# The SEED must be the exact same that was used when compiling the agent
# Garble version 0.5.2 or later must be installed and accessible in the PATH environment variable
darwin-garble:
	export GOGARBLE=${GOGARBLE};export GOOS=darwin GOARCH=amd64;garble -tiny -literals -seed ${SEED} build ${LDFLAGS} -o ${DIR}/${BINARY_NAME}-Darwin-x64.exe cmd/kubetunnel.go

clean:
	rm -rf ${DIR}*