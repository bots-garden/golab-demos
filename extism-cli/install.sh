#!/bin/bash
#go install github.com/extism/cli/extism@latest

EXTISM_VERSION="0.3.0"
EXTISM_ARCH="arm64"

wget https://github.com/extism/cli/releases/download/v${EXTISM_VERSION}/extism-v${EXTISM_VERSION}-linux-${EXTISM_ARCH}.tar.gz

sudo tar -xf extism-v${EXTISM_VERSION}-linux-${EXTISM_ARCH}.tar.gz -C /usr/bin
rm extism-v${EXTISM_VERSION}-linux-${EXTISM_ARCH}.tar.gz

extism --version
