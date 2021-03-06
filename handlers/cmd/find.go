package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/thamizhv/tgnutella/cache"
	"github.com/thamizhv/tgnutella/handlers/files"
	"github.com/thamizhv/tgnutella/helpers"
	"github.com/thamizhv/tgnutella/models"
)

type findHandler struct {
	findFilesList cache.ServentCache
	files         files.FileHandler
}

func NewFindHandler(serventCacheHelper *cache.ServentCacheHelper, files files.FileHandler) CMDHandler {
	return &findHandler{
		findFilesList: serventCacheHelper.FindFilesList,
		files:         files,
	}
}

func (f *findHandler) Handle(arg string) error {
	id := helpers.GetHash(arg)
	if f.files.Exists(id) {
		fmt.Printf("%s\n", id)
		return nil
	}

	findFile := &models.FindFile{
		Name:   arg,
		Active: true,
	}

	val, err := json.Marshal(findFile)
	if err != nil {
		return err
	}

	f.findFilesList.Set(id, val)

	helpers.FindChannel <- arg

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	select {
	case <-ctx.Done():
		findFile.Active = false
		val, err := json.Marshal(findFile)
		if err == nil {
			f.findFilesList.Set(id, val)
		}

		fmt.Printf("file %s not present locally. Sent the query to network and waited for five seconds to receive a hit. If found, message will be displayed.\n", arg)
	case c := <-helpers.FoundChannel:
		f.findFilesList.Remove(id)
		fmt.Printf("Found file %s with id:%s\n", arg, c)
	}

	return nil
}
