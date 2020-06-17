package main

import (
	"flag"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ceruleandatahub/wirepas-sink-bridge/promistel"
	"github.com/ceruleandatahub/wirepas-sink-bridge/wirepas"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultPort    = "/dev/ttyUSB0"
	defaultBitrate = 115200
	defaultTimeout = 5
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
	config.timeout = defaultTimeout

	if v := os.Getenv("WIREPAS_SINK_PORT"); v != "" {
		config.port = v
	}
	if v := os.Getenv("SOCKET_PATH"); v != "" {
		config.socket = v
	}
	if v, err := strconv.Atoi(os.Getenv("WIREPAS_SINK_BITRATE")); err == nil {
		config.bitrate = v
	}
	if v, err := strconv.Atoi(os.Getenv("SOCKET_TIMEOUT")); err == nil {
		config.timeout = v
	}

	flag.StringVar(&config.port, "port", config.port, "Serial port where the sink is connected")
	flag.IntVar(&config.bitrate, "bitrate", config.bitrate, "Serial bitrate used by the sink")
	flag.StringVar(&config.socket, "socket", config.socket, "Path to unix socket where data is written (write to stdout if empty)")
	flag.IntVar(&config.timeout, "timeout", config.timeout, "Timeout in seconds to wait for the socket to become available")

	// Comment this to get JSON logging, this is for pretty human-readable logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	flag.Parse()

	var socket net.Conn
	var err error

	if config.socket != "" {
		log.Info().Str("PATH", config.socket).Msg("Establishing socket connection")

		timer := time.NewTicker(time.Duration(config.timeout) * time.Second)
		defer timer.Stop()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		connected := make(chan bool, 1)

	socketwait:
		for {
			select {
			case <-connected:
				timer.Stop()
				ticker.Stop()
				break socketwait
			case <-timer.C:
				log.Error().Msg("Socket connection timeout")
			case <-ticker.C:
				socket, err = net.Dial("unix", config.socket)
				if err != nil {
					continue
				}

				log.Info().Msg("Socket connected")
				defer socket.Close()
				connected <- true
				continue
			}
		}
		// socket, err = net.Dial("unix", config.socket)
		// if err != nil {
		// 	log.Fatal().Err(err).Str("SOCKET", config.socket).Msg("Unable to open unix socket")
		// }
		// defer socket.Close()
	}

	conn, err := wirepas.ConnectSink(config.port, config.bitrate)
	if err != nil {
		log.Fatal().Err(err)
		os.Exit(1)
	}
	defer conn.Close()
	c := conn.Listen()

	go func() {
		for msg := range c {
			if msg.SrcEP != wirepas.PwsEpSrcPromistel {
				// We only support Promistel RuuviTags for now
				log.Warn().Int("SRC", int(msg.SrcEP)).Msg("Unsupported source EP")
				continue
			}
			info, err := promistel.DecodeMessage(msg)
			if err != nil {
				log.Warn().Err(err)
				continue
			}
			json, err := info.JSON()
			if err != nil {
				log.Printf("Unable to convert message to JSON: %v\n", err)
				continue
			}
			log.Info().Int("SRC", int(msg.SrcEP)).Int("DST", int(msg.DstEP)).Str("JSON", json).Msg("Message received")
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
