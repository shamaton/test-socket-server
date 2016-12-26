#!/bin/sh
export version="1.7.4"

# base update
yum -y update
yum -y install wget

# install git
cd /etc/yum.repos.d/
wget http://wing-repo.net/wing/6/EL6.wing.repo
yum -y --enablerepo=wing install git

# install golang
wget https://storage.googleapis.com/golang/go${version}.linux-amd64.tar.gz
tar -C /usr/local -xzf go${version}.linux-amd64.tar.gz
rm go${version}.linux-amd64.tar.gz

# install git-completion
cd /root
wget https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.bash
mv git-completion.bash .git-completion.bash

# warning : firewall stop
service iptables stop

# set pathes
mkdir ~/.go
echo "export GOPATH=/root/.go" >> /root/.bashrc
echo "export PATH=\$PATH:/usr/local/go/bin:\$GOPATH/bin" >> /root/.bashrc
echo "source ~/.git-completion.bash" >> /root/.bashrc

# install necessary package
source ~/.bashrc
go get github.com/mattn/gom

# for git clone
echo -e "Host github.com\n\tStrictHostKeyChecking no\n" >> ~/.ssh/config

# setup project
cd /root
git clone git@github.com:shamaton/test-socket-server.git
cd test-socket-server/front/src
gom install
rm -rf vendor/bin vendor/pkg
cd ../../back/src
gom install
rm -rf vendor/bin vendor/pkg