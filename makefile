SWIFT_SRC=AirDropBridge.swift
HEADER=airdrop.h
DYLIB=libairdrop.dylib
OLIB=AirDropBridge.o
BIN=airdrop

all: $(DYLIB) $(BIN)

$(DYLIB): $(SWIFT_SRC)
	# swiftc -emit-library -o $(DYLIB) -import-objc-header $(HEADER) $(SWIFT_SRC) -framework Cocoa -framework Foundation
	swiftc -c -parse-as-library $(SWIFT_SRC) -o $(OLIB)
	# nm AirDropBridge.o | grep main
	ar rcs libairdrop.a $(OLIB)

$(BIN): main.go $(DYLIB)
	go build -o $(BIN) main.go


run: $(DYLIB)
	DYLD_LIBRARY_PATH=. go run main.go

clean:
	rm -f $(DYLIB) $(BIN)

sign:
	codesign --force --sign - $(DYLIB)
	codesign --force --sign - $(BIN)

.PHONY: all run clean sign
