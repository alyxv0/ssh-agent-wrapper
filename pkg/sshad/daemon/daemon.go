package daemon

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"syscall"

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

	d := &Daemon{
		Name: name,
		PidFileName:  name + ".pid",
		LogFileName:  name + ".log",
		SockFileName: name + ".sock",
		WorkDir:      strings.Join([]string{u.HomeDir, ".local", "run", name}, "/"),
	}

	err = os.MkdirAll(d.WorkDir, fs.FileMode.Perm(0750))
	if err != nil {
		log.Fatalf("failed to create workdir '%v': %v\n", d.WorkDir, err)
	}

	return d, nil
}

func removeContents(dir string) error {
    d, err := os.Open(dir)
    if err != nil {
        return err
    }
    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }
    for _, name := range names {
        err = removeContents(strings.Join([]string{dir, name}, "/"))
        if err != nil {
            return err
        }
    }
    return nil
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

	c := make(chan os.Signal, 1)
	signal.Notify(c,
				// syscall.SIGHUP,
				syscall.SIGINT,
				syscall.SIGTERM,
				// syscall.SIGQUIT
			)
	go func() {
		s := <-c
	
		log.Printf("caught %v\n", s.String())
		err = removeContents(d.WorkDir)
		if err != nil {
			log.Fatalf("failed to delete workdir on shutdown: %v", err)
			os.Exit(1)
		}

		os.Exit(1)
	}()
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
