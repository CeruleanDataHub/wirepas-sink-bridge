package main

import (
	"flag"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ceruleandatahub/wirepas-sink-bridge/promistel"
	"github.com/ceruleandatahub/wirepas-sink-bridge/wirepas"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultPort    = "/dev/ttyUSB0"
	defaultBitrate = 115200
)

var config struct {
	port    string
	bitrate int
	socket  string
	timeout int
}

func init() {
	config.port = defaultPort
	config.bitrate = defaultBitrate

	if v := os.Getenv("WIREPAS_SINK_PORT"); v != "" {
		config.port = v
	}
	if v := os.Getenv("WIREPAS_SOCKET"); v != "" {
		config.socket = v
	}
	if v, err := strconv.Atoi(os.Getenv("WIREPAS_SINK_BITRATE")); err == nil {
		config.bitrate = v
	}

	flag.StringVar(&config.port, "port", config.port, "Serial port where the sink is connected")
	flag.IntVar(&config.bitrate, "bitrate", config.bitrate, "Serial bitrate used by the sink")
	flag.StringVar(&config.socket, "socket", config.socket, "Path to unix socket where data is written (write to stdout if empty)")
	flag.IntVar(&config.timeout, "timeout", config.timeout, "Timeout in seconds to wait for the socket to become available")

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	flag.Parse()

	var socket net.Conn
	var err error

	if config.socket != "" {
		socket, err = net.Dial("unix", config.socket)
		if err != nil {
			log.Fatal().Err(err).Str("SOCKET", config.socket).Msg("Unable to open unix socket")
		}
		defer socket.Close()
	}

	conn, err := wirepas.ConnectSink(config.port, config.bitrate)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer conn.Close()
	c := conn.Listen()

	go func() {
		for msg := range c {
			log.Printf("Message received on channel:\n%v\n", msg)
			if msg.SrcEP != wirepas.PwsEpSrcPromistel {
				// We only support Promistel RuuviTags for now
				log.Printf("Unsupported source EP %d", msg.SrcEP)
				continue
			}
			info, err := promistel.DecodeMessage(msg)
			if err != nil {
				log.Printf("Unable to decode message: %v\n", err)
				continue
			}
			json, err := info.JSON()
			if err != nil {
				log.Printf("Unable to convert message to JSON: %v\n", err)
				continue
			}
			log.Printf("Message received on channel:\n%s\n", json)
			if socket != nil {
				socket.Write([]byte(json + "\n"))
			}
		}
	}()

	sigs := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigs:
		log.Print("Received an interrupt, stopping...")
		close(done)
	}
	<-done
}
