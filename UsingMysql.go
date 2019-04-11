package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"time"
)

func main() {
	sqlconf := mysql.NewConfig()

	sqlconf.User = "username"
	sqlconf.Passwd = "password"
	sqlconf.Addr = "ip:port"
	sqlconf.Net = "tcp"
	sqlconf.DBName = "dbname"
	sqlconf.Timeout = 2 * time.Second

	db, err := sql.Open("mysql", sqlconf.FormatDSN()+"&charset=utf8")
	fmt.Println(sqlconf.FormatDSN())

	if err != nil {
		fmt.Println(err)
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	if _, err := db.Exec("use mdb;"); err != nil {
		fmt.Println(err)
		return
	}

	rows, err := db.Query("select * from table1;")
	if nil != err {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	var data []interface{}
	for rows.Next() {

		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		err = rows.Scan(scanArgs...)

		record := make(map[string]interface{})
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		data = append(data, record)
	}
	buf, _ := json.Marshal(data)
	fmt.Printf("%s\n", buf)

	var answer string
	fmt.Scanln(&answer)
}
