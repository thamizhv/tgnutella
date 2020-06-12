package servent

import (
	"fmt"
	"net"
	"time"

	"github.com/thamizhv/tgnutella/cache"
	"github.com/thamizhv/tgnutella/constants"
	"github.com/thamizhv/tgnutella/handlers/cmd"
	"github.com/thamizhv/tgnutella/handlers/descriptor"
	"github.com/thamizhv/tgnutella/helpers"
)

type Servent interface {
	Start() error
	Handle(cmd, arg string) error
	AddCmdHandler(cmd string, handler cmd.CMDHandler)
	AddDescriptorHandler(descriptor string, handler descriptor.DescriptorHandler)
}

type servent struct {
	address            string
	descriptorHandlers map[string]descriptor.DescriptorHandler
	cmdHandlers        map[string]cmd.CMDHandler
}

func NewServent(address string, serventCacheHelper *cache.ServentCacheHelper) Servent {
	return &servent{
		address:            address,
		cmdHandlers:        make(map[string]cmd.CMDHandler),
		descriptorHandlers: make(map[string]descriptor.DescriptorHandler),
	}
}

func (s *servent) AddCmdHandler(cmd string, handler cmd.CMDHandler) {
	s.cmdHandlers[cmd] = handler
}

func (s *servent) AddDescriptorHandler(descriptor string, handler descriptor.DescriptorHandler) {
	s.descriptorHandlers[descriptor] = handler
}

func (s *servent) Handle(cmd, arg string) error {
	handler, ok := s.cmdHandlers[cmd]
	if !ok {
		fmt.Println(constants.CommandNotFound)
		fmt.Println(constants.HelpText)
		return nil
	}

	return handler.Handle(arg)
}

func (s *servent) Start() error {
	address, err := net.ResolveTCPAddr(constants.NetworkTypeTCP, s.address)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP(constants.NetworkTypeTCP, address)
	if err != nil {
		return err
	}

	fmt.Printf("servent listening at %s\n", listener.Addr().String())

	go s.listen(listener)

	go s.pingPeers()

	go s.listenFindChannel()

	return nil
}

func (s *servent) listen(listener *net.TCPListener) {
	defer listener.Close()
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("error in reading tcp connection: %v\n", err)
			continue
		}

		if n < constants.DescriptorHeaderLength {

			if n == len(constants.ConnectionRequest) {
				err := s.handleConnectionRequest(conn, buf[:n])
				if err != nil {
					fmt.Printf("error in sending connection response: %v\n", err)
				}
			} else {
				fmt.Printf("invalid data from tcp connection: %v\n", err)
			}

			continue
		}

		descriptorHeader, err := descriptor.DecodeDescriptorHeader(buf[:constants.DescriptorHeaderLength])
		if err != nil {
			fmt.Printf("error in decoding descriptor header: %v\n", err)
			continue
		}

		descriptor, err := helpers.GetDescriptorString(descriptorHeader.PayloadDescriptor)
		if err != nil {
			fmt.Printf("error in getting descriptor string: %v\n", err)
			continue
		}

		handler := s.descriptorHandlers[descriptor]
		err = handler.Receive(conn.LocalAddr().String(), descriptorHeader, buf[:n])
		if err != nil {
			fmt.Printf("error in handling tcp connection: %v\n", err)
		}
	}
}

func (s *servent) handleConnectionRequest(conn *net.TCPConn, data []byte) error {
	if string(data) != constants.ConnectionRequest {
		return fmt.Errorf("invalid version or connection request: %s", string(data))
	}

	_, err := conn.Write([]byte(constants.ConnectionResponse))
	if err != nil {
		return err
	}

	return nil
}

func (s *servent) pingPeers() {
	for {
		pingHandler := s.descriptorHandlers[constants.PingDescriptor]
		err := pingHandler.Send("", nil, nil)
		if err != nil {
			fmt.Printf("error in sending ping descriptor payload: %v\n", err)
		}

		time.Sleep(constants.PingInterval)
	}
}

func (s *servent) listenFindChannel() {
	queryHandler := s.descriptorHandlers[constants.QueryDescriptor]
	for c := range helpers.FindChannel {
		err := queryHandler.Send("", nil, []byte(c))
		if err != nil {
			fmt.Printf("error in sending query descriptor payload: %v\n", err)
		}
	}
}
