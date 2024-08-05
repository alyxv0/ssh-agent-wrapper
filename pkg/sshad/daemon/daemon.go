package daemon

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/sevlyar/go-daemon"
)

type Runner func() error

type Daemon struct {
	Name         string
	PidFileName  string
	LogFileName  string
	SockFileName string
	WorkDir      string
}

func NewDaemon(name string) (*Daemon, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	return &Daemon{
		PidFileName:  name + ".pid",
		LogFileName:  name + ".log",
		SockFileName: name + ".sock",
		WorkDir:      strings.Join([]string{u.HomeDir, ".local", "run", name}, "/"),
	}, nil
}

func (d *Daemon) Run(runf Runner) error {

	cntxt := &daemon.Context{
		PidFileName: strings.Join([]string{d.WorkDir, d.PidFileName}, "/"),
		PidFilePerm: 0644,
		LogFileName: strings.Join([]string{d.WorkDir, d.LogFileName}, "/"),
		// LogFileName: d.LogFileName,
		LogFilePerm: 0640,
		WorkDir:     d.WorkDir,
		Umask:       027,
		Args:        []string{d.Name},
	}

	d2, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d2 != nil {
		return nil
	}
	defer cntxt.Release()


	content, err := os.ReadFile(strings.Join([]string{d.WorkDir, d.PidFileName}, "/"))
	if err != nil {
		return fmt.Errorf("failed to read pid file content: %v", err)
	}
	log.Print("- - - - - - - - - - - - - - -")
	log.Printf("sshwd started; pid[%v]\n", string(content))
	// log.Printf("sshwd pid = %v\n", string(content))
	err = runf()
	if err != nil {
		return fmt.Errorf("runner function failed: %v", err)
	}

	return nil
}
