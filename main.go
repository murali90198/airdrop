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

// tempFiles stores temporary files created during runtime that must be removed on exit
var tempFiles []string

func colorize(text, colorCode string) string {
	reset := "\033[0m"
	return colorCode + text + reset
}

// logf prints messages only if verbose mode is enabled
func logf(format string, a ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, format+"\n", a...)
	}
}

// cleanup removes all temporary files created by the program
func cleanup() {
	for _, f := range tempFiles {
		_ = os.Remove(f) // best effort cleanup
	}
}

// detectExtFromBytes tries to guess a file extension based on its first bytes
func detectExtFromBytes(sample []byte) string {
	ct := http.DetectContentType(sample)
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
	default:
		return ""
	}
}

func main() {
	// Ensure Go main goroutine stays bound to the current OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// CLI flags
	flag.BoolVar(&verbose, "v", false, "verbose logging")

	flag.Usage = func() {
		// ANSI escape codes для цветов
		cyan := "\033[36m"
		yellow := "\033[33m"
		bold := "\033[1m"

		fmt.Fprintf(os.Stderr, "%s\n\n", colorize("Copyright © 2025, Vitalii Tereshchuk | DOTOCA.NET All rights reserved | https://dotoca.net/airdrop", cyan))
		fmt.Fprintf(os.Stderr, "%s%s%s\n\n", bold, colorize("CLI Utility for Apple AirDrop — version 0.2.0", yellow), "\033[0m")
		fmt.Fprintf(os.Stderr, "%s\n", colorize("Usage:", bold))
		fmt.Fprintf(os.Stderr, "  %s [options] file1 file2 ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  cat file.pdf | %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s\n", colorize("Options:", bold))
		flag.PrintDefaults()
	}

	flag.Parse()

	// Check if no arguments and no stdin → show usage and exit
	info, _ := os.Stdin.Stat()
	noArgs := flag.NArg() == 0
	stdinIsPipe := (info.Mode() & os.ModeCharDevice) == 0

	if noArgs && !stdinIsPipe {
		flag.Usage()
		os.Exit(1)
	}

	// Ensure cleanup is called at exit
	defer cleanup()

	// Catch SIGINT/SIGTERM to cleanup temp files
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-c
		logf("caught signal %v — cleaning up", s)
		cleanup()
		os.Exit(130) // exit code 128 + signal number
	}()

	var files []string

	if flag.NArg() > 0 {
		// Files provided as command-line arguments
		files = append(files, flag.Args()...)
	} else {
		// Read from stdin into a temporary file
		const sniffLen = 512
		buf := make([]byte, sniffLen)
		n, err := io.ReadAtLeast(os.Stdin, buf, 1)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			fmt.Fprintln(os.Stderr, "failed to read stdin for sniffing:", err)
			os.Exit(2)
		}
		sample := buf[:n]

		ext := detectExtFromBytes(sample)
		tmp, err := os.CreateTemp("", "airdrop_stdin_*"+ext)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to create temp file:", err)
			os.Exit(3)
		}
		tmpName := tmp.Name()
		tempFiles = append(tempFiles, tmpName)

		if _, err = tmp.Write(sample); err != nil {
			_ = tmp.Close()
			fmt.Fprintln(os.Stderr, "failed to write to temp file:", err)
			os.Exit(4)
		}
		if _, err = io.Copy(tmp, os.Stdin); err != nil {
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

	// Check all files exist and are readable
	for _, f := range files {
		if _, err := os.Stat(f); err != nil {
			fmt.Fprintln(os.Stderr, "file not accessible:", f, ":", err)
			os.Exit(10)
		}
	}

	// Convert file paths to C strings
	cFiles := make([]*C.char, len(files))
	for i, f := range files {
		abs := f
		if !filepath.IsAbs(f) {
			if a, err := filepath.Abs(f); err == nil {
				abs = a
			}
		}
		cFiles[i] = C.CString(abs)
		defer C.free(unsafe.Pointer(cFiles[i]))
		logf("prepared file %d -> %s", i, abs)
	}

	if len(cFiles) == 0 {
		fmt.Fprintln(os.Stderr, "no files to share")
		os.Exit(11)
	}

	// Call AirDrop function from C
	ptr := (**C.char)(unsafe.Pointer(&cFiles[0]))
	res := C.ShareViaAirDrop(ptr, C.int(len(cFiles)))
	if res != 0 {
		fmt.Fprintln(os.Stderr, "AirDrop failed with code:", int(res))
		os.Exit(20 + int(res))
	}
	fmt.Println("AirDrop succeeded")
}
