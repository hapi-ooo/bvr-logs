package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// command line parameters
var (
    LogPath string
)

func init() {
    flag.StringVar(&LogPath, "p", "", "Path to the log file")
    flag.Parse()
}

// Interface for the client program --
type BvrClient struct{
    buf *bytes.Buffer
    conn net.Conn
}

// Returns a new BvrClient instance --
func newBvrClient(buf *bytes.Buffer, conn net.Conn) *BvrClient {
    return &BvrClient{ buf: buf, conn: conn }
}

// Starts the client --
// Open the log file, read from the log, and write to server.
func (bc *BvrClient) start() {
    // Open the log file provided by the user.
    file, err := os.Open(LogPath)
    if err != nil {
        log.Fatal(err)
    }
    // When this function returns the file should be closed.
    defer file.Close()

    // Read the file indefinitely
    // Write its contents to the connection as well as any additions
    // to the file.
    bc.readForever(file)
}

// Read the data and writes to the server -- 
// Read the data until an error occurs.
// Store the position we last read in the data to only look at new data
// Send the bytes to the server when data is read
func (bc *BvrClient) readForever(r io.Reader) {
    // When new bytes are available, begin reading at startByte
    // initialized to the beginning of the byte slice.
    startByte := int64(0)
    for {
        // Crash on an error reading the data source because
        // io.EOF is not returned as an error from (*Buffer) ReadFrom.
        n, err := bc.buf.ReadFrom(r)
        if err != nil {
            fmt.Println(err)
            break
        }
        // Only if bytes were read:
        if n != 0 {
            // Read the number of bytes past the startByte.
            endByte := n + startByte
            go bc.sendBytes(n, bc.buf.Bytes()[startByte:endByte])
            // Update the startByte to the end of what was read.
            startByte = endByte
        }
    }
}

// Send bytes to the server --
// Write the size of the data we intend to send to the connection.
// Write the data to the connection.
// Returns an error if one occurred.
func (bc *BvrClient) sendBytes(size int64, data []byte) error {
    // First write the size of the data being sent.
    err := binary.Write(bc.conn, binary.LittleEndian, size)
    // Write the data we wish to send to the server.
    n, err := io.CopyN(bc.conn, bytes.NewReader(data), size)
    if err != nil {
        return err
    }
    fmt.Printf("%d bytes have been written over the network.\n", n)
    return nil
}

func main() {
    // Handle when no LogPath value is passed.
    if LogPath == "" {
        log.Fatal("No value for LogPath provided! Please pass a path to the log file with the '-p' flag")
    }

    // Create a new connection to pass to the client.
    conn, err := net.Dial("tcp", ":3000")
    defer conn.Close()
    if err != nil {
        log.Fatal(err)
    }

    // Create and start a new BvrClient.
    client := newBvrClient(new(bytes.Buffer), conn)
    client.start()
}
