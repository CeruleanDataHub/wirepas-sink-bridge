package main

import (
	"flag"
	"net"
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

   currentLoop "../proto"
//	"github.com/ceruleandatahub/wirepas-sink-bridge/promistel"
	"github.com/ceruleandatahub/wirepas-sink-bridge/wirepas"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

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
	log.Printf("looping here")
	grpcConnection, err := grpc.Dial("localhost:10024", grpc.WithInsecure())
	grpcClient := currentLoop.NewCurrentLoopClient(grpcConnection)
	connection, err := grpcClient.SendTelemetry(context.Background()) //creates connection, does not send
	if err != nil {
		log.Warn().Err(err)
	}

  //receive loop
	go func() {
		for {
			in, err := connection.Recv()
			if err != nil {
				log.Warn().Err(err)
			}
			log.Printf("return value: %v", in.Hash)
		}
	}()

	// write loop
	go func() {
		log.Printf("looping here")
		for msg := range c {
			if msg.SrcEP != wirepas.PwsEpSrcPromistel {
				// We only support Promistel RuuviTags for now
				log.Printf("wirepas detection failed")
				log.Warn().Int("SRC", int(msg.SrcEP)).Msg("Unsupported source EP")
				continue
			}
//			info, err := promistel.DecodeMessage(msg)
//			if err != nil {
//				log.Printf("decode failed")
//				log.Warn().Err(err)
//				continue
//			}
//			json, err := info.JSON()
//			if err != nil {
//				log.Printf("Unable to convert message to JSON: %v\n", err)
//				continue
//			}
//			log.Printf(json)
//			log.Info().Int("SRC", int(msg.SrcEP)).Int("DST", int(msg.DstEP)).Str("JSON", json).Msg("Message received")
			telemetry := &currentLoop.CurrentLoopRequest{
				Hash:      "10",
				Timestamp: "10",
				Value:     int32(10),
				Voltage:   float32(10.3),
				Current:   float32(10.3)}
			connection.Send(telemetry)

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
