package main

import (
	"fmt"
	"goServer/buffer"
	capture "goServer/capturePackets"
	"goServer/config"
	"goServer/database"
	"goServer/email"
	"goServer/ssh"
	"goServer/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Initializing CyberSentinel...")

	fmt.Println("Initializing Redis...")

	buffer.InitRedis()

	go websocket.StartWebsocket()
	go websocket.StartNetworkStatsBroadcaster()

	//Start SSH Tunnel
	tunnel, err := ssh.CreateSSHTunnel()
	if err != nil {
		log.Fatalf("Failed to create SSH tunnel: %v", err)
	}
	defer tunnel.Close()

	fmt.Println("SSH Tunnel established on local port:", tunnel.LocalPort)

	// Connect to Database via SSH Tunnel
	database.ConnectDB("packets", tunnel.LocalPort, config.DBName)

	//// Start Packet Capture
	interfaceName := config.NetworkInterface
	go capture.StartPacketCapture(interfaceName)

	database.ConnectDB("emails", tunnel.LocalPort, config.IMDBName)

	// Start Email Listener
	go email.StartEmailListener()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs // block until signal is received
}
