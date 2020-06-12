package cmd

import (
	"fmt"
	"tgnutella/constants"
)

type helpHandler struct{}

func NewHelpHandler() CMDHandler {
	return &helpHandler{}
}

func (h *helpHandler) Handle(_ string) error {
	fmt.Println(constants.HelpText)
	return nil
}
