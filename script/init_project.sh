#!/bin/bash

user=$USER

echo -e "\e[33m############################"
echo          "# INITIALISATION DU PROJET #"
echo -e       "############################ \e[39m"

### USER DOES THIS PART
#Create application arborescence
mkdir /home/$user/Go
mkdir /home/$user/Go/src
mkdir /home/$user/Go/bin
mkdir /home/$user/Go/pkg

#Clone sources
git clone https://github.com/Ataww/SDCA-Makefile.git /home/$user/Go/src/SDCA-Makefile

#export GOPATH
export GOPATH=/home/$user/Go

#Get dependencies
echo "Installing Apache Thrift Go library..."
go get git.apache.org/thrift.git/lib/go/thrift/...
echo "Installing crypto/ssh library..."
go get golang.org/x/crypto/ssh

#Generate Go code from thrift file
cd /home/$user/Go/src/SDCA-Makefile/
thrift -r --gen go -out . compilation.thrift

#Compile the application
go install SDCA-Makefile/main

