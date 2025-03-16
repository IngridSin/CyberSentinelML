package capture

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"

	"goServer/models"
)

// StartPacketCapture begins packet sniffing on the given network interface
func StartPacketCapture(interfaceName string) {

	// Be aware, that OpenLive only supports microsecond resolution.
	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	fmt.Println("Capturing packets on:", interfaceName)

	packetFlows := make(map[string]*models.PacketFlow)

	for packet := range packetSource.Packets() {
		processPacket(packet, packetFlows)
	}

}

func processPacket(packet gopacket.Packet, flows map[string]*models.PacketFlow) {
	for _, layer := range packet.Layers() {
		fmt.Println("PACKET LAYER:", layer.LayerType())
	}
	////fmt.Println(packet)
	//if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
	//	fmt.Println("This is a TCP packet!")
	//	// Get actual TCP data from this layer
	//	tcp, _ := tcpLayer.(*layers.TCP)
	//	fmt.Printf("From src port %d to dst port %d\n", tcp.SrcPort, tcp.DstPort)
	//} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
	//	fmt.Println("This is a TCP packet!")
	//	// Get actual UDP data from this layer
	//	tcp, _ := udpLayer.(*layers.TCP)
	//	fmt.Printf("From src port %d to dst port %d\n", tcp.SrcPort, tcp.DstPort)
	//}
}
