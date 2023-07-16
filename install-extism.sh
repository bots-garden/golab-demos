#!/bin/bash

# -----------------------
# Install Extism
# -----------------------
sudo apt-get update -y
sudo apt-get install -y pkg-config

sudo apt install python3-pip -y
pip3 install poetry
pip3 install git+https://github.com/extism/cli

echo "export EXTISM_HOME=\"\$HOME/.local\"" >> ${HOME}/.bashrc
echo "export PATH=\"\$EXTISM_HOME/bin:\$PATH\"" >> ${HOME}/.bashrc

source ${HOME}/.bashrc

#extism install latest

extism --prefix=/usr/local install latest
pip3 install extism


# -----------------------
# Install Extism JS PDK
# -----------------------
#curl -O https://raw.githubusercontent.com/extism/js-pdk/main/install.sh
#sh install.sh

export TAG="v0.5.0"
export ARCH="aarch64"
export  OS="linux"
curl -L -O "https://github.com/extism/js-pdk/releases/download/$TAG/extism-js-$ARCH-$OS-$TAG.gz"
gunzip extism-js*.gz
sudo mv extism-js-* /usr/local/bin/extism-js
chmod +x /usr/local/bin/extism-js
