package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/stathat/stathatgo"
)

var (
	listenAddr = flag.String("addr", "127.0.0.1:6007", "listen address")
	userKey    = flag.String("key", "", "stathat key")
)

func main() {
	flag.Parse()
	if *userKey == "" {
		fmt.Fprintln(os.Stderr, "You must provide a key.")
		flag.Usage()
		os.Exit(2)
	}
	log.Fatal(run())
}

func run() error {
	defer stathat.WaitUntilFinished(time.Minute)

	l, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		return err
	}
	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		go handle(c)
	}
	panic("unreachable")
}

func handle(c net.Conn) {
	defer c.Close()
	r := bufio.NewReader(c)
	b, err := r.ReadBytes('\n')
	if err != nil {
		log.Println(err)
		return
	}
	statKey := strings.TrimSpace(string(b))
	go stathat.PostEZCountOne(statKey, *userKey)
}
