package wirepas

/*
#cgo CFLAGS: -I${SRCDIR}/../include
#cgo LDFLAGS: -L${SRCDIR}/../libs -lwirepasmeshapi

#define LOG_MODULE_NAME "WSR"
#define MAX_LOG_LEVEL DEBUG_LOG_LEVEL

#include "wpc.h"
#include "logger.h"

bool onDataReceived_cgo(const uint8_t*, uint8_t, app_addr_t, app_addr_t, app_qos_e, uint8_t, uint8_t, uint32_t, uint8_t, unsigned long long);
bool onDiagReceived_cgo(const uint8_t*, uint8_t, app_addr_t, app_addr_t, app_qos_e, uint8_t, uint8_t, uint32_t, uint8_t, unsigned long long);
*/
import "C"

import (
	"errors"
	"fmt"
	"log"
	"unsafe"
)

type Conn struct {
	// ctx context.Context
}

//export onDataReceived
func onDataReceived(bytes *uint8, num_bytes uint8, src_addr C.app_addr_t, dst_addr C.app_addr_t, qos C.app_qos_e, src_ep uint8, dst_ep uint8, travel_time uint32, hop_count uint8, timestamp_ms_epoch uint64) bool {
	fmt.Println("onDataReceived::")
	fmt.Printf("    dst_ep: %d, src_ep: %d, len: %d, src_addr: 0x%x, dst_addr: 0x%x\n", dst_ep, src_ep, num_bytes, src_addr, dst_addr)
	var bs = C.GoBytes(unsafe.Pointer(bytes), C.int(num_bytes))
	fmt.Printf("    data: %v %v\n", bytes, bs)
	fmt.Print("    bytes: ")
	for _, b := range bs {
		fmt.Printf("0x%02x ", b)
	}
	fmt.Println()

	return false
}

//export onDiagReceived
func onDiagReceived(bytes *uint8, num_bytes uint8, src_addr C.app_addr_t, dst_addr C.app_addr_t, qos C.app_qos_e, src_ep uint8, dst_ep uint8, travel_time uint32, hop_count uint8, timestamp_ms_epoch uint64) bool {
	fmt.Println("onDiagReceived::")
	fmt.Printf("    dst_ep: %d, src_ep: %d, len: %d, src_addr: 0x%x, dst_addr: 0x%x\n", dst_ep, src_ep, num_bytes, src_addr, dst_addr)
	fmt.Print("    data: ")
	var bs = C.GoBytes(unsafe.Pointer(bytes), C.int(num_bytes))
	fmt.Println(bs)
	for _, b := range bs {
		fmt.Printf("0x%02x ", b)
	}
	fmt.Println()

	return false
}

func ConnectSink(port string, bitrate int) (*Conn, error) {
	conn := new(Conn)

	log.Println("Connecting to Wirepas sink")
	log.Println("Bitrate is ", bitrate)
	log.Println("Port is", port)

	if C.WPC_initialize(C.CString(port), C.ulong(bitrate)) != C.APP_RES_OK {
		return nil, errors.New("Failed to connect to Wirepas sink")
	}

	var mesh_version C.ushort

	// Do sanity check to test connectivity with sink
	if C.WPC_get_mesh_API_version(&mesh_version) != C.APP_RES_OK {
		return nil, errors.New("Cannot establish communication with sink over UART")
	}
	log.Printf("Wirepas sink connected, node is running mesh API version %d\n", mesh_version)

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
	log.Println("Wirepas stack started")

	// Register for diagnostics data on EP 255
	C.WPC_register_for_data(255, (C.onDataReceived_cb_f)(unsafe.Pointer(C.onDiagReceived_cgo)))

	// Register for data on all other EPs (EP 0 to 254)
	for i := 0; i < 255; i++ {
		C.WPC_register_for_data(C.uchar(i), (C.onDataReceived_cb_f)(unsafe.Pointer(C.onDataReceived_cgo)))
	}

	return conn, nil
}

func (c *Conn) Close() {
	log.Println("Closing Wirepas sink connection")
	C.WPC_close()
}
