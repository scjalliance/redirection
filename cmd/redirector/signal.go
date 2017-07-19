package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("redirector: received signal \"%v\"", s)
}
