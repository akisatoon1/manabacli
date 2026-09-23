#!/usr/bin/env bash

# e2eテストを行うためのスクリプト.
# 使い方: 
# 1. 認証ファイル ~/.manabacli/config を用意する.
# 2. `./test.sh url` と, 第1引数にアップロード先のurlを指定する.

if [[ $# -ne 1 ]]; then
    echo "usage: test.sh [url]"
    exit 1
fi

../build.sh
url="$1"
../dist/manabacli "$url" file1.txt file2.txt
