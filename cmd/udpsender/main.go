package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	address, err := net.ResolveUDPAddr("udp", "localhost:8182")

	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, address)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Write([]byte(str))

		if err != nil {
			log.Fatal(err)
		}

	}
}
