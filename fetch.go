package db

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"strings"
)

func fetch(uri string, query string, callback func(row *sql.Rows)) error {
	db, err := sql.Open("mysql", uri)
	if err != nil { return err }
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil { return err }

	for rows.Next() {
		callback(rows)
	}

	return nil
}

func fetchSSH(sshuri string, uri string, query string, callback func(row *sql.Rows)) error {
	pemBytes, err := ioutil.ReadFile("ssh.pem")
	if err != nil { return err }
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil { return err }

	sshuriParts := strings.Split(sshuri, "@")
	sshcon, err := ssh.Dial("tcp", sshuriParts[1], &ssh.ClientConfig {
		User: sshuriParts[0],
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil { return err }
	defer sshcon.Close()

	mysql.RegisterDialContext("mysql+tcp", func(_ context.Context, addr string) (net.Conn, error) {
		return (&ViaSSHDialer{sshcon}).Dial(addr)
	})

	return fetch(uri, query, callback)
}

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Dial(addr string) (net.Conn, error) {
	return self.client.Dial("tcp", addr)
}
