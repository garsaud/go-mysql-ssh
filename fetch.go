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

func Fetch(uri string, query string, callback func(row *sql.Rows)) error {
	mysqlDb, err := sql.Open("mysql", uri)
	if err != nil { return err }
	defer mysqlDb.Close()

	rows, err := mysqlDb.Query(query)
	if err != nil { return err }

	for rows.Next() {
		callback(rows)
	}

	return nil
}

func FetchSSH(sshUri string, uri string, query string, callback func(row *sql.Rows)) error {
	pemBytes, err := ioutil.ReadFile("ssh.pem")
	if err != nil { return err }
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil { return err }

	sshUriParts := strings.Split(sshUri, "@")
	sshcon, err := ssh.Dial("tcp", sshUriParts[1], &ssh.ClientConfig {
		User: sshUriParts[0],
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil { return err }
	defer sshcon.Close()

	mysql.RegisterDialContext("mysql+tcp", func(_ context.Context, addr string) (net.Conn, error) {
		return sshcon.Dial("tcp", addr)
	})

	return Fetch(uri, query, callback)
}
