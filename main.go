package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/houston-inc/wirepas-sink-bridge/wirepas"
)

func main() {
	bitrate := flag.Int("bitrate", 115200, "bitrate")
	port := flag.String("port", "/dev/ttyUSB0", "port")

	flag.Parse()

	conn, err := wirepas.ConnectSink(*port, *bitrate)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Printf("Got signal: %v\n", sig)
		log.Println("Shutting down")
		done <- true
	}()

	log.Println("Waiting for messages...")
	<-done
}
