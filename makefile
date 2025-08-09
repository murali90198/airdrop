SWIFT_SRC=AirDropBridge.swift
HEADER=airdrop.h
OBJ=AirDropBridge.o
STATIC_LIB=libairdrop.a
BIN=airdrop

all: $(BIN)

$(OBJ): $(SWIFT_SRC)
	swiftc -c -import-objc-header $(HEADER) $(SWIFT_SRC) -o $(OBJ)
	# swiftc -emit-library -o $(DYLIB) -import-objc-header $(HEADER) $(SWIFT_SRC) -framework Cocoa -framework Foundation

$(STATIC_LIB): $(OBJ)
	ar rcs $(STATIC_LIB) $(OBJ)
	# nm AirDropBridge.o | grep main

$(BIN): main.go $(STATIC_LIB)
	go build -o $(BIN) main.go
	rm -f $(OBJ) $(STATIC_LIB)

clean:
	rm -f $(OBJ) $(STATIC_LIB) $(BIN)

sign:
	codesign --force --sign - $(STATIC_LIB)
	codesign --force --sign - $(BIN)


.PHONY: all clean
