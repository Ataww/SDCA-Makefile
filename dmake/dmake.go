package main

import (
	"bufio"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"os"
	"strings"
)

func Usage() {
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
}

/*
Main function
*/
func main() {
	flag.Usage = Usage
	server := flag.Bool("server", false, "Run server")
	protocol := flag.String("P", "binary", "Specify the protocol (binary, compact, json, simplejson)")
	framed := flag.Bool("framed", false, "Use framed transport")
	buffered := flag.Bool("buffered", false, "Use buffered transport")
	addr := flag.String("addr", "localhost:9090", "Address to listen to")
	secure := flag.Bool("secure", false, "Use tls secure transport")
	hostfile := flag.String("hostfile", "hostfile.txt", "Specify hostfile")
	makefile := flag.String("makefile", "Makefile", "Specify the Makefile path.")

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
		// Launch server
		if err := runServer(transportFactory, protocolFactory, *addr, *secure); err != nil {
			fmt.Println("error running server:", err)
		}
	} else {
		var hosts []string

		fmt.Println("Checking hosts inside " + *hostfile + " file")

		if _, err := os.Stat(*hostfile); os.IsNotExist(err) {
			hosts = append(hosts, "localhost:9090")
			fmt.Println("Using localhost as server : localhost:9090 ")
		} else {
			f, _ := os.Open(*hostfile)
			scanner := bufio.NewScanner(f)

			// Read hostfile
			for scanner.Scan() {
				str := scanner.Text()
				if strings.HasPrefix(str, "#") == false {
					hosts = append(hosts, str)
					fmt.Println("Found host : " + str)
				}
			}

			if len(hosts) == 0 {
				// No host = exit
				fmt.Println("No hosts were found")
				os.Exit(1)
			}
		}

		// Launch client
		if err := runClient(transportFactory, protocolFactory, *secure, hosts, *makefile); err != nil {
			fmt.Println("error running client:", err)
		}
	}
}
