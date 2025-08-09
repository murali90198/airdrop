SWIFT_SRC=AirDropBridge.swift
HEADER=airdrop.h
OBJ=AirDropBridge.o
STATIC_LIB=libairdrop.a
BIN=airdrop

all: $(BIN)

$(OBJ): $(SWIFT_SRC)
	# swiftc -c -import-objc-header $(HEADER) $(SWIFT_SRC) -o $(OBJ)
	# swiftc -emit-library -o $(DYLIB) -import-objc-header $(HEADER) $(SWIFT_SRC) -framework Cocoa -framework Foundation
	swiftc -c -import-objc-header $(HEADER) -parse-as-library $(SWIFT_SRC) -o $(OBJ)
	# Compile Swift source into an object file without main symbol to avoid duplicate symbols

$(STATIC_LIB): $(OBJ)
	ar rcs $(STATIC_LIB) $(OBJ)
	# Create static library from the object file

$(BIN): main.go $(STATIC_LIB)
	go build -o $(BIN) main.go
	rm -f $(OBJ) $(STATIC_LIB)
	# Build Go binary linking with static Swift library and clean intermediate files

clean:
	rm -f $(OBJ) $(STATIC_LIB) $(BIN)

sign:
	codesign --force --sign - $(STATIC_LIB)
	codesign --force --sign - $(BIN)

.PHONY: all clean sign
