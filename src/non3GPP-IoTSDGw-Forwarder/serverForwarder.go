package main

import (
	"bytes"
	"net"
	"time"
)

type connn struct {
	net.Conn
	readBuf *bytes.Buffer

	lastSendEnd time.Time
	delay time.Duration
}

var now = time.Now
var sleep = time.Sleep