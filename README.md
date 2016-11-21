# SDCA-Makefile

----

##Présentation

Makefile distribué codé en [Golang](www.golang.org) et utilisant la librairie [Thrift](www.http://thrift.apache.org/).

----

##Pré-requis

* Avoir au minimum la version 1.7.x de Go et avoir setté les variables d'environnement Go
* Avoir installé thrift

----

##Procédure de déploiement sur Grid5000

1. Réserver un cluster de machines :

		# <nb_machine> = Nombre de machines à réserver
		# <time> = durée de la réservation
		oarsub -I -l nodes=<nb_machine>,walltime=<time> -t deploy

2. Déployer Debian Jessie sur le cluster :

		uniq $OAR_NODE_FILE > hostfile.txt
		kadeploy3 -f ./hostfile.txt -e jessie-x64-base -k

3. Cloner le projet :

		git clone https://github.com/Ataww/SDCA-Makefile.git

4. Lancer le script de configuration des machines :
	
		bash ./SDCA-Makefile/script/configure_host.sh

5. Lancer le script d'installation des packages avec taktuk :

		# <user_name> = votre nom d'utilisateur
		taktuk -l root -f hostfile.txt broadcast exec [ /home/<user_name>/SDCA-Makefile/script/taktuk_root.sh $USER $(id -u) ]

6. Se connecter à l'une des machine, qui sera la machine cliente, et lancer le script d'initialisation du projet :

		# Connection à la machine
		ssh <machine>
			
		# Lancer le script
		bash ./SDCA-Mafefile/script/init_project.sh

7. Aller dans un dossier contenant un Makefile. Y ajouter un fichier hostfile.txt (format d'une ligne host:9090).

		# Si le Makefile et le hostfile sont dans le dossier courant
		dmake
		
		# Si aucun hostfile n'est présent ou spécifié alors l'exécution sera lancée en local
		dmake -makefile=<path_to_makefile>
		
		# Sinon il est possible de les spécifier
		dmake -hostfile=<path_to_file> -makefile=<path_to_makefile>
		

##IMPORTANT

Nos scripts de déploiement ouvrent seulement le port 9090 sur les machines du cluster.
Il est donc impératif que les lignes du fichiers hostfile.txt soit de la forme :

	host_1:9090
	host_2:9090
	...
	host_N:9090

----

##Bugs connus

- Il se peut (dans 90% des cas ...) que taktuk ne remplisse pas entièrement sa tâche d'installation des packages sur toutes les machines. Dans ce cas, un message d'error avec le status 127 sera affiché. Pour régler le problème il faut relancer le script "init_taktuk.sh" sur la machine concernée.

----

##Contacts

- [Clément Taboulot](mailto:clement.taboulot@grenoble-inp.org)
- [Vincent Chenal](mailto:vincent.chenal@grenoble-inp.org)
- [Maxime Hagenbourger](mailto:maxime.hagenbourger@grenoble-inp.org)
- [Nathanaël Couret](mailto:nathanael.couret@grenoble-inp.org)