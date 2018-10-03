package main

import (
	"bytes"
	"math/rand"
	"net"
	"os"
	"strconv"
)

func main() {
	var host string
	if len(os.Args) >= 2 {
		host = os.Args[1]
	}
	serv := guessServer{
		Host: host,
	}
	err := serv.listenAndServe()
	if err != nil {
		println(err.Error())
		return
	}
}

type guessServer struct {
	Host string
}

func (serv guessServer) listenAndServe() (err error) {
	if serv.Host == "" {
		serv.Host = ":8080"
	}
	l, err := net.Listen("tcp", serv.Host)
	if err != nil {
		return
	}
	defer l.Close()

	for {
		var conn net.Conn
		conn, err = l.Accept()
		if err != nil {
			return
		}

		go serv.handleRequest(conn)
	}

}

// Handles incoming requests
func (serv guessServer) handleRequest(conn net.Conn) {
	defer conn.Close()

	number := rand.Int()%90 + 10
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			println(err.Error())
			return
		}
		if !bytes.Equal(buf[:6], []byte("guess ")) {
			println(err.Error())
			return
		}
		guessNumber, err := strconv.Atoi(string(buf[6:n]))
		if err != nil {
			println(err.Error())
			return
		}

		if guessNumber < number {
			conn.Write([]byte("more"))
		} else if guessNumber > number {
			conn.Write([]byte("less"))
		} else {
			conn.Write([]byte("correct"))
			return
		}
	}
}
