# Output File Location
ROOT_DIR=$(shell git rev-parse --show-toplevel)
DIR=$(ROOT_DIR)/build
BINARY=kubetunnel
$(shell mkdir -p ${DIR})

APP_VERSION="0.2.8" # TODO: get this one from env var
OPERATOR_VERSION="0.0.15"
KUBETUNNEL_VERSION="1.1.4"
# Go build flags
LDFLAGS=-ldflags "-X main.Version=${APP_VERSION} -X main.OperatorVersion=${OPERATOR_VERSION}"

default-cli:
	go build ${LDFLAGS} -o ${DIR}/${BINARY} ${ROOT_DIR}/cmd/cli

# Compile CLI - Windows x64
windows-cli:
	mkdir -p ${DIR}/windows
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${DIR}/windows/${BINARY}.exe ${ROOT_DIR}/cmd/cli
	zip ${DIR}/${BINARY}-win.zip ${DIR}/windows/*
# Compile CLI - Linux x64
linux-cli:
	mkdir -p ${DIR}/linux
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY} ${ROOT_DIR}/cmd/cli
	zip ${DIR}/${BINARY}-linux.zip ${DIR}/linux/*

# Compile CLI - Darwin x64
darwin-cli:
	mkdir -p ${DIR}/mac
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY} ${ROOT_DIR}/cmd/cli
	zip ${DIR}/${BINARY}-mac.zip ${DIR}/mac/*

all-cli: darwin-cli linux-cli windows-cli

build_kubetunnel_server:
	docker build -t dcodetech/kubetunnel:${KUBETUNNEL_VERSION} . && docker push dcodetech/kubetunnel:${KUBETUNNEL_VERSION}
	#sed -i 's/KUBETUNNEL_SERVER_VERSION/${KUBETUNNEL_VERSION}/g' ${ROOT_DIR}/pkg/operator/helm-charts/templates

build_operator:
	docker build -t dcodetech/kubetunnel-operator:${OPERATOR_VERSION} -f Dockerfile.operator . && docker push dcodetech/kubetunnel-operator:${OPERATOR_VERSION}
#	sed -i 's/KUBETUNNEL_OPERATOR_VERSION/${OPERATOR_VERSION}/g' ${ROOT_DIR}/charts/kubetunnel-operator/Chart.yaml
#	sed -i 's/KUBETUNNEL_OPERATOR_VERSION/${OPERATOR_VERSION}/g' ${ROOT_DIR}/charts/kubetunnel-operator/templates/operator.yaml




clean:
	rm -rf ${DIR}*
