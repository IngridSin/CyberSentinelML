package capture

import (
	"fmt"
	"goServer/buffer"
	"log"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
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
	var headerLen int
	var flags map[string]int

	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		srcPort, dstPort, protocol = uint16(tcp.SrcPort), uint16(tcp.DstPort), 6 // TCP
		headerLen = int(tcp.DataOffset) * 4
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
		headerLen = 8                                                             // Fixed UDP header
		flags = map[string]int{}
	} else {
		return // Skip non-TCP/UDP packets
	}

	flowID := fmt.Sprintf("%s-%s-%d-%d-%d-%d", ip.SrcIP, ip.DstIP, srcPort, dstPort, protocol, time.Now().UnixNano())
	packetSize := len(packet.Data())
	currentTime := time.Now()

	flowMutex.Lock()
	defer flowMutex.Unlock()

	flow, exists := packetFlowData[flowID]
	if exists {
		flow.EndTime = currentTime

		if flow.SourceIP == ip.SrcIP.String() {
			// Forward
			flow.TotalFwdPackets++
			flow.TotalFwdBytes += packetSize
			flow.FwdPacketLengths = append(flow.FwdPacketLengths, packetSize)
			flow.FwdTimestamps = append(flow.FwdTimestamps, currentTime)
			flow.FwdHeaderLength += headerLen

			// TCP flag count (Fwd)
			if protocol == 6 {
				flow.FwdFIN += flags["FIN"]
				flow.FwdSYN += flags["SYN"]
				flow.FwdRST += flags["RST"]
				flow.FwdPSH += flags["PSH"]
				flow.FwdACK += flags["ACK"]
				flow.FwdURG += flags["URG"]
			}
		} else {
			// Backward
			flow.TotalBwdPackets++
			flow.TotalBwdBytes += packetSize
			flow.BwdPacketLengths = append(flow.BwdPacketLengths, packetSize)
			flow.BwdTimestamps = append(flow.BwdTimestamps, currentTime)
			flow.BwdHeaderLength += headerLen

			// TCP flag count (Bwd)
			if protocol == 6 {
				flow.BwdFIN += flags["FIN"]
				flow.BwdSYN += flags["SYN"]
				flow.BwdRST += flags["RST"]
				flow.BwdPSH += flags["PSH"]
				flow.BwdACK += flags["ACK"]
				flow.BwdURG += flags["URG"]
			}
		}
	} else {
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
			FwdTimestamps:    []time.Time{currentTime},
			FwdHeaderLength:  headerLen,
			TCPFlags:         flags,
		}
		flow = packetFlowData[flowID]
	}

	// calculate features after flow is updated
	avgBytes, avgPkts, bulkRate := detectBulkFeatures(flow.FwdTimestamps, flow.FwdPacketLengths)
	flow.FwdAvgBytesBulk = avgBytes
	flow.FwdAvgPacketsBulk = avgPkts
	flow.FwdAvgBulkRate = bulkRate

	bBytes, bPkts, bRate := detectBulkFeatures(flow.BwdTimestamps, flow.BwdPacketLengths)
	flow.BwdAvgBytesBulk = bBytes
	flow.BwdAvgPacketsBulk = bPkts
	flow.BwdAvgBulkRate = bRate

	activeMean, activeStd, activeMax, activeMin, idleMean, idleStd, idleMax, idleMin := detectActiveIdleFeatures(
		mergeAndSort(flow.FwdTimestamps, flow.BwdTimestamps),
	)
	flow.ActiveMean = activeMean
	flow.ActiveStd = activeStd
	flow.ActiveMax = activeMax
	flow.ActiveMin = activeMin
	flow.IdleMean = idleMean
	flow.IdleStd = idleStd
	flow.IdleMax = idleMax
	flow.IdleMin = idleMin

	// Push updated flow to Redis
	buffer.InsertPacket(flow)
}
