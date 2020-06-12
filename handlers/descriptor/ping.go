package descriptor

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"

	"github.com/thamizhv/tgnutella/constants"
	"github.com/thamizhv/tgnutella/helpers"
	"github.com/thamizhv/tgnutella/models"

	"github.com/google/uuid"
)

type pingDescriptor struct {
	*DescriptorHelper
	pongDescriptor DescriptorHandler
}

func NewPingDescriptor(helper *DescriptorHelper) DescriptorHandler {
	return &pingDescriptor{
		DescriptorHelper: helper,
		pongDescriptor:   NewPongDescriptor(helper),
	}
}

func (p *pingDescriptor) Send(_ string, _ *models.DescriptorHeader, _ []byte) error {

	payloadDescriptor, err := helpers.GetPayloadDescriptor(constants.PingDescriptor)
	if err != nil {
		return err
	}

	buf := p.Encode()
	header := &models.DescriptorHeader{
		ID:                uuid.New(),
		PayloadDescriptor: payloadDescriptor,
		TTL:               constants.DefaultServentTTL,
		Hops:              0,
		Length:            uint32(len(buf)),
	}

	data, err := p.EncodeDescriptorHeader(header)
	if err != nil {
		return err
	}

	data = append(data, buf...)

	// fmt.Println("ping send:", header.ID.String())

	// save my ping with my address
	key := fmt.Sprintf("%s-%s-%s", header.ID.String(), constants.PingDescriptor, "self")
	val := net.JoinHostPort(p.ip, strconv.FormatUint(uint64(p.port), 10))
	p.descriptorList.Set(key, []byte(val))

	return p.PropagateToPeers("", data)
}

func (p *pingDescriptor) Receive(localAddress string, header *models.DescriptorHeader, buf []byte) error {

	// fmt.Println("ping receive:", header.ID.String())
	// fmt.Println("ping receive data:", buf)

	// my ping I wont forward.
	key := fmt.Sprintf("%s-%s-%s", header.ID.String(), constants.PingDescriptor, "self")
	if p.descriptorList.Exists(key) {
		return nil
	}

	// if I have forwarded a ping, i wont forward again and also save intended address
	key = fmt.Sprintf("%s-%s", header.ID.String(), constants.PingDescriptor)
	if p.descriptorList.Exists(key) {
		return nil
	}

	ping := p.Decode(buf[constants.DescriptorHeaderLength:])
	remoteAddress := ping.IPAddress

	if remoteAddress == localAddress {
		return nil
	}

	err := p.pongDescriptor.Send(remoteAddress, header, nil)
	if err != nil {
		return err
	}

	p.descriptorList.Set(key, []byte(ping.IPAddress))

	// fmt.Println(localAddress)

	ping.IPAddress = localAddress
	encoded := p.Encode()

	header.Hops += 1
	header.TTL -= 1
	header.Length = uint32(len(encoded))

	if header.TTL == 0 {
		return nil
	}

	data, err := p.EncodeDescriptorHeader(header)
	if err != nil {
		return err
	}

	data = append(data, encoded...)

	return p.PropagateToPeers(remoteAddress, data)
}

func (p *pingDescriptor) Encode() []byte {
	buf := [14]byte{}
	ip := p.getEncodedIP(p.ip)
	copy(buf[0:], ip)
	binary.LittleEndian.PutUint16(buf[4:], p.port)

	// fmt.Println("encoded ping", buf[:])

	return buf[:]
}

func (p *pingDescriptor) Decode(buf []byte) *models.Ping {

	ping := &models.Ping{}
	ip := p.getDecodedIP(buf[:4])
	port := binary.LittleEndian.Uint16(buf[4:])
	ping.IPAddress = net.JoinHostPort(ip, strconv.FormatUint(uint64(port), 10))

	// fmt.Println("decoded ping", ping.IPAddress)
	return ping
}
