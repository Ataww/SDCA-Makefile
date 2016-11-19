package main

import (
	"SDCA-Makefile/compilationInterface"
	"crypto/tls"
	"fmt"
	"log"
	"path/filepath"
	"git.apache.org/thrift.git/lib/go/thrift"
	"sync"
	"os"
	"time"
	"golang.org/x/crypto/ssh"
	"bytes"
	"strings"
	"io/ioutil"
	"os/user"
)

var busy []bool
var mutex sync.Mutex
var current_server_id int = 0
var workingDir string

/*
Create thrift transport
 */
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

/*
Open a thrift connection
 */
func open_connection(t *thrift.TTransport) {
	err := (*t).Open()
	if err != nil {
		log.Fatal(err)
	}
}

/*
Close a thrift connection
 */
func close_connection(t *thrift.TTransport) {
	err := (*t).Close()
	if err != nil {
		log.Fatal(err)
	}
}

/*
Send an action to an other host
 */
func handleTarget(transport *thrift.TTransport, protocolFactory thrift.TProtocolFactory, target *Target, serverName string) (err error) {

	// Configuration of the command
	open_connection(transport)
	client := compilationInterface.NewCompilationServiceClientFactory(*transport, protocolFactory)
	command := compilationInterface.NewCommand()
	command.CommandLine = target.lineCommand
	command.WorkingDir = workingDir
	command.ID = target.id

	// Send the command
	status, err := client.ExecuteCommand(command)
	close_connection(transport)
	if err != nil {
		fmt.Println(serverName ," : There was a problem while running target ",target.id,": ", err.Error())
	}
	fmt.Println(serverName ," : Execute target ",target.id," and return status ",status)

	mutex.Lock()
	target.computing = false
	target.done = true
	busy[target.serverId] = false
	defer mutex.Unlock()

	return err
}

/*
Stop server
 */
func handleStop(transport *thrift.TTransport, protocolFactory thrift.TProtocolFactory, serverName string) (err error) {

	// Configuration
	open_connection(transport)
	client := compilationInterface.NewCompilationServiceClientFactory(*transport, protocolFactory)

	// Send the command
	err = client.Stop()
	close_connection(transport)

	if err != nil {
		fmt.Println(serverName, " : There was a problem while stoping server ", err.Error())
	}
	fmt.Println(serverName , " stop")

	return err
}

/*
Find an available server
 */
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

/*
Launch one server via SSH
 */
func launchServer(command string, hostname string, port string, config *ssh.ClientConfig) string {
	connection, error := ssh.Dial("tcp", fmt.Sprintf("%s:%s", strings.Split(hostname, ":")[0], port), config)
	if error != nil{
		fmt.Println("ERROR NEW CONNECTION ",fmt.Sprintf("%s:%s", strings.Split(hostname, ":")[0], port)," : ", error.Error())
	}
	defer connection.Close()
	session, error:= connection.NewSession()
	if error != nil{
		fmt.Println("ERROR NEW SESSION : ", error)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(command)

	return hostname + " : " + stdoutBuf.String()
}

/*
Get id_rsa key
 */
func getKeyFile() (key ssh.Signer, err error){
	usr, _ := user.Current()
	file := usr.HomeDir + "/.ssh/id_rsa"
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return
	}
	return
}

/*
Starts servers
 */
func startServers(hosts []string){
	usr, _ := user.Current()

	// "&>" is normal -> don't correct with "& >"
	cmd := "bash -c '"+usr.HomeDir+"/Go/bin/main -server=True -addr=0.0.0.0:9090 &> $(hostname)_server.out &'"

	results := make(chan string, 10)
	timeout := time.After(10 * time.Second)

	port := "22"

	key, err := getKeyFile();
	if err !=nil {
		panic(err)
	}
	// Create ssh config
	config := &ssh.ClientConfig{
		User: usr.Name,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	// Launch on each server
	for i := 0; i < len(hosts); i++{
		go func(hostname string, port string){
			results <- launchServer(cmd, hostname, port, config)
		}(hosts[i], port)
	}

	// Wait return
	for i := 0; i < len(hosts); i++ {
		select {
		case res := <- results:
			fmt.Println("Launch server on ",res)
		case <- timeout:
			fmt.Println("Time out ! : Exit ")
			os.Exit(1)
		}
	}


	time.Sleep(2 * time.Second)
}

/*
Main client function
 */
func runClient(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, secure bool, hosts []string, makefile string) error {
	var servers []*thrift.TTransport
	debut := time.Now()

	// Start all servers
	startServers(hosts);

	// Create thrift connection
	for i := 0; i < len(hosts); i++ {
		if err, server := createConnection(&transportFactory, hosts[i], secure); err != nil {
			fmt.Println("There was a problem while connecting to host " + hosts[i])
			log.Fatal(err)
			os.Exit(1) // Exit
		} else {
			servers = append(servers, server)
			busy = append(busy, false)
		}
	}

	// Parse Makefile
	root_target, _ := Parse(makefile)

	// Calculte working directory
	dir, _ := filepath.Abs(filepath.Dir(makefile))
	workingDir = dir


	// Job distribution while the target is not done
	for root_target.done != true {
		var leaf = root_target.Get_Leaf()
		if leaf != nil {
			if id_server := find_available_server(); id_server != -1 {
				if leaf.lineCommand != ""{
					// Execute the node command
					fmt.Println(hosts[id_server], " : Going to execute target ",leaf.id)
					mutex.Lock()
					leaf.computing = true
					leaf.serverId = id_server
					busy[id_server] = true
					mutex.Unlock()
					go handleTarget(servers[id_server], protocolFactory, leaf, hosts[id_server])
				}else{
					// There is no command to execute so this target is done
					fmt.Println("No command for target :", leaf.id)
					mutex.Lock()
					busy[id_server] = false
					leaf.done = true
					mutex.Unlock()
				}

			}
		}
	}

	// Stop each server
	for i := 0; i < len(servers); i++ {
		handleStop(servers[i], protocolFactory, hosts[i])
	}

	fin := time.Now()

	duree := fin - debut

	fmt.Println("==> TOTAL execution : ", duree.Second(), " sec.")

	// End
	return nil
}
