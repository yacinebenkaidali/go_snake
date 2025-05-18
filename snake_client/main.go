package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/gdamore/tcell/v2"
)

const (
	// snake controls
	UP    = 0x000
	RIGHT = 0x001
	LEFT  = 0x010
	DOWN  = 0x011

	// game controls
	QUIT    = 0x100
	PAUSE   = 0x101
	RESTART = 0x110
)

var debugLogger *log.Logger

func readConn(conn net.Conn, s *tcell.Screen) {
	for {
		buff := make([]byte, 4)
		_, err := io.ReadFull(conn, buff)
		if err != nil {
			continue
		}
		res := binary.BigEndian.Uint32(buff)
		debugLogger.Printf("received %d\n", res)

		printAt(0, 0, fmt.Sprint(res), s)
	}
}

func printAt(x, y int, s string, screen *tcell.Screen) {
	for i, ch := range s {
		(*screen).SetContent(x+i, y, ch, nil, tcell.StyleDefault)
	}
	(*screen).Show()
}

func main() {
	logFile, err := os.OpenFile("debug2.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log file: %v\n", err)
		os.Exit(1)
	}

	// Initialize the logger
	debugLogger = log.New(logFile, "DEBUG: ", log.Ltime|log.Lshortfile)

	serverAddr := ":3000"
	// Connect to the server
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer conn.Close()

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset))
	s.Clear()
	defer s.Fini()

	go readConn(conn, &s)

	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			{
				s.Sync()
			}
		case *tcell.EventKey:
			{
				buff := make([]byte, 4)
				send := false
				quit := false

				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
					binary.BigEndian.PutUint32(buff, QUIT)
					send = true
					quit = true
				} else if ev.Key() == tcell.KeyCtrlL {
					s.Sync()
				} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
					s.Clear()
				} else if ev.Rune() == 'R' || ev.Rune() == 'r' {
					binary.BigEndian.PutUint32(buff, RESTART)
					send = true
				} else if ev.Rune() == 'P' || ev.Rune() == 'p' {
					binary.BigEndian.PutUint32(buff, PAUSE)
					send = true
				}

				if !send {
					switch ev.Key() {
					case tcell.KeyUp:
						binary.BigEndian.PutUint32(buff, UP)
						send = true
					case tcell.KeyRight:
						binary.BigEndian.PutUint32(buff, RIGHT)
						send = true
					case tcell.KeyDown:
						binary.BigEndian.PutUint32(buff, DOWN)
						send = true
					case tcell.KeyLeft:
						binary.BigEndian.PutUint32(buff, LEFT)
						send = true
					}
				}

				if !send {
					switch ev.Rune() {
					case 'w', 'W':
						// up
						binary.BigEndian.PutUint32(buff, UP)
						send = true
					case 'd', 'D':
						// right
						binary.BigEndian.PutUint32(buff, RIGHT)
						send = true
					case 's', 'S':
						// down
						binary.BigEndian.PutUint32(buff, DOWN)
						send = true
					case 'a', 'A':
						// left
						binary.BigEndian.PutUint32(buff, LEFT)
						send = true
					}
				}

				if send {
					conn.Write(buff)
				}
				if quit {
					return
				}

			}
		}
	}
}
