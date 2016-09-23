package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "fmt"
import "crypto/rand"

func randString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func main() {
	db, err := sql.Open("mysql", "root:@/db")
	defer db.Close()
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("insert into idx (name) values(?);")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("begin")
	fmt.Println(err)
	for i := 0; i < 5000000; i++ {
		v := randString(25)
		stmt.Exec(v)
		if i%100000 == 0 {
			fmt.Println(i)
		}
	}
	_, err = db.Exec("commit")
	fmt.Println(err)
}
