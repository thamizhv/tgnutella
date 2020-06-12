package descriptor

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
	"strings"

	"github.com/thamizhv/tgnutella/cache"
	"github.com/thamizhv/tgnutella/handlers/files"
	"github.com/thamizhv/tgnutella/helpers"
	"github.com/thamizhv/tgnutella/models"
)

type DescriptorHandler interface {
	Send(destinationAddress string, header *models.DescriptorHeader, data []byte) error
	Receive(localAddress string, descriptorHeader *models.DescriptorHeader, data []byte) error
	PropagateToPeers(destinationAddress string, data []byte) error
}

type DescriptorHelper struct {
	ip             string
	port           uint16
	httpPort       uint16
	peerList       cache.ServentCache
	descriptorList cache.ServentCache
	peerFilesList  cache.ServentCache
	findFilesList  cache.ServentCache
	files          files.FileHandler
}

func NewDescriptorHelper(ip string, port, httpPort uint16, serventCacheHelper *cache.ServentCacheHelper, files files.FileHandler) *DescriptorHelper {
	return &DescriptorHelper{
		ip:             ip,
		port:           port,
		httpPort:       httpPort,
		files:          files,
		peerList:       serventCacheHelper.PeerList,
		descriptorList: serventCacheHelper.DescriptorList,
		peerFilesList:  serventCacheHelper.PeerFilesList,
		findFilesList:  serventCacheHelper.FindFilesList,
	}
}

func (d *DescriptorHelper) EncodeDescriptorHeader(header *models.DescriptorHeader) ([]byte, error) {

	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodeDescriptorHeader(data []byte) (*models.DescriptorHeader, error) {
	descriptorHeader := &models.DescriptorHeader{}

	err := binary.Read(bytes.NewReader(data), binary.LittleEndian, descriptorHeader)
	if err != nil {
		return nil, err
	}

	return descriptorHeader, nil
}

func (d *DescriptorHelper) PropagateToPeers(destinationAddress string, data []byte) error {
	peers := d.peerList.GetAll()

	for k, v := range peers {

		peer := &models.Peer{}

		err := json.Unmarshal(v, peer)
		if err != nil {
			continue
		}

		if !peer.Active {
			continue
		}

		peerAddress := net.JoinHostPort(peer.IP, peer.Port)

		if peerAddress == destinationAddress {
			continue
		}

		// fmt.Println("propagated peer address", peerAddress)
		// fmt.Println("propagated peer data", data)

		_, err = helpers.SendAndReceiveData(peerAddress, data, false)
		if err != nil {
			if strings.Contains(err.Error(), "error in establishing TCP connection") {
				peer.Active = false
				val, err := json.Marshal(peer)
				if err == nil {
					d.peerList.Set(k, val)
				}
			}
		}
	}

	return nil
}

func (d *DescriptorHelper) getEncodedIP(ipAddress string) []byte {

	buf := make([]byte, 4)
	ip := net.ParseIP(ipAddress)

	ipUint := binary.BigEndian.Uint32(ip.To4())
	binary.BigEndian.PutUint32(buf, ipUint)

	return buf
}

func (d *DescriptorHelper) getDecodedIP(buf []byte) string {
	ip := make(net.IP, 4)
	ipUint := binary.BigEndian.Uint32(buf)

	binary.BigEndian.PutUint32(ip, ipUint)
	return ip.String()
}
