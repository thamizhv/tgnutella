package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/thamizhv/tgnutella/cache"
	"github.com/thamizhv/tgnutella/constants"
	"github.com/thamizhv/tgnutella/helpers"
	"github.com/thamizhv/tgnutella/models"
)

type openHandler struct {
	myTCPAddress  string
	myHTTPAddress string
	peerList      cache.ServentCache
}

func NewOpenHandler(myTCPAddress, myHTTPAddress string, serventCacheHelper *cache.ServentCacheHelper) CMDHandler {
	return &openHandler{
		myTCPAddress:  myTCPAddress,
		myHTTPAddress: myHTTPAddress,
		peerList:      serventCacheHelper.PeerList,
	}
}

func (o *openHandler) Handle(arg string) error {
	connectionString, err := o.validateArg(arg)
	if err != nil {
		return err
	}

	if connectionString == o.myHTTPAddress || connectionString == o.myTCPAddress {
		return fmt.Errorf("couldn't connect to ports in which current server runs: %s", arg)
	}

	data, err := helpers.SendAndReceiveData(connectionString, []byte(constants.ConnectionRequest), true)
	if err != nil {
		return err
	}

	if string(data) != constants.ConnectionResponse {
		return fmt.Errorf("invalid connection response '%s'", string(data))
	}

	ip, port, err := net.SplitHostPort(connectionString)
	if err != nil {
		return fmt.Errorf("error in splitting destination host and port %s: %v", arg, err)
	}

	id := helpers.GetHash(connectionString)
	peer := &models.Peer{}

	value := o.peerList.Get(id)
	if value == nil {
		peer = &models.Peer{
			IP:     ip,
			Port:   port,
			Active: true,
		}
	} else {
		err = json.Unmarshal(value, peer)
		if err != nil {
			return fmt.Errorf("error in unmarshalling peer data for connection string %s: %v", connectionString, err)
		}
		peer.Active = true
	}

	data, err = json.Marshal(peer)
	if err != nil {
		return fmt.Errorf("error in marshalling peer data host %s and port %s: %v", peer.IP, peer.Port, err)
	}

	o.peerList.Set(id, data)

	fmt.Printf("%s\n", id)

	return nil
}

func (o *openHandler) validateArg(arg string) (string, error) {

	split := strings.Split(arg, ":")
	if len(split) == 0 {
		return "", fmt.Errorf("invalid connection address %s", arg)
	}

	if split[0] == "" {
		return constants.LocalHost + arg, nil
	}

	return arg, nil
}
