package helpers

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/thamizhv/tgnutella/constants"
)

func SendAndReceiveData(address string, data []byte, receive bool) ([]byte, error) {

	addr, err := net.ResolveTCPAddr(constants.NetworkTypeTCP, address)
	if err != nil {
		return nil, fmt.Errorf("error in resolving TCP address %s: %v\n", address, err)
	}

	conn, err := net.DialTCP(constants.NetworkTypeTCP, nil, addr)
	if err != nil {
		return nil, fmt.Errorf("error in establishing TCP connection %s: %v\n", address, err)
	}

	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("error in writing data in TCP connection %s: %v\n", address, err)
	}

	if receive {
		buf := make([]byte, 2048)

		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}

		return buf[:n], nil
	}

	return nil, nil
}

func SendHTTPGetRequest(url string) ([]byte, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}
