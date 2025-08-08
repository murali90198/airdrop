package main

/*
#cgo LDFLAGS: -L. -lairdrop -framework Cocoa
#include "airdrop.h"
#include <stdlib.h>
*/
import "C"

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

var (
	verbose bool
)

// tempFiles keeps track of temp files we must remove on exit
var tempFiles []string

func logf(format string, a ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, format+"\n", a...)
	}
}

func cleanup() {
	for _, f := range tempFiles {
		// best-effort remove
		_ = os.Remove(f)
	}
}

func mkTempFromStdin(suggestExt string) (string, error) {
	tmp, err := os.CreateTemp("", "airdrop_stdin_*"+suggestExt)
	if err != nil {
		return "", err
	}
	name := tmp.Name()
	tempFiles = append(tempFiles, name)

	// copy stdin to file
	_, err = io.Copy(tmp, os.Stdin)
	if cerr := tmp.Close(); cerr != nil && err == nil {
		err = cerr
	}
	if err != nil {
		return "", err
	}
	return name, nil
}

func detectExtFromBytes(sample []byte) string {
	// DetectContentType needs at least the initial bytes.
	ct := http.DetectContentType(sample)
	// map some common content-types to ext
	switch {
	case strings.HasPrefix(ct, "image/png"):
		return ".png"
	case strings.HasPrefix(ct, "image/jpeg"):
		return ".jpg"
	case strings.HasPrefix(ct, "image/gif"):
		return ".gif"
	case strings.HasPrefix(ct, "application/pdf"):
		return ".pdf"
	case strings.HasPrefix(ct, "text/"):
		return ".txt"
	// add more as needed
	default:
		return ""
	}
}

func main() {
	runtime.LockOSThread() // Закрепляем main за текущим ОС-потоком
	defer runtime.UnlockOSThread()

	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.Parse()

	// ensure cleanup on exit and on signals
	defer cleanup()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-c
		logf("caught signal %v — cleaning up", s)
		cleanup()
		os.Exit(130) // 128 + SIGINT
	}()

	var files []string

	if flag.NArg() > 0 {
		// Accept all provided args as files
		files = append(files, flag.Args()...)
	} else {
		// No args — read stdin into temp file.
		// Read a small sample to detect MIME / extension
		const sniffLen = 512
		buf := make([]byte, sniffLen)
		n, err := io.ReadAtLeast(os.Stdin, buf, 1)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			fmt.Fprintln(os.Stderr, "failed to read stdin for sniffing:", err)
			os.Exit(2)
		}
		sample := buf[:n]

		ext := detectExtFromBytes(sample)
		// create temp file with extension (we'll append remaining stdin)
		tmp, err := os.CreateTemp("", "airdrop_stdin_*"+ext)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to create temp file:", err)
			os.Exit(3)
		}
		tmpName := tmp.Name()
		tempFiles = append(tempFiles, tmpName)

		// write the sample and then the rest
		_, err = tmp.Write(sample)
		if err != nil {
			_ = tmp.Close()
			fmt.Fprintln(os.Stderr, "failed to write to temp file:", err)
			os.Exit(4)
		}
		_, err = io.Copy(tmp, os.Stdin)
		if err != nil {
			_ = tmp.Close()
			fmt.Fprintln(os.Stderr, "failed to write rest of stdin:", err)
			os.Exit(5)
		}
		if err := tmp.Close(); err != nil {
			fmt.Fprintln(os.Stderr, "failed to close temp file:", err)
			os.Exit(6)
		}
		files = append(files, tmpName)
		logf("stdin saved to %s (detected ext %s)", tmpName, ext)
	}

	// Validate files exist and are readable
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			fmt.Fprintln(os.Stderr, "file not accessible:", f, ":", err)
			os.Exit(10)
		}
	}

	// Build C array
	cFiles := make([]*C.char, len(files))
	// Keep track to free
	for i, f := range files {
		abs := f
		// convert to absolute path — safer for C code
		if !filepath.IsAbs(f) {
			if a, err := filepath.Abs(f); err == nil {
				abs = a
			}
		}
		cFiles[i] = C.CString(abs)
		// defer frees only after we call C function - OK here.
		defer C.free(unsafe.Pointer(cFiles[i]))
		logf("prepared file %d -> %s", i, abs)
	}

	// Call into C
	if len(cFiles) == 0 {
		fmt.Fprintln(os.Stderr, "no files to share")
		os.Exit(11)
	}

	// Pointer to first element
	ptr := (**C.char)(unsafe.Pointer(&cFiles[0]))
	res := C.ShareViaAirDrop(ptr, C.int(len(cFiles)))
	if res != 0 {
		fmt.Fprintln(os.Stderr, "AirDrop failed with code:", int(res))
		os.Exit(20 + int(res))
	}
	fmt.Println("AirDrop succeeded")
}
