import Foundation
import Cocoa

@_cdecl("ShareViaAirDrop")
public func ShareViaAirDrop(_ cPaths: UnsafePointer<UnsafePointer<CChar>?>, _ count: Int32) -> Int32 {
    let args = UnsafeBufferPointer(start: cPaths, count: Int(count))
    let paths: [String] = args.compactMap { $0.flatMap { String(cString: $0) } }
    let urls = paths.map { URL(fileURLWithPath: $0) }

    print("🧾 Files to send via AirDrop:")
    for url in urls {
        print("  - \(url.path)")
    }

    guard let service = NSSharingService(named: .sendViaAirDrop) else {
        print("❌ Failed to get AirDrop service")
        return 1
    }

    if !service.canPerform(withItems: urls) {
        print("❌ AirDrop cannot perform with given items.")
        return 1
    }

    class AirDropDelegate: NSObject, NSSharingServiceDelegate {
        func sharingService(_ sharingService: NSSharingService, willShareItems items: [Any]) {
            print("📤 Preparing to share \(items.count) item(s)")
        }

        func sharingService(_ sharingService: NSSharingService, didShareItems items: [Any]) {
            print("✅ Shared successfully")
            exit(0)
        }

        func sharingService(_ sharingService: NSSharingService, didFailToShareItems items: [Any], error: Error) {
            print("❌ Failed to share: \(error.localizedDescription)")
            exit(1)
        }

        func sharingService(_ sharingService: NSSharingService, sourceFrameOnScreenForShareItem item: Any) -> NSRect {
            return NSRect(x: 0, y: 0, width: 400, height: 100)
        }

        func sharingService(_ sharingService: NSSharingService, sourceWindowForShareItems items: [Any], sharingContentScope: UnsafeMutablePointer<NSSharingService.SharingContentScope>) -> NSWindow? {
            let window = NSWindow(contentRect: NSRect(x: 0, y: 0, width: 1, height: 1),
                                  styleMask: [],
                                  backing: .buffered,
                                  defer: false)
            window.level = .floating
            window.center()
            window.orderFrontRegardless()
            return window
        }
    }

    let delegate = AirDropDelegate()
    service.delegate = delegate

    // Ensure app is running with full UI access
    let app = NSApplication.shared
    app.setActivationPolicy(.regular)
    app.activate(ignoringOtherApps: true)

    service.perform(withItems: urls)

    // Run event loop to allow AirDrop UI interaction
    let runLoop = RunLoop.current
    while true {
        runLoop.run(until: Date(timeIntervalSinceNow: 0.1))
    }
}
