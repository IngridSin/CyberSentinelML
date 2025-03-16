package models

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
