package main

import (
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
	log.SetPrefix(fmt.Sprintf("[%d] [api-gateway] ", os.Getpid()))
}

func main() {
	log.Println("Hello, world!")
}
