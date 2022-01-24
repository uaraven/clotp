#!/bin/sh

OSES="darwin linux windows"
ARCHS="amd64 arm64"

rm -rf builds > /dev/null
mkdir builds

go test ./...

if [ $? -ne 0 ]; then
    echo "Tests failed"
    exit
fi

for OS in $OSES 
do
    for ARCH in $ARCHS 
    do
        echo "Building for ${OS}_$ARCH"
        OUT_NAME="builds/clotp_${OS}_${ARCH}"
        if [ "$OS" == "windows" ]; then
            OUT_NAME="${OUT_NAME}.exe"
        fi
        GOOS="$OS" GOARCH="$ARCH" go build -ldflags "-s -w" -o $OUT_NAME
    done
done

