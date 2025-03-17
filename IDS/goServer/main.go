package main

import (
	"fmt"
	"goServer/database"
	"log"

	"goServer/capture"
	"goServer/config"
	"goServer/ssh"
)

func main() {
	fmt.Println("Initializing Packet Capture Application...")

	// Start SSH Tunnel
	tunnel, err := ssh.CreateSSHTunnel()
	if err != nil {
		log.Fatalf("Failed to create SSH tunnel: %v", err)
	}
	defer tunnel.Close()

	fmt.Println("SSH Tunnel established on local port:", tunnel.LocalPort)

	// Connect to Database via SSH Tunnel
	database.ConnectDB(tunnel.LocalPort)

	// Start Packet Capture
	interfaceName := config.NetworkInterface // Change based on your system
	go capture.StartPacketCapture(interfaceName)

	// Keep running
	select {}
}
