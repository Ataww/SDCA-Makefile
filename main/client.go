package main

import (
	"SDCA-Makefile/compilationInterface"
	"crypto/tls"
	"fmt"
	"log"

	"git.apache.org/thrift.git/lib/go/thrift"
)

func createConnection(transportFactory *thrift.TTransportFactory, addr string, secure bool) (error, *thrift.TTransport) {
	var transport = new(thrift.TTransport)
	var err error

	if secure {
		cfg := new(tls.Config)
		cfg.InsecureSkipVerify = true
		*transport, err = thrift.NewTSSLSocket(addr, cfg)
	} else {
		*transport, err = thrift.NewTSocket(addr)
	}
	if err != nil {
		return err, nil
	}

	*transport = (*transportFactory).GetTransport(*transport)
	defer (*transport).Close()
	if err := (*transport).Open(); err != nil {
		return err, nil
	}

	return nil, transport
}

func open_connection(t *thrift.TTransport) {
	err := (*t).Open()
	if err != nil {
		log.Fatal(err)
	}
}

func close_connection(t *thrift.TTransport) {
	err := (*t).Close()
	if err != nil {
		log.Fatal(err)
	}
}

func handleTarget(client *compilationInterface.CompilationServiceClient, target *Target) (err error) {
	command := compilationInterface.NewCommand()
	command.Program = target.program
	command.Arguments = target.args
	status, err := client.ExecuteCommand(command)
	if err != nil {
		fmt.Println("There was a problem while running command : ", err)
	}
	fmt.Print("Server execute coommand and return status ", status, "\n")
	return err
}

func runClient(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, secure bool, hosts []string) error {
	var servers []*thrift.TTransport

	for i := 0; i < len(hosts); i++ {
		if err, server := createConnection(&transportFactory, hosts[i], secure); err != nil {
			fmt.Println("There was a problem while connecting to host " + hosts[i])
			log.Fatal(err)
		} else {
			servers = append(servers, server)
		}
	}

	//Parsing (Création de l'arbre de target)
	var root_target *Target = NewTarget("target1", "touch", "toto.txt")

	//Distribute tasks
	/*for !root_target.done {

	}*/

	open_connection(servers[0])
	handleTarget(compilationInterface.NewCompilationServiceClientFactory(*servers[0], protocolFactory), root_target)
	close_connection(servers[0])

	return nil
}
