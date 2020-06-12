package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/thamizhv/tgnutella/cache"
	"github.com/thamizhv/tgnutella/constants"
	"github.com/thamizhv/tgnutella/models"
)

type infoHandler struct {
	peerList cache.ServentCache
}

func NewInfoHandler(serventCacheHelper *cache.ServentCacheHelper) CMDHandler {
	return &infoHandler{
		peerList: serventCacheHelper.PeerList,
	}
}

func (i *infoHandler) Handle(arg string) error {
	if arg != constants.ArgTypeConnections {
		return fmt.Errorf("invalid argument type for info command: %s", arg)
	}

	out := i.peerList.GetAll()

	fmt.Printf("id\taddress\tTotalFiles\tTotalSize\n")
	for k, v := range out {
		peer := &models.Peer{}

		err := json.Unmarshal(v, peer)
		if err != nil {
			continue
		}

		if !peer.Active {
			continue
		}

		fmt.Printf("%s\t%s\t%s\t%s\n", k, net.JoinHostPort(peer.IP, peer.Port), strconv.FormatUint(uint64(peer.Files), 10),
			strconv.FormatUint(uint64(peer.Size), 10))
	}

	return nil
}
