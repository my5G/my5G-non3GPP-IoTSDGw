package udp_server

const (
	defaultForwaderPort = 1980
)

const (
	ChannelID1 = iota
	ChannelID2
	ChannelID3
	ChannelID4
	ChannelID5
	ChannelID6
	ChannelID7
	ChannelID8
	ChannelIDRecv
)

type Event int

const (
	EventUDPSender Event = iota
	EventUDPRecv
)
