#!/bin/bash

user=$USER
user_id=$(id -u)

echo -e "\e[33m#####################################"
echo          "# SCRIPT DE CONFIGURATION DES HOTES #"
echo -e       "##################################### \e[39m"

for hostname in $( uniq $OAR_NODE_FILE);
do
   echo -e "\e[31m##$hostname \e[39m"

   echo -e "\e[34m#Ouverture du port 9090 sur $hostname \e[39m"
   ssh root@$hostname 'iptables -A INPUT -p tcp --dport 9090 -j ACCEPT && iptables-save'

   echo -e "\e[34m#Montage du home nfs sur $hostname \e[39m"
   ssh root@$hostname "apt-get install -y nfs-common &&  mount -o rw,nfsvers=3,hard,intr,async,noatime,nodev,nosuid,auto,rsize=32768,wsize=32768 nfs.grenoble.grid5000.fr:/export/home/ /home/"
 
   echo -e "\e[34m#Creation de $user sur $hostname \e[39m"
   ssh root@$hostname "adduser -uid $user_id --no-create-home --disabled-password --gecos \"\" $user"
done
