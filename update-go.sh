#!/bin/bash

GOLANG_VERSION="1.20"
GOLANG_OS="linux"
GOLANG_ARCH="arm64"
TINYGO_VERSION="0.28.1"
TINYGO_ARCH="arm64"

# -----------------------
# Install GoLang
# -----------------------

echo "Installing Go & TinyGo"

wget https://go.dev/dl/go${GOLANG_VERSION}.${GOLANG_OS}-${GOLANG_ARCH}.tar.gz

sudo rm -rf /usr/local/go 
sudo tar -C /usr/local -xzf go${GOLANG_VERSION}.${GOLANG_OS}-${GOLANG_ARCH}.tar.gz

echo "" >> ~/.bashrc
echo 'export GOLANG_HOME="/usr/local/go"' >> ~/.bashrc
echo 'export PATH="\$GOLANG_HOME/bin:\$PATH"'  >> ~/.bashrc
source ~/.bashrc

rm go${GOLANG_VERSION}.${GOLANG_OS}-${GOLANG_ARCH}.tar.gz               

# -----------------------
# Install TinyGo
# -----------------------
wget https://github.com/tinygo-org/tinygo/releases/download/v${TINYGO_VERSION}/tinygo_${TINYGO_VERSION}_${TINYGO_ARCH}.deb
sudo dpkg -i tinygo_${TINYGO_VERSION}_${TINYGO_ARCH}.deb
rm tinygo_${TINYGO_VERSION}_${TINYGO_ARCH}.deb

export GOLANG_HOME="/usr/local/go"
export PATH="$GOLANG_HOME/bin:$PATH"

go version
go install -v golang.org/x/tools/gopls@latest
go install -v github.com/ramya-rao-a/go-outline@latest
go install -v github.com/stamblerre/gocode@v1.0.0
