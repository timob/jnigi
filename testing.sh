#!/bin/bash

java/test/javac.sh && \
go test -v --tags=cgo_testing
