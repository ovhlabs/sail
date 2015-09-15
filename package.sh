#!/bin/bash

for GOOS in darwin linux windows; do
    for GOARCH in 386 amd64; do
        architecture="${GOOS}-${GOARCH}"
        echo "Building ${path}"
        export GOOS=$GOOS
        export GOARCH=$GOARCH
        go build -ldflags "-X ${PROJECT_PATH}/${PROJECT_NAME}/update.architecture=${architecture} -X ${PROJECT_PATH}/${PROJECT_NAME}/update.urlUpdateRelease=${URL_UPDATE_RELEASE} -X ${PROJECT_PATH}/${PROJECT_NAME}/update.urlUpdateSnapshot=${URL_UPDATE_SNAPSHOT}" -o bin/${architecture}/sail
    done
done
