#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=cloudmind.core-api
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}