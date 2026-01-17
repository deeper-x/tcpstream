package main

import (
	"log"
	"os"

	"github.com/deeper-x/tcpstream/client"
	"github.com/deeper-x/tcpstream/errorcs"
	"github.com/deeper-x/tcpstream/server"
	"github.com/deeper-x/tcpstream/settings"
)

func main() {
	if len(os.Args) != 2 {
		log.Println(settings.UsageStr)
		os.Exit(errorcs.ARGS)
	}

	cmdType := os.Args[1]

	if cmdType == "client" {
		err := client.Run()
		if err != nil {
			log.Println(err)
			os.Exit(errorcs.CLIENT)
		}
		return
	}

	if cmdType == "server" {
		err := server.Run()
		if err != nil {
			log.Println(err)
			os.Exit(errorcs.SERVER)
		}
		return
	}
}
