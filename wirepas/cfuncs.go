package wirepas

/*

#include <stdio.h>
#include "wpc.h"

// The gateway functions

bool onDataReceived_cgo(const uint8_t* bytes, uint8_t num_bytes, app_addr_t src_addr, app_addr_t dst_addr, app_qos_e qos, uint8_t src_ep, uint8_t dst_ep, uint32_t travel_time, uint8_t hop_count, unsigned long long timestamp_ms_epoch)
{
	void onDataReceived(const uint8_t*, uint8_t, app_addr_t, app_addr_t, app_qos_e, uint8_t, uint8_t, uint32_t, uint8_t, unsigned long long);
	onDataReceived(bytes, num_bytes, src_addr, dst_addr, qos, src_ep, dst_ep, travel_time, hop_count, timestamp_ms_epoch);
	return true;
}
*/
import "C"
