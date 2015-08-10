#!/bin/bash

apt-get install g++ gcc -y

apt-get install gdb -y

apt-get install screen vim -y

apt-get install linux-tools -y

apt-get install curl -y

apt-get install git -y

apt-get install cmake -y


apt-get install libssl-dev -y

apt-get install libevent -y

apt-get install libevent-dev -y

apt-get install libcurl4-dev -y

apt-get install libcurl4-openssl-dev -y


apt-get install libicu-dev -y

apt-get install libsnappy-dev -y

apt-get install libv8-dev -y

apt-get install erlang -y

wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
tar -xzf go1.4.2.linux-amd64.tar.gz -C /usr/local/
echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/gopath.sh

curl https://storage.googleapis.com/git-repo-downloads/repo > /usr/bin/repo
chmod a+x /usr/bin/repo
