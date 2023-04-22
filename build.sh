#!/bin/bash

case "$1" in
    "raspi" ) echo "Raspberry Pi向けにGOOS=linux, GOARCH=arm64でビルドします" 1>&2
        GOOS=linux GOARCH=arm64 go build;;
    * ) echo "デフォルトのOSとアーキテクチャ向けにビルドします" 1>&2
    go build;;
esac

exit 0
