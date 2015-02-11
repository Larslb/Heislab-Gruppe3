package Driver
/*
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import "C"

func io_init()int {return int(C.io_init())}
func io_set_bit(channel int) {
	Set_conv_channel := C.int(channel)
	C.io_set_bit(Set_conv_channel)
}
func io_clear_bit(channel int) {
	Clear_conv_channel := C.int(channel)
	C.io_clear_bit(conv_channel)
}

func io_write_analog(int channel, int value){
	Write_conv_channel := C.int(channel)
	Write_conv_value := C.int(value)
	
	C.io_write_analaog(Write_conv_channel, Write_conv_value)
}


func io_read_bit(int channel) {
	Read_conv_channel := C.int(channel)
	return int(C.io_read_bit(Read_conv_channel))
}

func io_read_analog(int channel) {
	Read_conv_channel_2 := C.int(channel)
	return int(C.io_read_bit(Read_conv_channel_2))
}

