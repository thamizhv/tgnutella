package cmd

import (
	"fmt"

	"github.com/thamizhv/tgnutella/constants"
)

type helpHandler struct{}

func NewHelpHandler() CMDHandler {
	return &helpHandler{}
}

func (h *helpHandler) Handle(_ string) error {
	fmt.Println(constants.HelpText)
	return nil
}
