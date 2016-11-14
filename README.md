# SDCA-Makefile

----

##Présentation

Makefile distribué codé en [Golang](www.golang.org) en utilisant la librairie [Thrift](www.http://thrift.apache.org/).

----

##Pré-requis

* Avoir au minimum la version 1.7.x de Go
* Avoir installer thrift

----

##Compilation

1. Cloner le projet

		git clone https://github.com/Ataww/SDCA-Makefile.git
	
2. Générer les fichiers Go de service

		thrift -gen go compilation.thrift

3. Déplacer le package généré à la source du projet

		mv gen-go/compilationInterface ./

4. Déplacer les fichiers go à compiler dans l'environnement de travail Go

		cp SDCA-Makefile/* $GOPATH/src

5. Récupérer la librairie Go de Thrift

		go get git.apache.org/thrift.git/lib/go/thrift/...

6. Compiler le package compilationInterface

		go build compilationInterface

7. Compiler le main

		#L'exécutable est généré dans $GOBIN
		go install main

8. Lancer le serveur et le client

		#Pour svoir les options de lancement
		./main -help
		
		#Exemple d'exécution en mode serveur
		./main -server=True -addr=localhost:9090
		
		#Exemple d'éxécution en mode client
		./main -server=False -addr=localhost:9090

----

##Important

Si le client et le serveur sont lancer sur deux machines physique distincte. Il est important que le port sur lequel écoute le serveur soit ouvert.

###Exemple

Si mon serveur est lancé sur le port 9090 en tcp

	#Ouverture du port
	sudo iptables -A INPUT -p tcp --dport 9090 ACCEPT
	sudo iptables-save
	
	#Vous pouvez vérifier qur le port est bien ouvert
	netstat -pln