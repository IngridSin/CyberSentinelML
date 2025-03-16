package capture

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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
	// Use the constants.PacketLayers struct
	var packetLayers models.PacketLayers

	// Create a Decoding Layer Parser
	parser := gopacket.NewDecodingLayerParser(
		layers.LayerTypeEthernet,
		&packetLayers.EthLayer, &packetLayers.IPLayer, &packetLayers.IPv6Layer, &packetLayers.TCPLayer,
		&packetLayers.UDPLayer, &packetLayers.ICMPLayer, &packetLayers.ICMPv6Layer, &packetLayers.Payload,
	)

	// Decode packet
	err := parser.DecodeLayers(packet.Data(), &packetLayers.Decoded)
	if err != nil {
		log.Printf("Decoding failed: %v", err)
		return
	}

	// Print decoded layers
	fmt.Println("=== New Packet ===")
	for _, layerType := range packetLayers.Decoded {
		switch layerType {
		case layers.LayerTypeEthernet:
			fmt.Printf("Ethernet Src: %s, Dst: %s\n", packetLayers.EthLayer.SrcMAC, packetLayers.EthLayer.DstMAC)
		case layers.LayerTypeIPv4:
			fmt.Printf("IPv4 Src: %s, Dst: %s\n", packetLayers.IPLayer.SrcIP, packetLayers.IPLayer.DstIP)
		case layers.LayerTypeIPv6:
			fmt.Printf("IPv6 Src: %s, Dst: %s\n", packetLayers.IPv6Layer.SrcIP, packetLayers.IPv6Layer.DstIP)
		case layers.LayerTypeTCP:
			fmt.Printf("TCP SrcPort: %d, DstPort: %d, SYN: %t, ACK: %t\n",
				packetLayers.TCPLayer.SrcPort, packetLayers.TCPLayer.DstPort, packetLayers.TCPLayer.SYN, packetLayers.TCPLayer.ACK)
		case layers.LayerTypeUDP:
			fmt.Printf("UDP SrcPort: %d, DstPort: %d\n", packetLayers.UDPLayer.SrcPort, packetLayers.UDPLayer.DstPort)
		case layers.LayerTypeICMPv4:
			fmt.Printf("ICMPv4 Type: %d, Code: %d\n", packetLayers.ICMPLayer.TypeCode.Type(), packetLayers.ICMPLayer.TypeCode.Code())
		case layers.LayerTypeICMPv6:
			fmt.Printf("ICMPv6 Type: %d, Code: %d\n", packetLayers.ICMPv6Layer.TypeCode.Type(), packetLayers.ICMPv6Layer.TypeCode.Code())
		case gopacket.LayerTypePayload:
			fmt.Printf("Payload: %x\n", packetLayers.Payload.Payload())
		}
	}
}

//for _, layer := range packet.Layers() {
//	fmt.Println("PACKET LAYER:", layer.LayerType())
//}
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
//}
