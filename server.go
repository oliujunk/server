package main

import (
	"log"
	"oliujunk/server/apiserver"
	"oliujunk/server/commandserver"
	_ "oliujunk/server/config"
	_ "oliujunk/server/database"
)

func init() {
	// 日志信息添加文件名行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {

	go apiserver.Start()

	go commandserver.Start()

	//go communication.Start()

	select {}
}
