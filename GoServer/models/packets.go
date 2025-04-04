package models

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"time"
)

type PacketFlow struct {
	FlowID           string
	SourceIP         string
	SourcePort       uint16
	DestinationIP    string
	DestinationPort  uint16
	Protocol         uint8
	StartTime        time.Time
	EndTime          time.Time
	TotalFwdPackets  int
	TotalBwdPackets  int
	TotalFwdBytes    int
	TotalBwdBytes    int
	FwdPacketLengths []int
	BwdPacketLengths []int
	TCPFlags         map[string]int
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
