#!/usr/bin/env sh

# install mozilla sops
wget https://github.com/mozilla/sops/releases/download/3.2.0/sops_3.2.0_amd64.deb
dpkg -i sops_3.2.0_amd64.deb

/oya/install.sh

cd /oya/
oya run test
