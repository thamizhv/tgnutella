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

type queryDescriptor struct {
	*DescriptorHelper
	queryHitDescriptor DescriptorHandler
}

func NewQueryDescriptor(helper *DescriptorHelper) DescriptorHandler {
	return &queryDescriptor{
		DescriptorHelper:   helper,
		queryHitDescriptor: NewQueryHitDescriptor(helper),
	}
}

func (q *queryDescriptor) Send(_ string, _ *models.DescriptorHeader, searchString []byte) error {

	payloadDescriptor, err := helpers.GetPayloadDescriptor(constants.QueryDescriptor)
	if err != nil {
		return err
	}

	buf := q.Encode(0, string(searchString))
	header := &models.DescriptorHeader{
		ID:                uuid.New(),
		PayloadDescriptor: payloadDescriptor,
		TTL:               constants.DefaultServentTTL,
		Hops:              0,
		Length:            uint32(len(buf)),
	}

	data, err := q.EncodeDescriptorHeader(header)
	if err != nil {
		return err
	}

	data = append(data, buf...)

	// fmt.Println("query send:", header.ID.String())

	// save my query with my address
	key := fmt.Sprintf("%s-%s-%s", header.ID.String(), constants.QueryDescriptor, "self")
	val := net.JoinHostPort(q.ip, strconv.FormatUint(uint64(q.port), 10))
	q.descriptorList.Set(key, []byte(val))

	return q.PropagateToPeers("", data)
}

func (q *queryDescriptor) Receive(localAddress string, header *models.DescriptorHeader, buf []byte) error {

	// fmt.Println("query receive:", header.ID.String())
	// fmt.Println("query receive data:", buf)

	key := fmt.Sprintf("%s-%s-%s", header.ID.String(), constants.QueryDescriptor, "self")
	if q.descriptorList.Exists(key) {
		return nil
	}

	key = fmt.Sprintf("%s-%s", header.ID.String(), constants.QueryDescriptor)
	if q.descriptorList.Exists(key) {
		return nil
	}

	query := q.Decode(buf[constants.DescriptorHeaderLength:])
	remoteAddress := query.IPAddress

	if remoteAddress == localAddress {
		return nil
	}

	id := helpers.GetHash(query.Search)
	// fmt.Printf("query search: '%s'\n", query.Search)
	// fmt.Println("id", id)
	// fmt.Println("bytess query search", []byte(query.Search))
	if q.files.Exists(id) {

		err := q.queryHitDescriptor.Send(remoteAddress, header, []byte(id))
		if err != nil {
			return err
		}

		return nil
	}

	q.descriptorList.Set(key, []byte(query.IPAddress))

	// fmt.Println(localAddress)

	query.IPAddress = localAddress
	encoded := q.Encode(query.MinSpeedInKbps, query.Search)

	header.Hops += 1
	header.TTL -= 1
	header.Length = uint32(len(encoded))

	if header.TTL == 0 {
		return nil
	}

	data, err := q.EncodeDescriptorHeader(header)
	if err != nil {
		return err
	}

	data = append(data, encoded...)

	return q.PropagateToPeers(remoteAddress, data)
}

func (q *queryDescriptor) Encode(minSpeed uint16, keyword string) []byte {
	keyword += "\x00"

	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf[0:], minSpeed)
	buf = append(buf, []byte(keyword)...)

	ip := q.getEncodedIP(q.ip)
	buf = append(buf, ip...)

	temp := [2]byte{}
	binary.LittleEndian.PutUint16(temp[0:], q.port)
	buf = append(buf, temp[:]...)

	// fmt.Println("encoded query", buf[:])
	n := len(buf)
	return buf[:n]
}

func (q *queryDescriptor) Decode(buf []byte) *models.Query {

	query := &models.Query{}
	n := len(buf)

	query.MinSpeedInKbps = binary.LittleEndian.Uint16(buf[0:2])
	query.Search = string(buf[2 : n-7])

	ip := q.getDecodedIP(buf[n-6 : n-2])
	port := binary.LittleEndian.Uint16(buf[n-2:])
	query.IPAddress = net.JoinHostPort(ip, strconv.FormatUint(uint64(port), 10))

	// fmt.Println("decoded query", query)
	return query
}
