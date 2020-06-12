package constants

import "time"

const (
	UsageText = "Usage: tgnutella <servent_port>"

	HelpText = `	help			lists details of available commands
	open <host:port>	connects to a node on the network
	close <id>		closes connection by connection id (see info command)
	info connections	prints list of connected hosts with an id for each
	find <keyword>		search files on the network and lists results with an id for each entry
	get <id>		download a file by id
	`
	CommandNotFound = "Command not found. Please try below commands"

	CmdTypeHelp = "help"

	CmdTypeOpen = "open"

	CmdTypeClose = "close"

	CmdTypeInfo = "info"

	CmdTypeFind = "find"

	CmdTypeGet = "get"

	ArgTypeConnections = "connections"

	NetworkTypeTCP = "tcp"

	LocalHost = "127.0.0.1"

	PeersList = "peers"

	ProtocolVersion = "0.4"

	ConnectionRequest = "GNUTELLA CONNECT/" + ProtocolVersion

	ConnectionResponse = "GNUTELLA OK"

	PingInterval = time.Second * 10

	DefaultServentTTL = 10
)
