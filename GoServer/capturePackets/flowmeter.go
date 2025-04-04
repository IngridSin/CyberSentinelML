package capture

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"goServer/models"
)

var (
	flowMap   = make(map[string]*models.PacketFlow)
	flowMutex sync.Mutex
)

func processPacket(packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return
	}
	ip := ipLayer.(*layers.IPv4)

	var srcPort, dstPort uint16
	var protocol uint8
	var headerLen int
	var isTCP, isUDP bool
	flags := make(map[string]int)

	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp := tcpLayer.(*layers.TCP)
		srcPort, dstPort = uint16(tcp.SrcPort), uint16(tcp.DstPort)
		protocol = 6
		isTCP = true
		headerLen = int(tcp.DataOffset) * 4
		flags = map[string]int{
			"FIN": boolToInt(tcp.FIN),
			"SYN": boolToInt(tcp.SYN),
			"RST": boolToInt(tcp.RST),
			"PSH": boolToInt(tcp.PSH),
			"ACK": boolToInt(tcp.ACK),
			"URG": boolToInt(tcp.URG),
		}
	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp := udpLayer.(*layers.UDP)
		srcPort, dstPort = uint16(udp.SrcPort), uint16(udp.DstPort)
		protocol = 17
		isUDP = true
		headerLen = 8
	} else {
		return
	}

	timestamp := packet.Metadata().Timestamp
	length := len(packet.Data())
	key := fmt.Sprintf("%s-%s-%d-%d-%d", ip.SrcIP, ip.DstIP, srcPort, dstPort, protocol)
	revKey := fmt.Sprintf("%s-%s-%d-%d-%d", ip.DstIP, ip.SrcIP, dstPort, srcPort, protocol)

	direction := "fwd"
	flowMutex.Lock()
	flow, exists := flowMap[key]
	if !exists {
		if rev, revExists := flowMap[revKey]; revExists {
			flow = rev
			direction = "bwd"
		} else {
			flow = &PacketFlow{
				FlowID:    key,
				StartTime: timestamp,
				EndTime:   timestamp,
				Protocol:  protocol,
				SrcIP:     ip.SrcIP.String(),
				DstIP:     ip.DstIP.String(),
				SrcPort:   srcPort,
				DstPort:   dstPort,
				Fwd:       FlowDirection{Flags: map[string]int{}},
				Bwd:       FlowDirection{Flags: map[string]int{}},
			}
			flowMap[key] = flow
		}
	}

	flow.EndTime = timestamp
	dir := &flow.Fwd
	if direction == "bwd" {
		dir = &flow.Bwd
	}
	dir.Bytes += length
	dir.Packets++
	dir.PacketLengths = append(dir.PacketLengths, length)
	dir.Timestamps = append(dir.Timestamps, timestamp)
	dir.HeaderLengths = append(dir.HeaderLengths, headerLen)
	for k, v := range flags {
		dir.Flags[k] += v
	}
	flowMutex.Unlock()
}

func calculateIATs(timestamps []time.Time) []float64 {
	var iats []float64
	for i := 1; i < len(timestamps); i++ {
		iat := timestamps[i].Sub(timestamps[i-1]).Seconds()
		iats = append(iats, iat)
	}
	return iats
}

func calculateStats(values []int) (min, max int, mean, std float64) {
	if len(values) == 0 {
		return
	}
	min, max = values[0], values[0]
	sum := 0
	for _, v := range values {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	mean = float64(sum) / float64(len(values))
	for _, v := range values {
		std += math.Pow(float64(v)-mean, 2)
	}
	std = math.Sqrt(std / float64(len(values)))
	return
}

func exportFlowsToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"Flow ID", "Src IP", "Dst IP", "Src Port", "Dst Port", "Protocol", "Flow Duration",
		"Fwd Packet Count", "Bwd Packet Count", "Fwd Bytes", "Bwd Bytes",
		"Fwd Pkt Len Max", "Fwd Pkt Len Min", "Fwd Pkt Len Mean", "Fwd Pkt Len Std",
		"Bwd Pkt Len Max", "Bwd Pkt Len Min", "Bwd Pkt Len Mean", "Bwd Pkt Len Std",
	}
	writer.Write(header)

	flowMutex.Lock()
	defer flowMutex.Unlock()
	for _, flow := range flowMap {
		duration := flow.EndTime.Sub(flow.StartTime).Seconds()

		fwdMin, fwdMax, fwdMean, fwdStd := calculateStats(flow.Fwd.PacketLengths)
		bwdMin, bwdMax, bwdMean, bwdStd := calculateStats(flow.Bwd.PacketLengths)

		row := []string{
			flow.FlowID, flow.SrcIP, flow.DstIP,
			strconv.Itoa(int(flow.SrcPort)), strconv.Itoa(int(flow.DstPort)), strconv.Itoa(int(flow.Protocol)),
			fmt.Sprintf("%.6f", duration),
			strconv.Itoa(flow.Fwd.Packets), strconv.Itoa(flow.Bwd.Packets),
			strconv.Itoa(flow.Fwd.Bytes), strconv.Itoa(flow.Bwd.Bytes),
			strconv.Itoa(fwdMax), strconv.Itoa(fwdMin), fmt.Sprintf("%.2f", fwdMean), fmt.Sprintf("%.2f", fwdStd),
			strconv.Itoa(bwdMax), strconv.Itoa(bwdMin), fmt.Sprintf("%.2f", bwdMean), fmt.Sprintf("%.2f", bwdStd),
		}
		writer.Write(row)
	}
	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
