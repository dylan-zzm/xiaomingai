package main

import (
	"log"

	"github.com/dylan-zzm/xiaomingai/cmd/chatlog"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	chatlog.Execute()
}
