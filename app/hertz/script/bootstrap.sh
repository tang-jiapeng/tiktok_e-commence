#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=hertz
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}