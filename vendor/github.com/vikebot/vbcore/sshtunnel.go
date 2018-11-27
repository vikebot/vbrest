package vbcore

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

// Endpoint defines the network identification for a party
type Endpoint struct {
	Host string
	Port int
}

// String returns the the endpoint in format Host:Port
func (ep *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", ep.Host, ep.Port)
}

// NewEndpoint creates a new endpoint struct ready to use in SSHTunnels
func NewEndpoint(host string, port int) *Endpoint {
	return &Endpoint{
		Host: host,
		Port: port,
	}
}

// NewEndpointAddr creates a new endpoint struct ready to use in SSHTunnels.
// The format of addr has to be Host:Port where Port is a valid int as described
// by strconv.Atoi
func NewEndpointAddr(addr string) *Endpoint {
	if !strings.Contains(addr, ":") {
		return nil
	}
	split := strings.Split(addr, ":")
	if len(split) != 2 {
		return nil
	}
	port, err := strconv.Atoi(split[1])
	if err != nil {
		return nil
	}
	return NewEndpoint(split[0], port)
}

// SSHTunnel is a wrapper for a secure encrypted tcp tunnel based the
// ssh protocol
type SSHTunnel struct {
	Local      *Endpoint
	Server     *Endpoint
	Remote     *Endpoint
	config     *ssh.ClientConfig
	listener   net.Listener
	errHandler func(error)
}

// NewSSHTunnel initializes a new instance of SSHTunnel
func NewSSHTunnel(user string, password string, local *Endpoint, server *Endpoint, remote *Endpoint, errorHandler func(error)) *SSHTunnel {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	return &SSHTunnel{
		Local:      local,
		Server:     server,
		Remote:     remote,
		config:     sshConfig,
		errHandler: errorHandler,
	}
}

// Bind tries to initialize a new listener on the local address. It doesn't start
// to listen for connections.
func (tunnel *SSHTunnel) Bind() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}

	tunnel.listener = listener
	return nil
}

// Start will accept new clients and forward them into new go routines
func (tunnel *SSHTunnel) Start(tunnelReady chan bool) error {
	if tunnel.listener == nil {
		return errors.New("sshtunnel: listener not ready. call Bind before Start")
	}
	defer tunnel.listener.Close()

	tunnelReady <- true
	for {
		conn, err := tunnel.listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.config)
	if err != nil {
		tunnel.errHandler(fmt.Errorf("sshtunnel: server dial error: %s", err.Error()))
		return
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		tunnel.errHandler(fmt.Errorf("sshtunnel: remote dial error: %s", err.Error()))
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			tunnel.errHandler(fmt.Errorf("sshtunnel: io.Copy error: %s", err.Error()))
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}
