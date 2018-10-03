package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
)

type guessClient struct {
	Host string
	conn net.Conn
}

func (client *guessClient) connect() (err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", client.Host)
	if err != nil {
		return
	}

	client.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}
	defer client.conn.Close()

	err = client.startSession()
	if err != nil {
		return
	}
	return
}

func (client guessClient) startSession() (err error) {
	var i int
	for {
		fmt.Print("guess ")
		_, err = fmt.Scanf("%d", &i)
		if err != nil {
			return
		}
		client.conn.Write([]byte("guess " + strconv.Itoa(i)))

		reply := make([]byte, 1024)
		var n int
		n, err = client.conn.Read(reply)
		if err != nil {
			return
		}
		replyStr := string(reply[:n])
		fmt.Println(replyStr)
		if replyStr == "correct" {
			return
		} else if !(replyStr == "less" || replyStr == "more") {
			err = errors.New("unexpected response")
			return
		}
	}
}

func main() {
	var host string
	if len(os.Args) >= 2 {
		host = os.Args[1]
	}

	client := guessClient{
		Host: host,
	}
	err := client.connect()
	if err != nil {
		println(err.Error())
	}
}
