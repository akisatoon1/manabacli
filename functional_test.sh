#!/usr/bin/env bash

source env.sh
go test -v -tags=functional ./src/logic
