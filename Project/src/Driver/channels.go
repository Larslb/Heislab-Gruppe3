 package Driver
/*
#cgo LDFLAGS: -lcomedi -lm
#include "channels.h"
*/
import (
	"C"
	"fmt"
)


func Printsomething() {

	fmt.println(OBSTRUCTION)
}
