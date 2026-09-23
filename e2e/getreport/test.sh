#!/usr/bin/env bash

if [[ $# -ne 1 ]]; then
    echo "usage: test.sh [url]"
    exit 1
fi

../../build.sh
url="$1"
../../dist/manabacli getreport "$url"
