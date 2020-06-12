package descriptor

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/thamizhv/tgnutella/constants"
	"github.com/thamizhv/tgnutella/helpers"
	"github.com/thamizhv/tgnutella/models"
)

type pongDescriptor struct {
	*DescriptorHelper
}

func NewPongDescriptor(helper *DescriptorHelper) DescriptorHandler {
	return &pongDescriptor{
		DescriptorHelper: helper,
	}
}

func (p *pongDescriptor) Send(remoteAddress string, previous *models.DescriptorHeader, _ []byte) error {

	// fmt.Println("pong send:", previous.ID.String())

	// first pong saves intended address and mainly for reference
	key := fmt.Sprintf("%s-%s", previous.ID.String(), constants.PongDescriptor)
	if p.descriptorList.Exists(key) {
		return nil
	}

	payloadDescriptor, err := helpers.GetPayloadDescriptor(constants.PongDescriptor)
	if err != nil {
		return err
	}

	buf := p.Encode()

	header := &models.DescriptorHeader{
		ID:                previous.ID,
		PayloadDescriptor: payloadDescriptor,
		TTL:               previous.Hops + 1,
		Hops:              0,
		Length:            uint32(len(buf)),
	}

	data, err := p.EncodeDescriptorHeader(header)
	if err != nil {
		return err
	}

	data = append(data, buf...)

	_, err = helpers.SendAndReceiveData(remoteAddress, data, false)
	if err != nil {
		return err
	}

	p.descriptorList.Set(key, []byte(remoteAddress))

	id := helpers.GetHash(remoteAddress)

	ip, port, err := net.SplitHostPort(remoteAddress)
	if err != nil {
		return err
	}

	peer := &models.Peer{
		IP:     ip,
		Port:   port,
		Active: true,
	}

	val := p.peerList.Get(id)
	if val != nil {
		temp := &models.Peer{}
		err := json.Unmarshal(val, temp)
		if err == nil {
			peer.Active = temp.Active
			peer.Files = temp.Files
			peer.Size = temp.Size
		}
	}

	data, err = json.Marshal(peer)
	if err != nil {
		return err
	}

	p.peerList.Set(id, data)

	return nil
}

func (p *pongDescriptor) Receive(localAddress string, header *models.DescriptorHeader, buf []byte) error {

	// fmt.Println("pong receive:", header.ID.String())

	// receive pong intended for me
	key := fmt.Sprintf("%s-%s-%s", header.ID.String(), constants.PingDescriptor, "self")
	if !p.descriptorList.Exists(key) {
		// receive pong intended for others
		key = fmt.Sprintf("%s-%s", header.ID.String(), constants.PongDescriptor)
		if !p.descriptorList.Exists(key) {
			return nil
		}
	}

	targetAddress := p.descriptorList.Get(key)
	myAddress := net.JoinHostPort(p.ip, strconv.FormatUint(uint64(p.port), 10))

	pong := p.Decode(buf[constants.DescriptorHeaderLength:])
	pongAddress := net.JoinHostPort(pong.IPAddress, strconv.FormatUint(uint64(pong.Port), 10))

	// deduplication of pong other than self
	key = fmt.Sprintf("%s-%s-%s", header.ID, constants.PongDescriptor, pongAddress)
	if p.descriptorList.Exists(key) {
		return nil
	}

	if pongAddress != myAddress {

		id := helpers.GetHash(pongAddress)
		peer := &models.Peer{
			IP:     pong.IPAddress,
			Port:   strconv.FormatUint(uint64(pong.Port), 10),
			Files:  pong.Files,
			Size:   pong.Size,
			Active: true,
		}

		value := p.peerList.Get(id)
		if value != nil {
			temp := &models.Peer{}
			err := json.Unmarshal(value, temp)
			if err == nil {
				peer.Active = temp.Active
			}
		}

		data, err := json.Marshal(peer)
		if err != nil {
			return err
		}

		p.peerList.Set(id, data)
	}

	if string(targetAddress) != myAddress {

		header.Hops += 1
		header.TTL -= 1

		if header.TTL == 0 {
			return nil
		}

		data, err := p.EncodeDescriptorHeader(header)
		if err != nil {
			return err
		}

		data = append(data, buf[constants.DescriptorHeaderLength:]...)

		// fmt.Println("pong received and forwarded to", header.ID.String(), pongAddress, string(targetAddress))

		_, err = helpers.SendAndReceiveData(string(targetAddress), data, false)
		if err != nil {
			return err
		}
	}

	p.descriptorList.Set(key, []byte{1})

	return nil
}

func (p *pongDescriptor) Encode() []byte {
	buf := [14]byte{}
	ip := p.getEncodedIP(p.ip)
	binary.LittleEndian.PutUint16(buf[0:], p.port)
	copy(buf[2:6], ip)
	binary.LittleEndian.PutUint32(buf[6:], p.files.Count())
	binary.LittleEndian.PutUint32(buf[10:], p.files.Size())

	// fmt.Println("encoded pong", buf[:])
	return buf[:]
}

func (p *pongDescriptor) Decode(buf []byte) *models.Pong {

	pong := &models.Pong{}
	pong.Port = binary.LittleEndian.Uint16(buf[:2])
	pong.IPAddress = p.getDecodedIP(buf[2:6])
	pong.Files = binary.LittleEndian.Uint32(buf[6:10])
	pong.Size = binary.LittleEndian.Uint32(buf[10:])

	// fmt.Println("deoded pong", pong.Port, pong.IPAddress, pong.Files, pong.Size)

	return pong
}
