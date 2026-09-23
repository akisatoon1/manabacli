#!/usr/bin/env bash

# 他のディレクトリで実行されてもビルドできるようにするため.
cd "$(dirname "$0")"

go build -o dist/manabacli ./src/
