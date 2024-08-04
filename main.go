package main

import (
	"log"
	"os"
	"sshwd/daemon"
	"sshwd/server"
	"strings"
	"time"
)

func main() {

	d, err := daemon.NewDaemon("sshwd")
	if err != nil {
		log.Fatal(err)
	}

	err = d.Run(func() error {

		s := server.NewUnixSocketServer(20 * time.Second)
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
