package main

import (
	"initialthree/node/client/mocker/mocker"
	"os"
	//"initialthree/node/node_game/log"
)

func main() {
	addr := os.Args[1]
	userID := os.Args[2]

	//logger := log.GetLogger()
	//kendynet.InitLogger(logger)

	client := mocker.NewClientUser(userID, addr)
	client.Run()
}
