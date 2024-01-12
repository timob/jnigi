#!/bin/bash

java/test/javac.sh && \
go test --tags=cgo_testing
