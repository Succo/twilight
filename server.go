package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type state int

const (
	waiting state = iota
	p1
	p2
	ready
)

type server struct {
	*Map
	name [2]string
}

func (s *server) run() {
	log.Println("Starting tcp server")
	l, err := net.Listen("tcp", ":5555")
	defer l.Close()
	if err != nil {
		panic(err.Error())
	}

	// first player
	con0, err := l.Accept()
	if err != nil {
		panic(err.Error())
	}
	log.Println("First player entered the game")
	ch0 := make(chan bool, 1)
	done0 := make(chan bool, 1)
	go s.runP(con0, 0, ch0, done0)

	// second player
	con1, err := l.Accept()
	if err != nil {
		panic(err.Error())
	}
	log.Println("Second player entered the game")
	ch1 := make(chan bool, 1)
	done1 := make(chan bool, 1)
	go s.runP(con1, 1, ch1, done1)

	ch0 <- true
	<-done0
	ch1 <- true
	<-done1
	for i := 0; i < 50; i++ {
		ch0 <- true
		<-done0
		ch1 <- true
		<-done1
	}
}

func (s *server) runP(c net.Conn, id int, ch chan bool, done chan bool) {
	buf := make([]byte, 10)
	c.Read(buf[:3])
	if bytes.Compare(buf[:3], []byte("NME")) != 0 {
		fmt.Errorf("Invalid first connexion value")
		return
	}

	c.Read(buf[:1])
	t := int(buf[0])
	if t > 10 {
		buf = make([]byte, t)
	}
	c.Read(buf[:t])
	s.name[id] = string(buf[:t])
	s.set(c)
	s.hum(c)
	s.hme(c, id)
	s.state(c, "MAP")
	log.Printf("Initialisation finished for player %d", id)
	<-ch
	s.state(c, "UPD")
	c.SetDeadline(time.Now().Add(10 * time.Second))
	s.upd(c, id)
	done <- true
	for _ = range ch {
		s.state(c, "UPD")
		log.Printf("Move time for player %d", id)
		c.SetDeadline(time.Now().Add(2 * time.Second))
		s.upd(c, id)
		done <- true
	}
}

func (s *server) set(c net.Conn) {
	msg := make([]byte, 5)
	copy(msg, []byte("SET"))
	msg[3] = byte(uint8(s.Rows))
	msg[4] = byte(uint8(s.Columns))
	c.Write(msg)
}

func (s *server) hum(c net.Conn) {
	n := len(s.humans)
	msg := make([]byte, 4+2*n)
	copy(msg, []byte("HUM"))
	msg[3] = byte(uint8(n))
	for i, h := range s.humans {
		hum := s.cells[h]
		msg[4+2*i] = byte(uint8(hum.X))
		msg[4+2*i+1] = byte(uint8(hum.Y))
	}
	c.Write(msg)
}

func (s *server) hme(c net.Conn, id int) {
	msg := make([]byte, 5)
	copy(msg, []byte("HME"))
	mon := s.cells[s.monster[id][0]]
	msg[3] = byte(uint8(mon.X))
	msg[4] = byte(uint8(mon.Y))
	c.Write(msg)
}

func (s *server) state(c net.Conn, trame string) {
	n := len(s.humans) + len(s.monster[0]) + len(s.monster[1])
	msg := make([]byte, 4+5*n)
	copy(msg, []byte(trame))
	msg[3] = byte(uint8(n))
	var i int
	// Send all humans data
	for _, h := range s.humans {
		hum := s.cells[h]
		msg[4+5*i] = byte(uint8(hum.X))
		msg[4+5*i+1] = byte(uint8(hum.Y))
		msg[4+5*i+2] = byte(uint8(hum.Count))
		i++
	}
	// Send all vamp data
	for _, m := range s.monster[1] {
		mon := s.cells[m]
		msg[4+5*i] = byte(uint8(mon.X))
		msg[4+5*i+1] = byte(uint8(mon.Y))
		msg[4+5*i+3] = byte(uint8(mon.Count))
		i++
	}
	// Send all wolf data
	for _, m := range s.monster[0] {
		mon := s.cells[m]
		msg[4+5*i] = byte(uint8(mon.X))
		msg[4+5*i+1] = byte(uint8(mon.Y))
		msg[4+5*i+4] = byte(uint8(mon.Count))
		i++
	}
	c.Write(msg)
}

func (s *server) upd(c net.Conn, id int) error {
	buf := make([]byte, 100)
	c.Read(buf[:3])
	if bytes.Compare(buf[:3], []byte("MOV")) != 0 {
		return errors.New("Invalid MOV trame value")
	}

	c.Read(buf[:1])
	t := int(buf[0])
	moves := make([]move, t)
	for i := 0; i < t; i++ {
		c.Read(buf[:5])
		moves[i] = move{
			oldx:  int(buf[0]),
			oldy:  int(buf[1]),
			count: int(buf[2]),
			newx:  int(buf[3]),
			newy:  int(buf[4]),
		}
	}
	fmt.Println(t, moves)
	return s.apply(moves, id)
}
