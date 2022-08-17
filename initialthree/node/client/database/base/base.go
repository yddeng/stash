package base

import (
	"initialthree/node/common/db"
)

var client *db.Client

func GetClient() *db.Client {
	return client
}

func InitClient(sqlType string, host string, port int, dbname string, user string, password string) *db.Client {
	c, err := db.NewClient(sqlType, host, port, dbname, user, password)
	if err != nil {
		panic(err)
	}
	client = c
	return client
}
