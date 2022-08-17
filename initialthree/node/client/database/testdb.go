package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

func pgsqlOpen(host string, port int, dbname string, user string, password string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password)
	return sqlx.Open("postgres", connStr)
}

func mysqlOpen(host string, port int, dbname string, user string, password string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname)
	return sqlx.Open("mysql", connStr)
}

func sqlOpen(sqlType string, host string, port int, dbname string, user string, password string) (*sqlx.DB, error) {
	if sqlType == "mysql" {
		return mysqlOpen(host, port, dbname, user, password)
	} else {
		return pgsqlOpen(host, port, dbname, user, password)
	}
}

/*
CREATE TABLE "testdb" (
"__key__" varchar(255) NOT NULL,
"__version__" int8 NOT NULL,
"value" int8 NOT NULL DEFAULT 0,
PRIMARY KEY ("__key__")
);

*/

func main() {

	k := os.Args[1]
	rand.Seed(time.Now().UnixNano())

	if k == "1" {
		wg := sync.WaitGroup{}
		for i := 0; i < 5; i++ {
			//xdb, err := sqlOpen("pgsql", "10.50.31.13", 5432, "initialthree2", "sgzr2", "sniperHWfeiyu2019")
			xdb, err := sqlOpen("pgsql", "127.0.0.1", 5432, "yidongdeng", "dbuser", "123456")
			if err != nil {
				panic(err)
			}

			wg.Add(1)

			start := i*2000 + 1
			go func(xdb *sqlx.DB, s int) {
				i := s
				e := s + 2000
				for ; i < e; i++ {
					inset(xdb, "testdb", fmt.Sprintf("%d", i), map[string]interface{}{"value": 0})
				}

				wg.Done()
			}(xdb, start)
		}

		wg.Wait()
		fmt.Println("inset ok")
	} else if k == "2" {
		wg := sync.WaitGroup{}

		//xdb, err := sqlOpen("pgsql", "10.50.31.13", 5432, "initialthree2", "sgzr2", "sniperHWfeiyu2019")
		//xdb, err := sqlOpen("pgsql", "127.0.0.1", 5432, "yidongdeng", "dbuser", "123456")
		//if err != nil {
		//	panic(err)
		//}

		for i := 0; i < 100; i++ {
			xdb, err := sqlOpen("pgsql", "127.0.0.1", 5432, "yidongdeng", "dbuser", "123456")
			//xdb, err := sqlOpen("pgsql", "10.50.31.13", 5432, "initialthree2", "sgzr2", "sniperHWfeiyu2019")
			if err != nil {
				panic(err)
			}

			wg.Add(1)

			start := i*100 + 1
			go func(xdb *sqlx.DB, s int) {
				i := s
				e := s + 100
				for i < e {
					num := rand.Int()%5 + 5
					keys := make([]string, 0, num)
					for j := 0; j < num && i+j < 10000; j++ {
						keys = append(keys, fmt.Sprintf("%d", i+j))
					}
					i += num

					if len(keys) != 0 {
						checkKey(xdb, "testdb", keys)
					}
				}

				wg.Done()
			}(xdb, start)
		}

		wg.Wait()
		fmt.Println("get ok")
	}
}

func inset(xdb *sqlx.DB, tableName, key string, fields map[string]interface{}) error {
	sqlStr := `
INSERT INTO "%s" (__key__,__version__,%s)
VALUES (%s);`

	columns, values := []string{}, []string{"$1", "$2"}
	args := []interface{}{key, 1}
	var i = 3
	for k, v := range fields {
		columns = append(columns, k)
		values = append(values, fmt.Sprintf("$%d", i))
		i++
		args = append(args, v)
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(columns, ","), strings.Join(values, ","))
	//fmt.Println(sqlStatement)
	smt, err := xdb.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec(args...)
	return err

}

func checkKey(cdb *sqlx.DB, table string, keys []string) {

	sqlStr := `
SELECT __key__ FROM "%s" 
WHERE `

	s := fmt.Sprintf(sqlStr, table)
	ss := make([]string, 0, len(keys))
	for _, k := range keys {
		ss = append(ss, fmt.Sprintf("__key__ = '%s'", k))
	}

	s += strings.Join(ss, " or ") + ";"

	//fmt.Println(s)

	rows, err := cdb.Query(s)
	if err != nil {
		fmt.Println(s, err)
		panic(err)
	}
	defer rows.Close()

	findk := make([]string, 0, len(keys))
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			panic(err)
		}

		findk = append(findk, key)
		//fmt.Println(findk, key)

		find := false
		for _, k := range keys {
			if k == key {
				find = true
				break
			}
		}

		if !find {
			fmt.Printf(" %v 多读到key %s \n", keys, key)
		}
	}

	for _, k := range keys {
		find := false
		for _, kk := range findk {
			if kk == k {
				find = true
				break
			}
		}
		if !find {
			fmt.Printf(" %v %v 少读到key %s \n", keys, findk, k)
		}
	}

	//fmt.Printf(" %v %v ok \n", keys, findk)

}
