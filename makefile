SWIFT_SRC=AirDropBridge.swift
HEADER=airdrop.h
DYLIB=libairdrop.dylib

# Default target
all: $(DYLIB)

# Build Swift dynamic library
$(DYLIB): $(SWIFT_SRC)
    swiftc -emit-library -o $(DYLIB) -import-objc-header $(HEADER) $(SWIFT_SRC) -framework Cocoa -framework Foundation

# Run the Go application
run: $(DYLIB)
	DYLD_LIBRARY_PATH=. go run main.go

# Clean up built files
clean:
	rm -f $(DYLIB) $(HEADER)


build: $(DYLIB)
	# swiftc -emit-library -o $(DYLIB) -import-objc-header $(HEADER) $(SWIFT_SRC) -framework Cocoa -framework Foundation
	go build -o airdrop main.go

.PHONY: all run clean
