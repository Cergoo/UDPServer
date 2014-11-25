/*
	UDPClient asynchronous udp client
	(c) 2014 Cergoo
	under terms of ISC license

	client := Connect(...) 		- create and run client
	pkg = <-client.ChRead       - get package from server
	client.ChWrite<-pkg         - send package to server
	client = nil                - stop and destroy client
*/
package UDPClient

import (
	"net"
	"runtime"
)

// chanLen Length of read write channel
const chanLen = 10

type (
	// Client main struct
	Client struct {
		conn      net.Conn    //
		frameSize uint16      //
		ChRead    chan []byte //
		ChWrite   chan []byte //
	}
)

// New constructor of a new server
func Connect(addr string, frameSize uint16) (t *Client, err error) {

	if frameSize == 0 {
		frameSize = 65507
	}

	t = &Client{
		ChRead:    make(chan []byte, chanLen),
		ChWrite:   make(chan []byte, chanLen),
		frameSize: frameSize,
	}

	t.conn, err = net.Dial("udp", addr)
	if err != nil {
		return
	}

	go t.reader()
	go t.writer()

	// destroy action
	stopAllGorutines := func(t *Client) {
		close(t.ChWrite)
		t.conn.Close()
	}
	runtime.SetFinalizer(t, stopAllGorutines)
	return
}

func (t *Client) reader() {
	var (
		e error
		n int
	)

	for {
		buf := make([]byte, t.frameSize)
		n, e = t.conn.Read(buf)
		if e == nil {
			t.ChRead <- buf[:n]
		} else {
			return
		}
	}
}

func (t *Client) writer() {
	var (
		v []byte
	)
	for v = range t.ChWrite {
		t.conn.Write(v)
	}
}
