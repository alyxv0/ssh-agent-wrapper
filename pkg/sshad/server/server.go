package server

import (
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mortytheshorty/ssh-agent-wrapper/pkg/sshad/database"
)

type UnixSocketServer struct {
	socket    net.Listener
	db        *database.Database
	resetTime time.Duration
}

func NewUnixSocketServer(resetTime time.Duration) *UnixSocketServer {
	db, err := database.NewDb()
	if err != nil {
		log.Println("failed to create database:", err)
	}
	err = db.Init()
	if err != nil {
		log.Println("failed to initialize db:", err)
	}

	return &UnixSocketServer{db: db, resetTime: resetTime}
}

func requestHandler(conn net.Conn, s *UnixSocketServer) {
	defer conn.Close()

	buf := make([]byte, 4096)

	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("failed to read: %v\n", err)
	}

	wait := false
	msg := string(buf[:n])

	if strings.HasSuffix(msg, "\n") {
		// log.Println("found newline")
		msg = strings.Replace(msg, "\n", "", -1)
		// log.Println(msg)
	}

	host := s.db.Get(msg)
	if host != nil {
		if host.Loaded {
			n = copy(buf, "OK")
		} else {
			n = copy(buf, []byte(host.Keypath))
			wait = true
		}
	} else {
		n = copy(buf, "NOT FOUND")
	}

	_, err = conn.Write(buf[:n])
	if err != nil {
		log.Printf("failed to write: %v\n", err)
	}

	if wait {
		n, err = conn.Read(buf)
		if err != nil {
			log.Printf("failed to read answer: %v\n", err)
		}

		msg = string(buf[:n])
		if strings.HasSuffix(msg, "\n") {
			msg = strings.Replace(msg, "\n", "", -1)
		}

		if msg == "OK" {
			hosts := s.db.GetEqualKeys(host.Keypath)
			for _, h := range hosts {
				h.Loaded = true
				log.Println("host:", h)
			}

			// s.db.Print()

			go func() {
				// log.Printf("creating timer for removal of %v\n", hosts)
				time.Sleep(s.resetTime)

				cmd := exec.Command("ssh-add", "-d", host.Keypath)
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin

				err = cmd.Run()
				if err != nil {
					log.Printf("failed to run command: %v\n", err)
				} else {
					log.Printf("successfully removed %v\n", host.Keypath)
				}

				for _, host := range hosts {
					host.Loaded = false
				}

				// s.db.Print()
			}()
		} else if msg == "FAILED" {
			log.Println("client failed to add ssh-key to ssh-agent")
		}
	}

	// s.db.Print()
}

func (s *UnixSocketServer) Listen(path string) (err error) {

	s.socket, err = net.Listen("unix", path)
	if err != nil {
		return err
	}

	for {
		conn, err := s.socket.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v\n", err)
		}

		go requestHandler(conn, s)
	}

}
