package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/mortytheshorty/ssh-agent-wrapper/pkg/sshad/daemon"
	"github.com/mortytheshorty/ssh-agent-wrapper/pkg/sshad/server"
)

func main() {

	d, err := daemon.NewDaemon("sshad")
	if err != nil {
		log.Fatal(err)
	}

	// already running inside daemon context
	err = d.Run(func() error {

		s := server.NewUnixSocketServer(600 * time.Second)
		if err != nil {
			return err
		}

		log.Println("WorkDir:", d.WorkDir)
		sockpath := strings.Join([]string{d.WorkDir, d.SockFileName}, "/")
		log.Println("SocketPath:", sockpath)

		err = os.Remove(sockpath)
		if err != nil {
			log.Println("failed to remove unix socket:", err)
		}

		err = s.Listen(sockpath)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalf("failed to run daemon context: %v\n", err)
		return
	}

}
