package models

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type PacketFlow struct {
	FlowID            string
	SourceIP          string
	SourcePort        uint16
	DestinationIP     string
	DestinationPort   uint16
	Protocol          uint8
	Timestamp         string
	FlowDuration      int64
	TotalFwdPackets   int
	TotalBwdPackets   int
	MinSegmentSizeFwd int
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
