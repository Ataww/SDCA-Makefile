package main

import (
	"compilationInterface"
	"crypto/tls"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
)

func handleClient(client *compilationInterface.CompilationServiceClient) (err error) {
	command := compilationInterface.NewCommand()
	command.Program = "touch"
	command.Arguments = "toto.txt"
	status, err := client.ExecuteCommand(command)
	if err != nil {
		fmt.Println("ERROR TA MERE :", err)
	}
	fmt.Print("Server execute coommand and return status ", status, "\n")
	return err
}

func runClient(transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, addr string, secure bool) error {
	var transport thrift.TTransport
	var err error
	/*transport, err = thrift.NewTSocket(addr)
	if err != nil {
		fmt.Println("Error opening socket:", err)
		return err
	}*/

	if secure {
		cfg := new(tls.Config)
		cfg.InsecureSkipVerify = true
		transport, err = thrift.NewTSSLSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTSocket(addr)
	}
	if err != nil {
		fmt.Println("Error opening socket:", err)
		return err
	}

	transport = transportFactory.GetTransport(transport)
	defer transport.Close()
	if err := transport.Open(); err != nil {
		return err
	}
	return handleClient(compilationInterface.NewCompilationServiceClientFactory(transport, protocolFactory))
}
