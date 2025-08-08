# swiftc -c AirDropBridge.swift -o AirDropBridge.o
swiftc -c -parse-as-library AirDropBridge.swift -o AirDropBridge.o

ar rcs libairdrop.a AirDropBridge.o

# runtime lib libairdrop.a
