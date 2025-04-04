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
