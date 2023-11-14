package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

// Interface for the server program --
type FileServer struct{}

// Start the server --
// Listen on port 3000.
// For each connection, read from that connection.
func (fs *FileServer) start() {
    // Create a listener on port 3000
    ln, err := net.Listen("tcp", ":3000")
    if err != nil {
        log.Fatal(err)
    }
    defer ln.Close()
    // For each connection to the listener, read the data on that connection.
    for {
        // Create a new connection from the listener.
        conn, err := ln.Accept()
        if err != nil {
            log.Fatal("Error creating a connection from the listener", err)
        }
        defer conn.Close()
        // Read data from the connection
        go fs.readLoop(conn)
    }
}

// Read data from the connection --
// Get the size of the incoming data, then copy that data to a buffer.
func (fs *FileServer) readLoop(conn net.Conn) {
    // Create a buffer for the data read from the connection.
    buf := new(bytes.Buffer)
    // While there are no connection errors, read from the connection
    for {
        // Read the size of the incoming data.
        // Break out of the loop if there is an error with the connection.
        var size int64
        err := binary.Read(conn, binary.LittleEndian, &size)
        if err != nil {
            fmt.Println("Connection error", err)
            break
        }
        // If no data is received continue to check for data again.
        if size == 0 {
            continue
        }
        fmt.Println("Incoming data size:", size)
        // Copy the received data to a buffer.
        n, err := io.CopyN(buf, conn, size)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println(string(buf.Bytes()))
        fmt.Printf("read %d bytes over the network\n", n)
        // Reset the buffer.
        buf = new(bytes.Buffer)
    }
}

func main() {
    // Create and start a new BvrDam.
    server := &FileServer{}
    server.start()
}
