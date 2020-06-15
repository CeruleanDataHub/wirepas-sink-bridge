package wirepas

/*
#cgo CFLAGS: -I${SRCDIR}/../include
#cgo LDFLAGS: -L${SRCDIR}/../libs -lwirepasmeshapi

#define LOG_MODULE_NAME "WSR"

#include "wpc.h"
#include "logger.h"

void onDataReceived_cgo(const uint8_t*, uint8_t, app_addr_t, app_addr_t, app_qos_e, uint8_t, uint8_t, uint32_t, uint8_t, unsigned long long);
*/
import "C"

import (
	"errors"
	"sync"
	"unsafe"

	"github.com/rs/zerolog/log"
)

type Conn struct {
	// ctx context.Context
	listenerLock sync.Mutex
	listener     chan *Message
}

// We need a package level access to the connection as we need to call methods from the cgo callbacks
var (
	once sync.Once
	conn *Conn
)

//export onDataReceived
func onDataReceived(bytes *uint8, num_bytes uint8, src_addr uint32, dst_addr uint32, qos C.app_qos_e, src_ep uint8, dst_ep uint8, travel_time uint32, hop_count uint8, timestamp_ms_epoch uint64) {
	// fmt.Println("onDataReceived::")
	// fmt.Printf("    dst_ep: %d, src_ep: %d, len: %d, src_addr: %d, dst_addr: %d\n", dst_ep, src_ep, num_bytes, src_addr, dst_addr)
	// var bs = C.GoBytes(unsafe.Pointer(bytes), C.int(num_bytes))
	// fmt.Printf("    data: %v %v\n", bytes, bs)
	// fmt.Print("    bytes: ")
	// for _, b := range bs {
	// 	fmt.Printf("0x%02x ", b)
	// }
	// fmt.Println()

	msg := &Message{
		DstEP:      dst_ep,
		SrcEP:      src_ep,
		SrcAddress: src_addr,
		DstAddress: dst_addr,
		Payload:    C.GoBytes(unsafe.Pointer(bytes), C.int(num_bytes)),
	}

	conn.listenerLock.Lock()
	if conn.listener != nil {
		select {
		case conn.listener <- msg:
		default:
		}
	}
	conn.listenerLock.Unlock()
}

func ConnectSink(port string, bitrate int) (*Conn, error) {
	once.Do(func() {
		conn = new(Conn)
	})

	log.Info().Str("PORT", port).Int("BITRATE", bitrate).Msg("Connecting to Wirepas sink")

	if C.WPC_initialize(C.CString(port), C.ulong(bitrate)) != C.APP_RES_OK {
		return nil, errors.New("Failed to connect to Wirepas sink")
	}

	var mesh_version C.ushort

	// Do sanity check to test connectivity with sink
	if C.WPC_get_mesh_API_version(&mesh_version) != C.APP_RES_OK {
		return nil, errors.New("Cannot establish communication with sink over UART")
	}
	log.Info().Int("MESH VERSION", int(mesh_version)).Msg("Wirepas sink connected")

	// Get app config
	// var seq_p C.uchar
	// var interval_p C.ushort
	// var config C.uchar

	// if C.WPC_get_app_config_data(&seq_p, &interval_p, &config, 80) != C.APP_RES_OK {
	// 	log.Fatalln("Cannot get app config data")
	// }

	// Start the stack
	if C.WPC_start_stack() != C.APP_RES_OK {
		return nil, errors.New("Failed to start the Wirepas stack")
	}
	log.Info().Msg("Wirepas stack started")

	// Register for diagnostics data on EP 255
	// C.WPC_register_for_data(255, (C.onDataReceived_cb_f)(unsafe.Pointer(C.onDiagReceived_cgo)))

	return conn, nil
}

func (c *Conn) Close() {
	log.Info().Msg("Closing listener channel")
	c.listenerLock.Lock()
	if c.listener != nil {
		close(c.listener)
	}
	c.listenerLock.Unlock()

	log.Info().Msg("Closing Wirepas sink connection")
	C.WPC_close()
}

func (c *Conn) Listen() chan *Message {
	if c.listener != nil {
		log.Info().Msg("Wirepas sink listener already started, reusing existing listener")
		return c.listener
	}

	log.Info().Msg("Starting Wirepas sink listener")

	conn.listener = make(chan *Message, 10)

	// Register for data on all EPs (EP 0 to 255)
	for i := 0; i <= 255; i++ {
		C.WPC_register_for_data(C.uchar(i), (C.onDataReceived_cb_f)(unsafe.Pointer(C.onDataReceived_cgo)))
	}

	return c.listener
}
