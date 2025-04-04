package main

import (
	"fmt"
	"goServer/config"
	"goServer/database"
	"goServer/email"
	"goServer/websocket"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Initializing CyberSentinel...")

	go websocket.StartWebsocket()
	// Start SSH Tunnel
	//tunnel, err := ssh.CreateSSHTunnel()
	//if err != nil {
	//	log.Fatalf("Failed to create SSH tunnel: %v", err)
	//}
	//defer tunnel.Close()
	//
	//fmt.Println("SSH Tunnel established on local port:", tunnel.LocalPort)

	// Connect to Database via SSH Tunnel
	//database.ConnectDB("packets", tunnel.LocalPort, config.DBName)

	//// Start Packet Capture
	//interfaceName := config.NetworkInterface
	//go capture.StartPacketCapture(interfaceName)
	//

	database.ConnectDB("emails", 0, config.IMDBName)

	// Start Email Listener
	go email.StartEmailListener()

	// Wait for interrupt signal (e.g., Ctrl+C)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs // block until signal is received
}
