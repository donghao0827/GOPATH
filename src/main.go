package main;

import (
	"bsserver"
	"middleware"
)

func main() {
	go bcserver.BCServer()
	middleware.Middleware("127.0.0.1:53")
}