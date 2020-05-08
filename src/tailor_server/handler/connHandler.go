package handler

import (
	"fmt"
	"net"
	"protocol"
	"tailor"
)

const (
	setex byte = iota
	setnx
	set
	get
	del
	unlink
	incr
	incrby
	ttl
	keys
	cnt
	save
	load
)

func HandleConn(conn net.Conn, cache *tailor.Cache, savingPath string, maxSizeOfDatagram int) {
	defer func() {
		kvs, _ := cache.Keys("[A-z]+")
		for i := range kvs {
			fmt.Printf("key: %s, val: %v\n", kvs[i].Key(), kvs[i].Val())
		}
	}()
	for {
		datagram, err := readDatagram(conn, maxSizeOfDatagram)
		if err != nil {
			break
		}
		switch datagram.Op {
		case setex:
			doSetex(cache, datagram, conn)
		case setnx:
			doSetnx(cache, datagram, conn)
		case set:
			doSet(cache, datagram, conn)
		case get:
			doGet(cache, datagram, conn)
		case del:
			doDel(cache, datagram, conn)
		case unlink:
			doUnlink(cache, datagram, conn)
		case incr:
			doIncr(cache, datagram, conn)
		case incrby:
			doIncrby(cache, datagram, conn)
		case ttl:
			doTtl(cache, datagram, conn)
		case cnt:
			doCnt(cache, conn)
		case save:
			doSave(savingPath, cache, conn)
		case load:
			doLoad(savingPath, cache, conn)
		}
	}
}

func readDatagram(conn net.Conn, maxSize int) (*protocol.Protocol, error) {
	buf := make([]byte, maxSize)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	datagram, err := protocol.GetDatagram(buf[:n])
	if err != nil {
		return nil, err
	}
	return datagram, nil
}
