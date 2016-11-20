#!/bin/bash

#taktuk -l root -f $OAR_FILE_NODES broadcast exec [ /home/vchenal/SDCA-Makefile/script/taktuk_root.sh $USER $(id -u) ]

user=$1
user_id=$2

echo -e "\e[33m#####################################"
echo          "# SCRIPT DE CONFIGURATION DES HOTES #"
echo -e       "##################################### \e[39m"

### ROOT DOES THIS PART

#Update sources
apt-get update

#Install go
wget https://storage.googleapis.com/golang/go1.7.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.7.3.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile
echo "GOPATH=/home/$user/Go" >> /etc/profile
source /etc/profile

#Install Apache thrift compile/install dependencies
apt-get install -y build-essential

#Compile/Install Apache Thrift
wget http://apache.crihan.fr/dist/thrift/0.9.3/thrift-0.9.3.tar.gz
tar xvzf thrift-0.9.3.tar.gz
cd thrift-0.9.3
./configure
make
make install

#Install git
apt-get install -y git



#Install application dependencies
wget https://johnvansickle.com/ffmpeg/builds/ffmpeg-git-64bit-static.tar.xz
tar xpvf ffmpeg-git-64bit-static.tar.xz -C /usr/local/bin --strip-components=1
apt-get install -y blender imagemagick unzip libav-tools

#Open server port
echo -e "\e[34m#Ouverture du port 9090 sur $hostname \e[39m"
iptables -A INPUT -p tcp --dport 9090 -j ACCEPT && iptables-save
