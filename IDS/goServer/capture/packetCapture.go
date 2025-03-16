package capture

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"goServer/database"
	"goServer/models"
	"goServer/utilities"
)

var (
	packetFlowData = make(map[string]*models.PacketFlow)
	flowMutex      sync.Mutex
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

	for packet := range packetSource.Packets() {
		processPacket(packet)
	}

}

func getPacketLayersInfo(packet gopacket.Packet) *models.PacketLayers {
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
		return nil
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
	return nil
}

// processPacket extracts flow features from packets
func processPacket(packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	udpLayer := packet.Layer(layers.LayerTypeUDP)

	if ipLayer == nil {
		return // Skip non-IP packets
	}

	ip, _ := ipLayer.(*layers.IPv4)

	// Extract transport layer details
	var srcPort, dstPort uint16
	var protocol uint8
	var flags map[string]int

	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		srcPort, dstPort, protocol = uint16(tcp.SrcPort), uint16(tcp.DstPort), 6 // TCP
		flags = map[string]int{
			"FIN": utilities.BoolToInt(tcp.FIN),
			"SYN": utilities.BoolToInt(tcp.SYN),
			"RST": utilities.BoolToInt(tcp.RST),
			"PSH": utilities.BoolToInt(tcp.PSH),
			"ACK": utilities.BoolToInt(tcp.ACK),
			"URG": utilities.BoolToInt(tcp.URG),
		}
	} else if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		srcPort, dstPort, protocol = uint16(udp.SrcPort), uint16(udp.DstPort), 17 // UDP
		flags = map[string]int{}
	} else {
		return // Skip non-TCP/UDP packets
	}

	// Construct Flow ID
	flowID := fmt.Sprintf("%s-%s-%d-%d-%d", ip.SrcIP, ip.DstIP, srcPort, dstPort, protocol)
	packetSize := len(packet.Data())
	currentTime := time.Now()

	flowMutex.Lock()
	defer flowMutex.Unlock()

	if flow, exists := packetFlowData[flowID]; exists {
		// Update existing flow
		flow.TotalFwdPackets++
		flow.TotalFwdBytes += packetSize
		flow.FwdPacketLengths = append(flow.FwdPacketLengths, packetSize)
		flow.EndTime = currentTime
	} else {
		// Create a new flow
		packetFlowData[flowID] = &models.PacketFlow{
			FlowID:           flowID,
			SourceIP:         ip.SrcIP.String(),
			SourcePort:       srcPort,
			DestinationIP:    ip.DstIP.String(),
			DestinationPort:  dstPort,
			Protocol:         protocol,
			StartTime:        currentTime,
			EndTime:          currentTime,
			TotalFwdPackets:  1,
			TotalFwdBytes:    packetSize,
			FwdPacketLengths: []int{packetSize},
			TCPFlags:         flags,
		}
	}

	// Store in database in real-time
	database.InsertPacket(packetFlowData[flowID])
}
