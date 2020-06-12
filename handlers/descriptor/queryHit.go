package descriptor

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"tgnutella/constants"
	"tgnutella/helpers"
	"tgnutella/models"
)

type queryHitDescriptor struct {
	*DescriptorHelper
}

func NewQueryHitDescriptor(helper *DescriptorHelper) DescriptorHandler {
	return &queryHitDescriptor{
		DescriptorHelper: helper,
	}
}

func (qh *queryHitDescriptor) Send(remoteAddress string, previous *models.DescriptorHeader, id []byte) error {

	// fmt.Println("queryHit send:", previous.ID.String())

	// my queryHit for reference
	key := fmt.Sprintf("%s-%s", previous.ID.String(), constants.QueryHitDescriptor)
	if qh.descriptorList.Exists(key) {
		return nil
	}

	payloadDescriptor, err := helpers.GetPayloadDescriptor(constants.QueryHitDescriptor)
	if err != nil {
		return err
	}

	buf := qh.Encode(string(id))

	header := &models.DescriptorHeader{
		ID:                previous.ID,
		PayloadDescriptor: payloadDescriptor,
		TTL:               previous.Hops + 1,
		Hops:              0,
		Length:            uint32(len(buf)),
	}

	data, err := qh.EncodeDescriptorHeader(header)
	if err != nil {
		return err
	}

	data = append(data, buf...)

	_, err = helpers.SendAndReceiveData(remoteAddress, data, false)
	if err != nil {
		return err
	}

	qh.descriptorList.Set(key, []byte(remoteAddress))
	return nil
}

func (qh *queryHitDescriptor) Receive(localAddress string, header *models.DescriptorHeader, buf []byte) error {

	// fmt.Println("queryHit receive:", header.ID.String())

	// receive queryhit for my query
	key := fmt.Sprintf("%s-%s-%s", header.ID.String(), constants.QueryDescriptor, "self")
	if !qh.descriptorList.Exists(key) {
		key = fmt.Sprintf("%s-%s", header.ID.String(), constants.QueryDescriptor)
		if !qh.descriptorList.Exists(key) {
			return nil
		}
	}

	targetAddress := qh.descriptorList.Get(key)
	myAddress := net.JoinHostPort(qh.ip, strconv.FormatUint(uint64(qh.port), 10))

	queryHit := qh.Decode(buf[constants.DescriptorHeaderLength:])
	queryHitAddress := net.JoinHostPort(queryHit.IPAddress, strconv.FormatUint(uint64(queryHit.Port), 10))

	// deduplication of query hit other than self
	key = fmt.Sprintf("%s-%s-%s", header.ID, constants.QueryHitDescriptor, queryHitAddress)
	if qh.descriptorList.Exists(key) {
		return nil
	}

	if queryHitAddress != myAddress {

		id := queryHit.ResultSet.FileIndex
		peerFile := &models.PeerFile{
			IPAddress: queryHit.IPAddress,
			Port:      queryHit.Port,
			HTTPPort:  queryHit.HTTPPort,
			Details: models.ResultSet{
				FileIndex: queryHit.ResultSet.FileIndex,
				FileSize:  queryHit.ResultSet.FileSize,
				FileName:  queryHit.ResultSet.FileName,
			},
		}

		peerFileList := make([]*models.PeerFile, 0)
		peerFileList = append(peerFileList, peerFile)
		value := qh.peerFilesList.Get(id)
		if value != nil {
			temp := make([]*models.PeerFile, 0)
			err := json.Unmarshal(value, &temp)
			if err == nil {
				peerFileList = append(peerFileList, temp...)
			}
		}

		data, err := json.Marshal(peerFileList)
		if err != nil {
			return err
		}

		qh.peerFilesList.Set(id, data)

		val := qh.findFilesList.Get(id)
		if val != nil {
			findFile := &models.FindFile{}
			err := json.Unmarshal(val, findFile)
			if err == nil {
				if findFile.Active {
					helpers.FoundChannel <- id
				} else {
					fmt.Printf("\nFound file %s with id:%s\n", queryHit.ResultSet.FileName, id)
				}
			}
		}
	}

	if string(targetAddress) != myAddress {

		header.Hops += 1
		header.TTL -= 1

		if header.TTL == 0 {
			return nil
		}

		data, err := qh.EncodeDescriptorHeader(header)
		if err != nil {
			return err
		}

		data = append(data, buf[constants.DescriptorHeaderLength:]...)

		// fmt.Println("query hit received and forwarded to", header.ID.String(), queryHitAddress, string(targetAddress))

		_, err = helpers.SendAndReceiveData(string(targetAddress), data, false)
		if err != nil {
			return err
		}
	}

	qh.descriptorList.Set(key, []byte{1})

	return nil
}

func (qh *queryHitDescriptor) Encode(id string) []byte {
	buf := make([]byte, 13)
	buf[0] = 1
	binary.LittleEndian.PutUint16(buf[1:], qh.httpPort)
	binary.LittleEndian.PutUint16(buf[3:5], qh.port)
	ip := qh.getEncodedIP(qh.ip)
	copy(buf[5:9], ip)
	binary.LittleEndian.PutUint32(buf[9:], 1000)

	file := qh.files.Get(id)
	fileSize := uint32(file.Size)
	fileName := file.Name + "\x00"

	resultSet := []byte{}
	resultSet = append(resultSet, []byte(id)...)
	temp := make([]byte, 4)
	binary.LittleEndian.PutUint32(temp, fileSize)
	resultSet = append(resultSet, temp...)
	resultSet = append(resultSet, []byte(fileName)...)

	buf = append(buf, resultSet...)

	// fmt.Println("encoded queryHit", buf[:])
	return buf[:]
}

func (qh *queryHitDescriptor) Decode(buf []byte) *models.QueryHit {
	n := len(buf)
	queryHit := &models.QueryHit{}
	queryHit.NumHits = uint8(buf[0])
	queryHit.HTTPPort = binary.LittleEndian.Uint16(buf[1:])
	queryHit.Port = binary.LittleEndian.Uint16(buf[3:5])
	queryHit.IPAddress = qh.getDecodedIP(buf[5:9])
	queryHit.Speed = binary.LittleEndian.Uint32(buf[9:13])

	resultSet := models.ResultSet{}

	// id length = 32 so.. 13+32
	resultSet.FileIndex = string(buf[13:45])
	resultSet.FileSize = binary.LittleEndian.Uint32(buf[45:49])
	resultSet.FileName = string(buf[49 : n-1])
	queryHit.ResultSet = resultSet

	// fmt.Println("deoded queryHit", queryHit)

	return queryHit
}
