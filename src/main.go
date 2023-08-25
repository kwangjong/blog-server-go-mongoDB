package main

import (
	"github.com/kwangjong/kwangjong.github.io/server"
)

func main() {
	server.LoadSecret()
	server.Run()
}