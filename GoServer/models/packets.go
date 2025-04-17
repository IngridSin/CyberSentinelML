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
