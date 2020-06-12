package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/thamizhv/tgnutella/cache"
	"github.com/thamizhv/tgnutella/handlers/files"
	"github.com/thamizhv/tgnutella/helpers"
	"github.com/thamizhv/tgnutella/models"
)

type getHandler struct {
	peerList      cache.ServentCache
	peerFilesList cache.ServentCache
	files         files.FileHandler
}

func NewGetHandler(serventCacheHelper *cache.ServentCacheHelper, files files.FileHandler) CMDHandler {
	return &getHandler{
		peerList:      serventCacheHelper.PeerList,
		peerFilesList: serventCacheHelper.PeerFilesList,
		files:         files,
	}
}

func (g *getHandler) Handle(arg string) error {
	if g.files.Exists(arg) {
		fmt.Printf("file already present in current working directory: %s\n", arg)
		return nil
	}

	if !g.peerFilesList.Exists(arg) {
		fmt.Printf("file not found in peers: %s\n", arg)
		return nil
	}

	val := g.peerFilesList.Get(arg)

	peerFileList := make([]*models.PeerFile, 0)

	err := json.Unmarshal(val, &peerFileList)
	if err != nil {
		return fmt.Errorf("error in unmarshalling peer files list: %s", arg)
	}

	for _, peerFile := range peerFileList {
		httpAddress := net.JoinHostPort(peerFile.IPAddress, strconv.FormatUint(uint64(peerFile.HTTPPort), 10))
		tcpAddress := net.JoinHostPort(peerFile.IPAddress, strconv.FormatUint(uint64(peerFile.Port), 10))

		id := helpers.GetHash(tcpAddress)
		val := g.peerList.Get(id)
		if val != nil {
			peer := &models.Peer{}
			err = json.Unmarshal(val, peer)
			if err != nil {
				fmt.Println("get command: error occurred in unmarshalling peerlist")
				continue
			}

			if !peer.Active {
				fmt.Printf("get command: peer %s is closed. Open connection to download\n", peer.IP+":"+peer.Port)
				continue
			}
		}

		bodyBytes, err := helpers.SendHTTPGetRequest("http://" + httpAddress + "/" + peerFile.Details.FileName)
		if err != nil {
			fmt.Println("get command: error occurred in sending http request " + httpAddress + "/" + peerFile.Details.FileName)
			continue
		}

		f, err := os.Create(peerFile.Details.FileName)
		if err != nil {
			fmt.Println("get command: error occurred in creating file " + peerFile.Details.FileName)
			continue
		}

		_, err = f.Write(bodyBytes)
		if err != nil {
			f.Close()
			fmt.Println("get command: error occurred in writing file ")
			continue
		}

		f.Close()

		g.files.UpdateFileList()
		fmt.Printf("Downloaded id %s in current working directory. File Name %s\n", arg, peerFile.Details.FileName)
		break
	}

	return nil
}
