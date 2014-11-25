/*
	UDPServer asynchronous udp server
	(c) 2014 Cergoo
	under terms of ISC license

	srv := UDPServerNew(...) - create and run server
	pkg = <-srv.ChRead       - get package from client
	srv.ChWrite<-pkg         - send package to client
	srv = nil                - stop and destroy server

*/
package UDPServer

import (
	"log"
	"net"
	"runtime"
)

// chanLen Length of read write channel
const chanLen = 10

type (
	// Server main struct
	Server struct {
		conn      *net.UDPConn //
		frameSize uint16       //
		log       *log.Logger  //
		ChRead    chan *TPack  //
		ChWrite   chan *TPack  //
	}
	// TPack pack struct
	TPack struct {
		Pack []byte
		Addr *net.UDPAddr
	}
)

// New constructor of a new server
func New(addr string, frameSize uint16, log *log.Logger) (srv *Server, err error) {
	var (
		laddr *net.UDPAddr
	)

	if frameSize == 0 {
		frameSize = 65507
	}

	srv = &Server{
		ChRead:    make(chan *TPack, chanLen),
		ChWrite:   make(chan *TPack, chanLen),
		frameSize: frameSize,
		log:       log,
	}

	laddr, err = net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return
	}
	srv.conn, err = net.ListenUDP("udp", laddr)
	if err != nil {
		return
	}

	go srv.reader()
	go srv.writer()

	// destroy action
	stopAllGorutines := func(t *Server) {
		close(t.ChWrite)
		t.conn.Close()
	}
	runtime.SetFinalizer(srv, stopAllGorutines)
	return
}

func (t *Server) reader() {
	t.log.Println("run reader")
	defer func() {
		t.log.Println("stop reader")
	}()

	var (
		e error
		n int
	)

	for {
		pack := &TPack{}
		pack.Pack = make([]byte, t.frameSize)
		n, pack.Addr, e = t.conn.ReadFromUDP(pack.Pack)
		if e != nil {
			return
		} else {
			pack.Pack = pack.Pack[:n]
			t.ChRead <- pack
		}
	}
}

func (t *Server) writer() {
	t.log.Println("run writer")
	defer func() {
		t.log.Println("stop writer")
	}()

	var (
		v *TPack
		e error
	)
	for v = range t.ChWrite {
		_, e = t.conn.WriteToUDP(v.Pack, v.Addr)
		if e != nil {
			t.log.Println(e)
		}

	}
}
