# garsaud/go-mysql-ssh 🔐

Package `garsaud/go-mysql-ssh` implements a thin layer above `database/sql` to run MySQL queries through ssh and make it easier to retrieve rows.

## Install

```sh
go get github.com/garsaud/go-mysql-ssh
```

## Examples

```go
package main

import (
    db "github.com/garsaud/go-mysql-ssh"
    "database/sql"
)

type User struct {
    id uint
    email string
}

func main() {
    // Optional. Default value is "ssh.pem"
    db.PrivateKeyFilename = "private-key.pem"

    users := make([]*User, 0)

    db.FetchSSH(
        "sshlogin@example.com:22",
        "mysqllogin:mysqlpassword@mysql+tcp(127.0.0.1:3306)/databasename",
        "select id, email from users",
        func(row *sql.Rows) {
            user := new(User)
            row.Scan(&user.id, &user.email)
            users = append(users, user)
        },
    )

    // Direct MySQL connection without SSH is also possible:

    db.Fetch(
        "mysqllogin:mysqlpassword@tcp(example.com:3306)/databasename",
        "select id, email from users",
        func(row *sql.Rows) {
            user := new(User)
            row.Scan(&user.id, &user.email)
            users = append(users, user)
        },
    )
}
```
