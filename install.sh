#!/bin/bash

jdk_path=$1

if [[ "$jdk_path" == "" ]]; then
	echo "jdk path not given" 
	echo "usage: ./install.sh <jdk path>"
	exit 1
fi

if [[ "$OSTYPE" == "linux-gnu" ]]; then
	osdir=linux
elif [[ "$OSTYPE" == "darwin"* ]]; then
	osdir=darwin
elif [[ "$OSTYPE" == "msys" ]]; then
	osdir=win32
else
	echo "ERROR: unkown OSTYPE: $OSTYPE"
	exit 1
fi


export CGO_CFLAGS="-I${jdk_path}/include -I${jdk_path}/include/${osdir}"

go install
