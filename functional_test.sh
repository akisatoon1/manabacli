#!/usr/bin/env bash

source env.sh
go test -v -count=1 -tags=functional ./src/logic
