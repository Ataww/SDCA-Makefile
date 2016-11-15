package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"git.apache.org/thrift.git/lib/go/thrift"
)

func Usage() {
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
}

func main() {
	flag.Usage = Usage
	server := flag.Bool("server", false, "Run server")
	protocol := flag.String("P", "binary", "Specify the protocol (binary, compact, json, simplejson)")
	framed := flag.Bool("framed", false, "Use framed transport")
	buffered := flag.Bool("buffered", false, "Use buffered transport")
	addr := flag.String("addr", "localhost:9090", "Address to listen to")
	secure := flag.Bool("secure", false, "Use tls secure transport")
	hostfile := flag.String("hostfile", "hostfile.txt", "Specify hostfile")

	flag.Parse()

	var protocolFactory thrift.TProtocolFactory
	switch *protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		fmt.Fprint(os.Stderr, "Invalid protocol specified", protocol, "\n")
		Usage()
		os.Exit(1)
	}

	var transportFactory thrift.TTransportFactory
	if *buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if *framed {
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	if *server {
		if err := runServer(transportFactory, protocolFactory, *addr, *secure); err != nil {
			fmt.Println("error running server:", err)
		}
	} else {
		var hosts []string

		fmt.Println("Checking hosts inside " + *hostfile + " file")
		f, _ := os.Open(*hostfile)
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			str := scanner.Text()
			hosts = append(hosts, str)
			fmt.Println("Found host : " + str)
		}

		if len(hosts) == 0 {
			fmt.Println("No hosts were found")
		}

		if err := runClient(transportFactory, protocolFactory, *secure, hosts); err != nil {
			fmt.Println("error running client:", err)
		}
	}
}
