package main

import (
	"SDCA-Makefile/compilationInterface"
	"crypto/tls"
	"fmt"
	"log"

	"git.apache.org/thrift.git/lib/go/thrift"
	"sync"
)

var busy []bool
var mutex sync.Mutex
var current_server_id int = 0

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

func handleTarget(transport *thrift.TTransport, protocolFactory thrift.TProtocolFactory, target *Target) (err error) {

	open_connection(transport)
	client := compilationInterface.NewCompilationServiceClientFactory(*transport, protocolFactory)

	command := compilationInterface.NewCommand()
	command.Program = target.program
	command.Arguments = target.args
	command.ID = target.id
	status, err := client.ExecuteCommand(command)
	close_connection(transport)
	if err != nil {
		fmt.Println(target.serverId ,"There was a problem while running target ",target.id,": ", err)
	}
	fmt.Print(target.serverId ," - Server execute target",target.id," and return status ", status, "\n")

	mutex.Lock()
	target.computing = false
	target.done = true
	busy[target.serverId] = false
	defer mutex.Unlock()

	return err
}

func find_available_server() int {
	mutex.Lock()
	var nb_tested_id int = 0

	for nb_tested_id != len(busy) {
		if busy[current_server_id] == false {
			selected_id := current_server_id
			current_server_id = (current_server_id + 1) % len(busy)
			defer mutex.Unlock()
			return selected_id
		}
		current_server_id = (current_server_id + 1) % len(busy)
		nb_tested_id++
	}
	defer mutex.Unlock()
	return -1
}

func runClient(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, secure bool, hosts []string, makefile string) error {
	var servers []*thrift.TTransport


	/////////////////////////////
	//   CONNECTIONS CREATION
	/////////////////////////////
	for i := 0; i < len(hosts); i++ {
		if err, server := createConnection(&transportFactory, hosts[i], secure); err != nil {
			fmt.Println("There was a problem while connecting to host " + hosts[i])
			log.Fatal(err)
		} else {
			servers = append(servers, server)
			busy = append(busy, false)
		}
	}

	/////////////////////////////
	//   PARSING
	/////////////////////////////
	root_target, _ := Parse(makefile)

	/////////////////////////////
	//   JOB DISTRIBUTION
	/////////////////////////////
	for root_target.done != true {
		var leaf = root_target.Get_Leaf()
		if leaf != nil {
			if id_server := find_available_server(); id_server != -1 {
				fmt.Print("-----------------------------------\n")
				fmt.Print("host ", id_server, " is executing : ",leaf.id,"\n")

				mutex.Lock()
				leaf.computing = true
				leaf.serverId = id_server
				busy[id_server] = true
				mutex.Unlock()

				if leaf.program != ""{
					go handleTarget(servers[id_server], protocolFactory, leaf)
				}else{
					fmt.Println("No command")

					mutex.Lock()
					busy[id_server] = false
					leaf.computing = false
					leaf.done = true
					mutex.Unlock()
				}

			}
		}
	}

	return nil
}
