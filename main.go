package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/thamizhv/tgnutella/cache"
	"github.com/thamizhv/tgnutella/constants"
	"github.com/thamizhv/tgnutella/handlers/cmd"
	"github.com/thamizhv/tgnutella/handlers/descriptor"
	"github.com/thamizhv/tgnutella/handlers/files"
	"github.com/thamizhv/tgnutella/servent"
)

func main() {
	serventPort, err := validateAndGetPort(os.Args[1:])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	parsedServentPort, err := strconv.ParseUint(serventPort, 10, 16)
	if err != nil {
		fmt.Printf("Please enter valid servent port number:%s\n", serventPort)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error in getting current working directory:%v\n", err)
		return
	}

	serventAddress := net.JoinHostPort(constants.LocalHost, serventPort)

	httpPort := parsedServentPort + 1
	httpPortString := strconv.FormatUint(httpPort, 10)
	httpAddress := net.JoinHostPort(constants.LocalHost, httpPortString)
	go startHTTPServer(httpPortString, dir)

	fileHelper := files.NewFileHelper(dir)
	serventCacheHelper := cache.NewServentCacheHelper()
	descriptorHelper := descriptor.NewDescriptorHelper(constants.LocalHost, uint16(parsedServentPort),
		uint16(httpPort), serventCacheHelper, fileHelper)

	servent := servent.NewServent(serventAddress, serventCacheHelper)
	addCMDHandlersTo(serventAddress, httpAddress, servent, serventCacheHelper, fileHelper)
	addDescriptorHandlersTo(servent, descriptorHelper)

	err = servent.Start()
	if err != nil {
		fmt.Printf("Error in starting servent :%s\n", err.Error())
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		fmt.Print("\nshutting down...\n")
		os.Exit(0)
	}()

	var command string
	cmdPlaceHolder := serventAddress + ">"
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(cmdPlaceHolder)
		command, _ = reader.ReadString('\n')
		cmd, arg := getCommandAndArgs(command)
		if cmd == "" {
			continue
		}
		err := servent.Handle(cmd, arg)
		if err != nil {
			fmt.Printf("Error in executing command %s, with argument %s :%v\n", cmd, arg, err)
		}
	}

}

func addCMDHandlersTo(serventAddres, httpAddress string, servent servent.Servent, serventCacheHelper *cache.ServentCacheHelper, fileHelper files.FileHandler) {
	servent.AddCmdHandler(constants.CmdTypeHelp, cmd.NewHelpHandler())
	servent.AddCmdHandler(constants.CmdTypeOpen, cmd.NewOpenHandler(serventAddres, httpAddress, serventCacheHelper))
	servent.AddCmdHandler(constants.CmdTypeClose, cmd.NewCloseHandler(serventCacheHelper))
	servent.AddCmdHandler(constants.CmdTypeInfo, cmd.NewInfoHandler(serventCacheHelper))
	servent.AddCmdHandler(constants.CmdTypeFind, cmd.NewFindHandler(serventCacheHelper, fileHelper))
	servent.AddCmdHandler(constants.CmdTypeGet, cmd.NewGetHandler(serventCacheHelper, fileHelper))
}

func addDescriptorHandlersTo(servant servent.Servent, descriptorHelper *descriptor.DescriptorHelper) {
	servant.AddDescriptorHandler(constants.PingDescriptor, descriptor.NewPingDescriptor(descriptorHelper))
	servant.AddDescriptorHandler(constants.PongDescriptor, descriptor.NewPongDescriptor(descriptorHelper))
	servant.AddDescriptorHandler(constants.QueryDescriptor, descriptor.NewQueryDescriptor(descriptorHelper))
	servant.AddDescriptorHandler(constants.QueryHitDescriptor, descriptor.NewQueryHitDescriptor(descriptorHelper))
}

func getCommandAndArgs(command string) (string, string) {
	command = strings.TrimSuffix(command, "\n")

	var cmd, arg string
	cmd = command

	split := strings.Split(command, " ")
	if len(split) > 1 {
		cmd = split[0]
		arg = split[1]
	}

	return cmd, arg
}

func validateAndGetPort(args []string) (string, error) {
	length := len(args)
	if length == 0 {
		return "", fmt.Errorf("Error: missing arguments servent port\n%s", constants.UsageText)
	}

	return args[0], nil
}

func startHTTPServer(port string, directory string) {
	fs := http.FileServer(http.Dir(directory))
	fmt.Printf("http server listening at port %s\n", port)

	err := http.ListenAndServe(":"+port, fs)
	if err != nil {
		fmt.Print("\nerror in invoking http server\n")
		os.Exit(1)
	}
}
