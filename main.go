package main

/*
#cgo LDFLAGS: -L. -lairdrop -framework Cocoa
#include "airdrop.h"
#include <stdlib.h>

*/
import "C"

import (
	"fmt"
	"io"
	"os"
	"unsafe"
)

// ./airdrop filename.png — передать файл
// cat somefile | ./airdrop

func main() {
	var files []string

	if len(os.Args) > 1 {
		// Есть аргумент - считаем, что это файл
		files = append(files, os.Args[1])
	} else {
		// Нет аргументов - читаем из stdin и сохраняем во временный файл
		tmpFile, err := os.CreateTemp("", "stdinfile_*")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to create temp file:", err)
			os.Exit(1)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		_, err = io.Copy(tmpFile, os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to read from stdin:", err)
			os.Exit(1)
		}

		// Обязательно закрываем файл, чтобы AirDrop мог его прочитать
		err = tmpFile.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to close temp file:", err)
			os.Exit(1)
		}

		files = append(files, tmpFile.Name())
	}

	// Конвертируем Go-строки в C-строки
	cFiles := make([]*C.char, len(files))
	for i, f := range files {
		cFiles[i] = C.CString(f)
		defer C.free(unsafe.Pointer(cFiles[i]))
	}

	result := C.ShareViaAirDrop((**C.char)(unsafe.Pointer(&cFiles[0])), C.int(len(cFiles)))

	if result != 0 {
		fmt.Println("AirDrop failed with code:", result)
		os.Exit(1)
	} else {
		fmt.Println("AirDrop succeeded")
	}
}
