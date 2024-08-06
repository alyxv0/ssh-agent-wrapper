package database

import (
	"log"
	"os"
	"os/user"
	"strings"
)

type HostKeyEntry struct {
	Host    string
	Keypath string
	Loaded  bool
}

type Database struct {
	db       []*HostKeyEntry
	ssh_path string
}

func NewDb() (*Database, error) {
	db := &Database{}

	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	db.ssh_path = strings.Join([]string{u.HomeDir, ".ssh"}, "/")

	return db, nil
}

func (db *Database) Init() error {
	// we read the ~/.ssh/config file for building the database
	bcontent, err := os.ReadFile(strings.Join([]string{db.ssh_path, "config"}, "/"))
	if err != nil {
		return err
	}

	config := strings.TrimSpace(string(bcontent))
	lines := strings.Split(config, "\n")

	host := ""
	keypath := ""

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// we skip every empty string slice if host and keypath are "" (empty)
		// otherwise we append a new host key entry to database
		if len([]byte(line)) == 0 {
			if host == "" && keypath == "" {
				continue
			}

			db.db = append(db.db, &HostKeyEntry{Host: host, Keypath: strings.Join([]string{db.ssh_path, keypath}, "/"), Loaded: false})
			host = ""
			keypath = ""
			log.Println("loaded ", db.db[len(db.db)-1])
			continue
		}

		// split by space
		splitted := strings.Split(line, " ")
		if splitted[0] == "Host" {
			host = splitted[1]
			continue
		}

		if splitted[0] == "IdentityFile" {
			// get dirname
			dsplit := strings.Split(splitted[1], "/")
			dsplit = dsplit[len(dsplit)-1:]
			keypath = dsplit[0]

			// if it is the last line from file, append last entry from here
			if i == len(lines)-1 {
				db.db = append(db.db, &HostKeyEntry{Host: host, Keypath: strings.Join([]string{db.ssh_path, keypath}, "/"), Loaded: false})
				log.Println("loaded ", db.db[len(db.db)-1])
			}
			continue
		}
	}

	return nil
}

func (db *Database) Get(host string) *HostKeyEntry {

	for _, h := range db.db {
		if h.Host == host {
			return h
		}
	}
	return nil
}

func (db *Database) GetEqualKeys(keypath string) []*HostKeyEntry {

	hke_list := []*HostKeyEntry{}

	for _, e := range db.db {
		if e.Keypath == keypath {
			hke_list = append(hke_list, e)
		}
	}

	return hke_list
}

func (db *Database) Print() {
	for i, n := range db.db {
		log.Printf("%v: Host: %v; Keypath: %v; Loaded: %v\n", i, n.Host, n.Keypath, n.Loaded)
	}
}
