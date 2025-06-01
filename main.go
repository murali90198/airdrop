package main

/*
#cgo LDFLAGS: -L. -lairdrop -framework Cocoa
#include "airdrop.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	files := []string{"./test.png"}
	cFiles := make([]*C.char, len(files))
	for i, f := range files {
		cFiles[i] = C.CString(f)
		defer C.free(unsafe.Pointer(cFiles[i]))
	}

	result := C.ShareViaAirDrop((**C.char)(unsafe.Pointer(&cFiles[0])), C.int(len(cFiles)))

	if result != 0 {
		fmt.Println("AirDrop failed with code:", result)
	} else {
		fmt.Println("AirDrop succeeded")
	}
}
