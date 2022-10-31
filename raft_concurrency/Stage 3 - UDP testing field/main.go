package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Starting the server...")
	companion := os.Getenv("COMPANION")
	if companion == "node2" {
		fmt.Println("Node 1 is running")
		conn, err := net.Dial("udp", "node2:8080")
		if err != nil {
			fmt.Println("Error connecting to node 2")
		}
		fmt.Println("Connected to node 2")
		fmt.Println("Attempting to print!")
		fmt.Fprintf(conn, "Hello from node 1")
		conn.Close()
	} else {
		fmt.Println("Node 2 is running")
		p := make([]byte, 2048)
		address, _ := net.ResolveUDPAddr("udp", "node2:8080")
		fmt.Println("Resolved address successfully")
		conn, _ := net.ListenUDP("udp", address)
		fmt.Println("Listening on node 2")
		for {
			_, _, _ = conn.ReadFromUDP(p)
			fmt.Println(string(p))
		}
	}
}
