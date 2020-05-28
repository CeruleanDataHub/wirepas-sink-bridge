package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/houston-inc/wirepas-sink-bridge/wirepas"
)

func main() {
	bitrate := flag.Int("bitrate", 115200, "bitrate")
	port := flag.String("port", "/dev/ttyUSB0", "port")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	flag.Parse()

	conn, err := wirepas.ConnectSink(ctx, *port, *bitrate)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	conn.OnDataReceived(129, func(s string) {
		log.Println("Data received: ", s)
	})

	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	// <-sig

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Received an interrupt, stopping...")
		close(done)
	}()
	<-done
}
