package main

import (
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/mortytheshorty/ssh-agent-wrapper/pkg/sshw/client"
)

func main() {

	if len(os.Args) < 1 {
		log.Println("missing argument")
		return
	}

	host := os.Args[1]

	u, err := user.Current()
	if err != nil {
		log.Fatalf("failed to retreive user info: %v\n", err)
		return
	}

	sockPath := strings.Join([]string{u.HomeDir, ".local/run/sshwd/sshwd.sock"}, "/")
	c, err := client.NewClient(sockPath)
	if err != nil {
		log.Fatalf("client failed to connect to unix socket: %v\n", err)
		return
	}

	keypath, err := c.Request(host)
	if err != nil {
		log.Fatalf("request failed: %v\n", err)
		return
	}

	if strings.Contains(keypath, "/") {
		cmd := exec.Command("ssh-add", keypath)

		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin

		err := cmd.Run()
		if err != nil {
			err = c.Failed()
			if err != nil {
				log.Fatalf("failed to acknowledge server: %v\n", err)
				return
			}
			log.Fatalf("failed to run ssh-add %v: %v", keypath, err)
			return
		}

		// log.Println("acknowledge")
		err = c.Acknowledge()
		if err != nil {
			log.Fatalf("failed to acknowledge daemon: %v", err)
			return
		}
	} else if keypath != "OK" {
		log.Fatalln("something failed")
		return
	}

	cmd := exec.Command("ssh", host)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		log.Fatalf("failed to run ssh %v: %v", host, err)
		return
	}
}
