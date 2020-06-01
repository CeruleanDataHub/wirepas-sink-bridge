package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/houston-inc/wirepas-sink-bridge/promistel"
	"github.com/houston-inc/wirepas-sink-bridge/wirepas"
)

var config struct {
	port    string
	bitrate int
	socket  string
}

func init() {
	flag.StringVar(&config.port, "port", "/dev/ttyUSB0", "Serial port where the sink is connected")
	flag.IntVar(&config.bitrate, "bitrate", 115200, "Serial bitrate used by the sink")
	flag.StringVar(&config.socket, "socket", "/tmp/wirepas.sock", "Path to unix socket where data is written")
}

func main() {
	flag.Parse()

	conn, err := wirepas.ConnectSink(config.port, config.bitrate)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := make(chan *wirepas.Message, 10)
	conn.Listen(c)

	socket, err := net.Dial("unix", config.socket)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to open unix socket %s (%v)", config.socket, err))
	}
	defer socket.Close()

	go func() {
		for msg := range c {
			info, err := promistel.DecodeWirepasMessage(msg)
			if err != nil {
				log.Print(err)
				continue
			}
			json, err := info.JSON()
			if err != nil {
				log.Print(err)
				continue
			}
			log.Printf("Message received on channel:\n%s\n", json)
			socket.Write([]byte(json + "\n"))
		}
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigs:
		log.Println("Received an interrupt, stopping...")
		close(done)
	}
	<-done
}
