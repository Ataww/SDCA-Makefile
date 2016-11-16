package main

import (
	"SDCA-Makefile/compilationInterface"
	"crypto/tls"
	"fmt"
	"log"

	"git.apache.org/thrift.git/lib/go/thrift"
)

var busy []bool
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
		fmt.Println("There was a problem while running command : ", err)
	}
	fmt.Print("Server execute coommand and return status ", status, "\n")
	target.computing = false
	target.done = true
	busy[target.serverId] = false

	return err
}

func find_available_server() int {
	var nb_tested_id int = 0

	for nb_tested_id != len(busy) {
		if busy[current_server_id] == false {
			selected_id := current_server_id
			current_server_id = (current_server_id + 1) % len(busy)
			return selected_id
		}
		current_server_id = (current_server_id + 1) % len(busy)
		nb_tested_id++
	}
	return -1
}

func runClient(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, secure bool, hosts []string) error {
	var servers []*thrift.TTransport

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
	t1 := NewTarget("target1", "sleep", "3")
	t2 := NewTarget("target2", "sleep", "3")
	t3 := NewTarget("target3", "sleep", "3")
	t4 := NewTarget("target4", "sleep", "3")
	t5 := NewTarget("target5", "sleep", "3")
	t6 := NewTarget("target6", "sleep", "3")
	t7 := NewTarget("target7", "sleep", "3")

	t1.Add_Dependency(t2)
	t1.Add_Dependency(t3)

	t2.Add_Dependency(t4)
	t2.Add_Dependency(t5)

	t3.Add_Dependency(t6)
	t6.Add_Dependency(t7)
	//t1.Print(0)

	fmt.Printf("t1 : %t\n", t1.Is_Computable())
	fmt.Printf("t2 : %t\n", t2.Is_Computable())
	fmt.Printf("t3 : %t\n", t3.Is_Computable())
	fmt.Printf("t4 : %t\n", t4.Is_Computable())
	fmt.Printf("t5 : %t\n", t5.Is_Computable())
	fmt.Printf("t6 : %t\n", t6.Is_Computable())
	fmt.Printf("t7 : %t\n", t7.Is_Computable())
	/////////////////////////////

	for t1.done != true {
		var leaf = t1.Get_Leaf()
		if leaf != nil {
			if id_server := find_available_server(); id_server != -1 {
				fmt.Print("host ", id_server, " is executing : \n")
				leaf.Print(0)
				leaf.computing = true
				leaf.serverId = id_server
				busy[id_server] = true
				go handleTarget(servers[id_server], protocolFactory, leaf)
			}
		}
	}

	return nil
}
