package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const (
	logPrefix     = "..."
	logEOF        = "EOF"
	logClosed     = "Server closed connection"
	logConnection = "Connected to %s"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	TelnetClient
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (t *telnetClient) Close() error {
	return t.conn.Close()
}

func (t *telnetClient) Send() error {
	return t.handleMessages(t.in, t.conn, logEOF)
}

func (t *telnetClient) Receive() error {
	return t.handleMessages(t.conn, t.out, logClosed)
}

func (t *telnetClient) Connect() error {
	var err error
	if t.conn, err = net.DialTimeout("tcp", t.address, t.timeout); err != nil {
		return err
	}
	return t.log(logConnection, t.address)
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	tc := &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
	return tc
}

func (t *telnetClient) handleMessages(in io.Reader, out io.Writer, eofMsg string) error {
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return t.log("%s", eofMsg)
}

func (t *telnetClient) log(s string, v ...interface{}) error {
	_, err := fmt.Fprintf(os.Stderr, "%s%s\n", logPrefix, fmt.Sprintf(s, v...))
	return err
}
