package models

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"time"
)

type PacketFlow struct {
	FlowID          string
	SourceIP        string
	SourcePort      uint16
	DestinationIP   string
	DestinationPort uint16
	Protocol        uint8
	StartTime       time.Time
	EndTime         time.Time

	TotalFwdPackets int
	TotalBwdPackets int
	TotalFwdBytes   int
	TotalBwdBytes   int

	FwdPacketLengths []int
	BwdPacketLengths []int
	FwdTimestamps    []time.Time
	BwdTimestamps    []time.Time

	FwdHeaderLength int
	BwdHeaderLength int

	// TCP Flags Count (per direction)
	FwdFIN int
	FwdSYN int
	FwdRST int
	FwdPSH int
	FwdACK int
	FwdURG int

	BwdFIN int
	BwdSYN int
	BwdRST int
	BwdPSH int
	BwdACK int
	BwdURG int

	// Bulk features
	FwdAvgBytesBulk   float64
	FwdAvgPacketsBulk float64
	FwdAvgBulkRate    float64
	BwdAvgBytesBulk   float64
	BwdAvgPacketsBulk float64
	BwdAvgBulkRate    float64

	// Active/Idle time features
	ActiveMean float64
	ActiveStd  float64
	ActiveMax  float64
	ActiveMin  float64
	IdleMean   float64
	IdleStd    float64
	IdleMax    float64
	IdleMin    float64

	// legacy support for raw TCP flags if needed
	TCPFlags map[string]int
}

type PacketLayers struct {
	EthLayer    layers.Ethernet
	IPLayer     layers.IPv4
	IPv6Layer   layers.IPv6
	TCPLayer    layers.TCP
	UDPLayer    layers.UDP
	ICMPLayer   layers.ICMPv4
	ICMPv6Layer layers.ICMPv6
	Payload     gopacket.Payload
	Decoded     []gopacket.LayerType
}

type NetworkDashboardStats struct {
	Type              string     `json:"type"`
	TotalFlows        int        `json:"total_flows"`
	MaliciousFlows    int        `json:"malicious_flows"`
	LastMaliciousTime time.Time  `json:"last_malicious_time"`
	LastMaliciousFlow FlowDetail `json:"last_malicious_flow"`
}

type FlowDetail struct {
	FlowID    string    `json:"flow_id"`
	SrcIP     string    `json:"source_ip"`
	DstIP     string    `json:"destination_ip"`
	Protocol  string    `json:"protocol"`
	RiskScore float64   `json:"risk_score"`
	Timestamp time.Time `json:"timestamp"`
}

type NetworkPacket struct {
	Timestamp        *time.Time `json:"timestamp"`
	FlowID           string     `json:"flow_id"`
	SourceIP         string     `json:"src_ip"`
	DestinationIP    string     `json:"dst_ip"`
	SourcePort       *int       `json:"source_port"`
	DestinationPort  *int       `json:"destination_port"`
	Protocol         *int       `json:"protocol"`
	FlowDuration     *float64   `json:"flow_duration"`
	TotalFwdPackets  *int       `json:"total_fwd_packets"`
	TotalBwdPackets  *int       `json:"total_bwd_packets"`
	BytesPerSecond   *float64   `json:"flow_bytes_per_sec"`
	PacketsPerSecond *float64   `json:"flow_packets_per_sec"`
	Prediction       *int       `json:"prediction"`
}

type PaginatedPacketsResponse struct {
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Total    int             `json:"total"`
	Packets  []NetworkPacket `json:"packets"`
}
