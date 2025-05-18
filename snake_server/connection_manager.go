package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type ConnectionManager struct {
	listener    net.Listener
	addr        string
	quiCh       chan struct{}
	connections map[net.Addr]net.Conn

	g *Game
}

func NewServer(host, port string) (*ConnectionManager, error) {
	addr := net.JoinHostPort(host, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	g, err := NewGame()
	if err != nil {
		return nil, err
	}

	return &ConnectionManager{
		listener:    listener,
		addr:        addr,
		connections: make(map[net.Addr]net.Conn),
		quiCh:       make(chan struct{}),
		g:           g,
	}, nil
}

func (cm *ConnectionManager) StartServer() {
	fmt.Printf("Starting server at addr: %v\n", cm.addr)
	for {
		conn, err := cm.listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go cm.HandleConnection(conn)
	}
}

func (cm *ConnectionManager) HandleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Accepting client connection")
	cm.connections[conn.RemoteAddr()] = conn

	for {
		lengthPrefixBytes := make([]byte, 4)
		_, err := io.ReadFull(conn, lengthPrefixBytes)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				// log.Printf("Client disconnected: %s", conn.RemoteAddr())
				return
			}
			// log.Printf("Error reading length prefix from %s: %v", conn.RemoteAddr(), err)
			return
		}
		cmdReceived := binary.BigEndian.Uint32(lengthPrefixBytes)

		switch cmdReceived {

		case UPcmd:
			fmt.Println(cm.g.s, "UP")
			if cm.g.currDirection != DOWN {
				cm.g.currDirection = UP
			}
		case RIGHTcmd:
			fmt.Println(cm.g.s, "right")
			if cm.g.currDirection != LEFT {
				cm.g.currDirection = RIGHT
			}
		case DOWNcmd:
			fmt.Println(cm.g.s, "down")
			if cm.g.currDirection != UP {
				cm.g.currDirection = DOWN
			}
		case LEFTcmd:
			fmt.Println(cm.g.s, "left")
			if cm.g.currDirection != RIGHT {
				cm.g.currDirection = LEFT
			}

		case QUIT:
			close(cm.g.quiCh)
		case PAUSE:
			cm.g.paused = true
		case RESTART:
			cm.g.resetGame()
		}
	}
}

func (cm *ConnectionManager) sendScore() {
	for {
		select {
		case score := <-cm.g.scoreCh:
			{
				buf := new(bytes.Buffer)
				score32 := int32(score)
				err := binary.Write(buf, binary.BigEndian, score32)
				if err != nil {
					fmt.Println("binary.Write failed:", err)
				}
				b := buf.Bytes()
				for _, conn := range cm.connections {
					n, _ := conn.Write(b)
					debugLogger.Printf("bytes wrote %d\n", n)
				}
			}
		case <-cm.quiCh:
			{
				return
			}
		}

	}
}
