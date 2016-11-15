# SDCA-Makefile

----

##Présentation

Makefile distribué codé en [Golang](www.golang.org) en utilisant la librairie [Thrift](www.http://thrift.apache.org/).

----

##Pré-requis

* Avoir au minimum la version 1.7.x de Go et avoir setté les variables d'environnement Go
* Avoir installé thrift

----

##Compilation

1. Aller dans le workspace Go

		cd $GOPATH/src

2. Cloner le projet

		git clone https://github.com/Ataww/SDCA-Makefile.git

3. Générer les fichiers Go de service

		cd SDCA-Makefile && thrift -gen go -out . compilation.thrift

4. Récupérer la librairie Go de Thrift

		go get git.apache.org/thrift.git/lib/go/thrift/...

5. Compiler le main

		#L'exécutable est généré dans $GOBIN
		go install SDCA-Makefile/main

6. Lancer le serveur et le client

		#Pour svoir les options de lancement
		./main -help

		#Exemple d'exécution en mode serveur
		./main -server=True -addr=localhost:9090

		#Exemple d'éxécution en mode client
		./main -server=False -addr=localhost:9090

----

##Important

Si le client et le serveur sont lancés sur deux machines physiques distinctes. Il est important que le port sur lequel écoute le serveur soit ouvert.

###Exemple

Si mon serveur est lancé sur le port 9090 en tcp

	#Ouverture du port
	sudo iptables -A INPUT -p tcp --dport 9090 ACCEPT
	sudo iptables-save

	#Vous pouvez vérifier qur le port est bien ouvert
	netstat -pln
