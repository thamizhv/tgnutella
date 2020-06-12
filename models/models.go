package models

import "github.com/google/uuid"

type DescriptorType uint8

type DescriptorHeader struct {
	ID                uuid.UUID      `json:"id"`
	PayloadDescriptor DescriptorType `json:"payloadDescriptor"`
	TTL               uint8          `json:"ttl"`
	Hops              uint8          `json:"hops"`
	Length            uint32         `json:"length"`
}

type Ping struct {
	IPAddress string
}

type Pong struct {
	Port      uint16
	IPAddress string
	Files     uint32
	Size      uint32
}

type Query struct {
	MinSpeedInKbps uint16
	Search         string
	IPAddress      string
}

type QueryHit struct {
	NumHits   uint8
	HTTPPort  uint16
	Port      uint16
	IPAddress string
	Speed     uint32
	ResultSet ResultSet
}

type ResultSet struct {
	FileIndex string
	FileSize  uint32
	FileName  string
}

type Peer struct {
	IP     string
	Port   string
	Files  uint32
	Size   uint32
	Active bool
}

type File struct {
	Name string
	Size int64
}

type PeerFile struct {
	HTTPPort  uint16
	Port      uint16
	IPAddress string
	Details   ResultSet
}

type FindFile struct {
	Name   string
	Active bool
}
