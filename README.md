# HTTP/1.1 Protocol Implementation in Go

A learning project that recreates the HTTP/1.1 protocol from scratch using Go, focusing on understanding low-level network programming and protocol implementation.

## Project Overview

This project implements core HTTP/1.1 functionality by building up from the TCP layer, providing hands-on experience with:
- TCP socket programming
- HTTP request parsing
- Protocol state management
- Network I/O handling

## Architecture

### Core Components

#### 1. Request Parser (`internal/request/`)
- **HTTP Method Validation**: Regex-based validation for HTTP methods (GET, POST, etc.)
- **Request Target Parsing**: URI validation and parsing
- **HTTP Version Parsing**: Version string validation (HTTP/1.1, HTTP/2, etc.)
- **Request Line Parser**: Complete HTTP request line parsing with error handling
- **State Management**: Tracks parsing state (initialized, processing, done)

#### 2. TCP Listener (`cmd/tcplistener/`)
- **TCP Server**: Listens on port 42069 for incoming connections
- **Chunked Reading**: Reads data in 8-byte chunks to simulate real network conditions
- **Line-based Processing**: Handles CRLF-delimited HTTP messages
- **Concurrent Handling**: Uses goroutines for connection management

#### 3. UDP Sender (`cmd/udpsender/`)
- **Test Client**: Sends messages to the TCP listener for testing
- **Interactive Mode**: Command-line interface for manual testing

## Key Learning Concepts

### I/O Readers and File Descriptors
- Understanding how `io.Reader` interface works
- File descriptor management in Unix-like systems
- Difference between file I/O (position-based) and network I/O (stream-based)

### HTTP/1.1 Protocol Details
- Request line format: `METHOD request-target HTTP-version`
- CRLF line termination (`\r\n`)
- Header parsing and validation
- State machine implementation

### Network Programming
- TCP vs UDP differences
- Socket programming fundamentals
- Chunked data reading
- Connection lifecycle management

## Project Structure

```
httpfromtcp/
├── cmd/
│   ├── tcplistener/     # TCP server implementation
│   └── udpsender/       # UDP test client
├── internal/
│   └── request/         # HTTP request parsing logic
├── go.mod              # Go module definition
└── messages.txt        # Test data
```

## Usage

### Start the TCP Listener
```bash
go run cmd/tcplistener/main.go
```

### Send Test Messages
```bash
go run cmd/udpsender/main.go
```

## Implementation Status

- ✅ HTTP method validation
- ✅ Request target parsing
- ✅ HTTP version parsing
- ✅ Request line parsing
- ✅ TCP server with chunked reading
- ✅ Basic state management
- 🚧 Header parsing (TODO)
- 🚧 Response generation (TODO)
- 🚧 Error handling improvements (TODO)

## Learning Goals

This project demonstrates:
1. **Low-level network programming** - Working directly with TCP sockets
2. **Protocol implementation** - Building HTTP from the ground up
3. **Go concurrency** - Using goroutines for connection handling
4. **I/O patterns** - Understanding readers, writers, and file descriptors
5. **State management** - Tracking protocol parsing states

## Dependencies

- Go 1.25.1+
- Standard library only (no external HTTP libraries)

## Development Notes

- Uses 8-byte buffer size for chunked reading to simulate real network conditions
- Implements custom `io.Reader` patterns for learning purposes
- Focus on understanding the protocol rather than production-ready implementation
