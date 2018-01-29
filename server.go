package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
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
	ch0 := make(chan []cell, 1)
	done0 := make(chan []cell, 1)
	go s.runP(con0, 0, ch0, done0)

	// second player
	con1, err := l.Accept()
	if err != nil {
		panic(err.Error())
	}
	log.Println("Second player entered the game")
	ch1 := make(chan []cell, 1)
	done1 := make(chan []cell, 1)
	go s.runP(con1, 1, ch1, done1)

	s.state = ready
	var update0, update1 []cell
	ch0 <- make([]cell, 0)
	update0 = <-done0
	ch1 <- update0
	update1 = <-done1
	// Play for 50 rounds
	for i := 0; i < 50; i++ {
		ch0 <- append(update0, update1...)
		update0 = <-done0
		if s.state > ready {
			break
		}
		ch1 <- append(update1, update0...)
		update1 = <-done1
		if s.state > ready {
			break
		}
	}
	close(ch0)
	close(ch1)
	switch s.state {
	case win0:
		log.Println("Player 0 won")
	case win1:
		log.Println("Player 0 won")
	case null:
		log.Println("Equality")
	}
}

func (s *server) runP(c net.Conn, id int, ch chan []cell, done chan []cell) {
	defer s.bye(c)
	// Initialisation stuff
	reader := bufio.NewReader(c)
	buf := make([]byte, 10)
	io.ReadFull(reader, buf[:4])
	if bytes.Compare(buf[:3], []byte("NME")) != 0 {
		panic("Invalid first connexion value")
		return
	}

	t := int(buf[3])
	if t > 10 {
		buf = make([]byte, t)
	}
	io.ReadFull(reader, buf[:t])
	s.name[id] = string(buf[:t])
	s.set(c)
	s.hum(c)
	s.hme(c, id)
	s.send_map(c)
	log.Printf("Initialisation finished for player %d", id)
	s.send_upd(c, make([]cell, 0))

	// First round
	<-ch
	c.SetReadDeadline(time.Now().Add(10 * time.Second))
	err, updated := s.upd(reader, id)
	if err != nil {
		fmt.Println(err.Error())
	}
	done <- updated
	for update := range ch {
		s.send_upd(c, update)
		c.SetReadDeadline(time.Now().Add(2 * time.Second))
		err, updated := s.upd(reader, id)
		if err != nil {
			fmt.Println(err.Error())
		}
		done <- updated
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

func (s *server) bye(c net.Conn) {
	msg := []byte("BYE")
	c.Write(msg)
}
func (s *server) send_map(c net.Conn) {
	n := len(s.humans) + len(s.monster[0]) + len(s.monster[1])
	msg := make([]byte, 4+5*n)
	copy(msg, []byte("MAP"))
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

func (s *server) send_upd(c net.Conn, update []cell) {
	update = s.reload(update)
	n := len(update)
	msg := make([]byte, 4+5*n)
	copy(msg, []byte("UPD"))
	msg[3] = byte(uint8(n))
	for i, cl := range update {
		switch cl.kind {
		case human:
			msg[4+5*i] = byte(uint8(cl.X))
			msg[4+5*i+1] = byte(uint8(cl.Y))
			msg[4+5*i+2] = byte(uint8(cl.Count))
		case wolf:
			msg[4+5*i] = byte(uint8(cl.X))
			msg[4+5*i+1] = byte(uint8(cl.Y))
			msg[4+5*i+4] = byte(uint8(cl.Count))
		case vamp:
			msg[4+5*i] = byte(uint8(cl.X))
			msg[4+5*i+1] = byte(uint8(cl.Y))
			msg[4+5*i+3] = byte(uint8(cl.Count))
		}
	}
	c.Write(msg)
}

func (s *server) upd(reader *bufio.Reader, id int) (err error, update []cell) {
	buf := make([]byte, 5)
	_, err = io.ReadFull(reader, buf[:3])
	if err != nil {
		return err, update
	}
	if bytes.Compare(buf[:3], []byte("MOV")) != 0 {
		fmt.Println(string(buf[:3]), buf)
		return errors.New("Invalid MOV trame value"), update
	}

	_, err = io.ReadFull(reader, buf[:1])
	if err != nil {
		return err, update
	}
	t := int(buf[0])
	moves := make([]move, t)
	for i := 0; i < t; i++ {
		_, e := io.ReadFull(reader, buf[:5])
		if e != nil {
			err = e
			continue
		}
		moves[i] = move{
			oldx:      int(uint(buf[0])),
			oldy:      int(uint(buf[1])),
			count:     int(uint(buf[2])),
			effective: int(uint(buf[2])),
			newx:      int(uint(buf[3])),
			newy:      int(uint(buf[4])),
		}
	}
	if err != nil {
		return err, update
	}
	return s.apply(moves, id)
}
