# HTTP/1.1 Server Implementation from TCP in Go

A complete HTTP/1.1 server implementation built from scratch using only Go's TCP layer. This project demonstrates low-level network programming by recreating the HTTP protocol without using any high-level HTTP libraries.

## What This Project Does

This is a fully functional HTTP/1.1 server that:
- **Listens for TCP connections** on port 42069
- **Parses raw HTTP requests** byte-by-byte from the TCP stream
- **Routes requests** to appropriate handlers
- **Generates proper HTTP/1.1 responses** with status lines, headers, and body
- **Handles errors gracefully** with appropriate HTTP status codes

All of this is built directly on top of TCP sockets without using Go's `net/http` package.

## Quick Start with Docker

The easiest way to run this project is using Docker:

```bash
# Build and start the server
docker-compose up --build

# Or run in detached mode (background)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the server
docker-compose down
```

Once running, test the server:

```bash
# Default endpoint - returns 200 OK
curl http://localhost:42069/

# Error endpoint - returns 400 Bad Request
curl http://localhost:42069/yourproblem

# Another error endpoint - returns 500 Internal Server Error
curl http://localhost:42069/myproblem
```

## Running Without Docker

If you prefer to run directly with Go:

```bash
# Run the HTTP server
go run cmd/httpserver/main.go

# Or run the basic TCP listener (development version)
go run cmd/tcplistener/main.go
```

## Architecture Deep Dive

### How It Works

When a client connects to the server, here's what happens:

1. **TCP Connection**: The server accepts a raw TCP connection on port 42069
2. **Chunked Reading**: Data is read from the socket in 8-byte chunks to simulate real network conditions
3. **Request Parsing**: The raw bytes are parsed into:
   - Request line (method, target, HTTP version)
   - Headers (key-value pairs)
   - Body (request payload)
4. **Handler Execution**: The parsed request is routed to your custom handler function
5. **Response Generation**: A proper HTTP/1.1 response is constructed with:
   - Status line (e.g., `HTTP/1.1 200 OK`)
   - Headers (`Content-Length`, `Connection`, `Content-Type`)
   - Body (the actual response data)
6. **Write Back**: The response is written back to the TCP connection
7. **Connection Close**: The connection is closed (HTTP/1.1 with `Connection: close`)

### Core Components

#### 📦 `internal/request/` - Request Parsing
The heart of HTTP request parsing:
- **Method Validation**: Regex-based validation for HTTP methods (GET, POST, PUT, DELETE, etc.)
- **Request Target Parsing**: URI validation and extraction
- **HTTP Version Parsing**: Validates HTTP version strings (HTTP/1.1, HTTP/2, etc.)
- **Header Parsing**: Parses key-value header pairs from the request
- **Body Reading**: Reads request body based on Content-Length header
- **State Machine**: Manages parsing state (request line → headers → body)

#### 📦 `internal/response/` - Response Generation
Handles creating proper HTTP responses:
- **Status Line Writing**: Formats status codes (200, 400, 500) into HTTP status lines
- **Header Writing**: Writes headers with proper CRLF formatting
- **Body Writing**: Writes response body with correct Content-Length
- **Default Headers**: Automatically adds `Content-Length`, `Connection: close`, and `Content-Type`

#### 📦 `internal/server/` - HTTP Server
The main server implementation:
- **TCP Listener**: Creates and manages the TCP listener socket
- **Connection Handling**: Accepts incoming connections and spawns goroutines
- **Handler Interface**: Defines the handler function signature
- **Error Management**: Converts handler errors into proper HTTP error responses
- **Graceful Shutdown**: Supports clean server shutdown

#### 📦 `internal/headers/` - Header Management
Manages HTTP header parsing and storage:
- **Header Parsing**: Extracts key-value pairs from header lines
- **Header Validation**: Ensures headers follow the `Key: Value` format
- **Storage**: Provides a map-based structure for header access

#### 📦 `internal/httpErrors/` - Error Handling
Centralized error management:
- **Exception Types**: Defines constants for different error types
- **Error Messages**: Maps exception types to error message generators
- **Formatted Errors**: Provides consistent error formatting with context

#### 📦 `internal/buffer/` - Dynamic Buffer
Custom buffer implementation for handling streaming data:
- **Dynamic Growth**: Automatically doubles buffer size when needed
- **Buffer Shifting**: Efficiently shifts data after partial reads
- **Capacity Management**: Maintains minimum buffer sizes for performance

### Project Structure

```
httpfromtcp/
├── cmd/
│   ├── httpserver/        # Main HTTP server with routing
│   ├── tcplistener/       # Basic TCP listener (development/debugging)
│   └── udpsender/         # UDP test client for sending test data
├── internal/
│   ├── request/           # HTTP request parsing
│   ├── response/          # HTTP response generation
│   ├── server/            # Server implementation
│   ├── headers/           # Header parsing and management
│   ├── httpErrors/        # Error handling
│   └── buffer/            # Dynamic buffer implementation
├── docker-compose.yaml    # Docker Compose configuration
├── Dockerfile             # Multi-stage Docker build
├── go.mod                 # Go module definition
└── messages.txt           # Test data
```

## Key Learning Concepts

### 1. Protocol Implementation
Understanding HTTP/1.1 at the byte level:
- Request line format: `METHOD request-target HTTP-version\r\n`
- Header format: `Key: Value\r\n`
- Header termination: Empty line `\r\n\r\n`
- Body length determined by `Content-Length` header

### 2. TCP Socket Programming
Working directly with TCP connections:
- Accepting connections from a listener
- Reading from a stream-based socket (no seek, no rewind)
- Writing responses back through the socket
- Managing connection lifecycle

### 3. State Machine Design
Parsing HTTP requests requires state management:
- **State 1**: Parse request line
- **State 2**: Parse headers (line by line)
- **State 3**: Parse body (based on Content-Length)
- Transitions happen as parsing progresses

### 4. I/O Patterns
Understanding Go's I/O interfaces:
- `io.Reader`: Reading from network connections
- `io.Writer`: Writing responses back to clients
- Buffering strategies for efficient I/O
- Chunked reading to simulate real network conditions

### 5. Concurrent Programming
Using goroutines for handling multiple connections:
- One goroutine per connection
- Non-blocking accept loop
- Clean resource cleanup with deferred connection closes

## Handler API

Create custom handlers by matching this signature:

```go
func handler(w io.Writer, req *request.Request) *server.HandlerError {
    // Check the request target (route)
    switch req.RequestLine.RequestTarget {
    case "/your-route":
        // Return an error response
        return &server.HandlerError{
            StatusCode: response.HttpStatusBadRequest,
            Message: []byte("Error message\n"),
        }
    }
    
    // Write success response
    w.Write([]byte("Success response\n"))
    return nil  // nil = 200 OK
}
```

The handler receives:
- `w io.Writer`: Write your response body here
- `req *request.Request`: Contains method, target, headers, and body

Return:
- `nil`: Generates a 200 OK response with your written content
- `*server.HandlerError`: Generates an error response with the specified status code

## Implementation Highlights

- ✅ **Full HTTP/1.1 request parsing** (method, target, version, headers, body)
- ✅ **HTTP response generation** (status line, headers, body)
- ✅ **Custom routing** via handler functions
- ✅ **Error handling** with proper HTTP status codes
- ✅ **Concurrent connection handling** with goroutines
- ✅ **Dynamic buffer management** for streaming data
- ✅ **Graceful shutdown** with signal handling
- ✅ **Docker support** for easy deployment

## Technologies Used

- **Go 1.25.1+**: The only dependency
- **Standard Library Only**: No external HTTP libraries
- **Docker**: For containerized deployment

## Why This Project?

This project exists to:
1. **Understand HTTP deeply** - By building it from scratch, you understand every detail
2. **Learn TCP networking** - Work directly with sockets and streams
3. **Master Go I/O** - Understand readers, writers, and buffering
4. **Practice systems programming** - Deal with bytes, state machines, and protocols

This is a learning project focused on understanding rather than production use. For production, use Go's excellent `net/http` package!

## Development

The project includes a TCP listener for debugging raw requests:

```bash
# Start the TCP listener (shows parsed request details)
go run cmd/tcplistener/main.go
```

This prints detailed information about each incoming request for debugging purposes.

## License

MIT

---