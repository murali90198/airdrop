SWIFT_SRC=AirDropBridge.swift
HEADER=airdrop.h
OBJ=AirDropBridge.o
STATIC_LIB=libairdrop.a
BIN=airdrop

ARCH ?= arm64

# make all ARCH=x86_64

all: $(BIN)

all-archs:
	$(MAKE) all ARCH=arm64
	mv $(BIN) $(BIN)
	$(MAKE) all ARCH=x86_64
	mv $(BIN) $(BIN)-x86_64


$(OBJ): $(SWIFT_SRC)
	# swiftc -c -import-objc-header $(HEADER) $(SWIFT_SRC) -o $(OBJ)
	# swiftc -emit-library -o $(DYLIB) -import-objc-header $(HEADER) $(SWIFT_SRC) -framework Cocoa -framework Foundation
	swiftc -c -import-objc-header $(HEADER) -parse-as-library $(SWIFT_SRC) -o $(OBJ)
	# Compile Swift source into an object file without main symbol to avoid duplicate symbols

$(STATIC_LIB): $(OBJ)
	ar rcs $(STATIC_LIB) $(OBJ)
	# Create static library from the object file

$(BIN): main.go $(STATIC_LIB)
	CGO_ENABLED=1 GOARCH=$(if $(filter arm64,$(ARCH)),arm64,amd64) go build -o $(BIN) main.go
	rm -f $(OBJ) $(STATIC_LIB)
	# Build Go binary linking with static Swift library and clean intermediate files


clean:
	rm -f $(OBJ) $(STATIC_LIB) $(BIN)

sign:
	codesign --force --sign - $(STATIC_LIB)
	codesign --force --sign - $(BIN)

.PHONY: all clean sign
