package cmd

import (
	"encoding/json"
	"fmt"
	"tgnutella/cache"
	"tgnutella/models"
)

type closeHandler struct {
	peerList cache.ServentCache
}

func NewCloseHandler(serventCacheHelper *cache.ServentCacheHelper) CMDHandler {
	return &closeHandler{
		peerList: serventCacheHelper.PeerList,
	}
}

func (c *closeHandler) Handle(arg string) error {
	data := c.peerList.Get(arg)
	if data == nil {
		return fmt.Errorf("id not present %s", arg)
	}

	peer := &models.Peer{}

	err := json.Unmarshal(data, peer)
	if err != nil {
		return fmt.Errorf("error in unmarshalling peer data for id %s: %v", arg, err)
	}
	peer.Active = false

	data, err = json.Marshal(peer)
	if err != nil {
		return fmt.Errorf("error in marshalling peer data host %s and port %s: %v", peer.IP, peer.Port, err)
	}

	c.peerList.Set(arg, data)
	return nil
}
