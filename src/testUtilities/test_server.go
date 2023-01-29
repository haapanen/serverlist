package testutilities

import (
	"net"
	"os"
	"strings"
	"time"
)

type TestServer struct {
	getServersResponse func() []byte
	getStatusResponse  func() []byte
	messageChan        chan string
}

func getDefaultServersResponse() []byte {
	buf, err := os.ReadFile("./testUtilities/getserversResponse.bin")
	if err != nil {
		panic(err)
	}
	return buf
}

func getDefaultStatusResponse() []byte {

	buf, err := os.ReadFile("./testUtilities/getstatusResponse.bin")
	if err != nil {
		panic(err)
	}
	return buf
}

func NewTestServer() *TestServer {
	server := TestServer{
		getServersResponse: getDefaultServersResponse,
		getStatusResponse:  getDefaultStatusResponse,
		messageChan:        make(chan string),
	}

	return &server
}

func (t *TestServer) Start(port int) {
	go func() {
		conn, err := net.ListenUDP("udp", &net.UDPAddr{
			IP:   net.IPv4zero,
			Port: port,
		})

		if err != nil {
			panic(err)
		}

		defer conn.Close()

		for {
			select {
			case msg := <-t.messageChan:
				if msg == "stop" {
					return
				}
			default:
			}

			buf := make([]byte, 1024)
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, addr, err := conn.ReadFrom(buf)
			if err, ok := err.(net.Error); ok && err.Timeout() {
				continue
			} else if err != nil {
				panic(err)
			}

			if strings.HasPrefix(string(buf[:n]), "\xff\xff\xff\xffgetstatus") {
				go t.handleGetStatus(conn, addr)
			} else if strings.HasPrefix(string(buf[:n]), "\xff\xff\xff\xffgetservers") {
				go t.handleGetServers(conn, addr)
			}
		}
	}()

}

func (t *TestServer) handleGetServers(conn net.PacketConn, addr net.Addr) {
	conn.WriteTo(t.getServersResponse(), addr)
}

func (t *TestServer) handleGetStatus(conn net.PacketConn, addr net.Addr) {
	conn.WriteTo(t.getStatusResponse(), addr)
}

func (t *TestServer) Stop() {
	t.messageChan <- "stop"
}

func (t *TestServer) ResetGetServersResponse() {
	t.getServersResponse = getDefaultServersResponse
}

func (t *TestServer) ResetGetStatusResponse() {
	t.getStatusResponse = getDefaultStatusResponse
}

func (t *TestServer) SetGetServersResponse(response []byte) {
	t.getServersResponse = func() []byte {
		return response
	}
}

func (t *TestServer) SetGetStatusResponse(response []byte) {
	t.getStatusResponse = func() []byte {
		return response
	}
}
