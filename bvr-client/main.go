package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Command line parameters
var (
    LogPath string
)
func init() {
    flag.StringVar(&LogPath, "p", "", "Path to the log file")
    flag.Parse()
    // Handle when no LogPath value is passed.
    if LogPath == "" {
        log.Fatal("No value for LogPath provided! Please pass a path to the log file with the '-p' flag")
    }
}

// Interface for the client program --
type BvrClient struct{
    buf *bytes.Buffer
    http *http.Client
}

// Returns a new BvrClient instance --
func newBvrClient(buf *bytes.Buffer, http *http.Client) *BvrClient {
    return &BvrClient{ buf: buf, http: http }
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
// Store the position we last read in the data to only look at new data.
// Send the bytes to the server when data is read.
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
            // go bc.sendBytes(n, bc.buf.Bytes()[startByte:endByte])
            // Create the request body
            body := bytes.NewReader(bc.buf.Bytes()[startByte:endByte])
            // POST the logs
            bc.http.Post("http://localhost:3000/log", "text/plain", body)
            // Update the startByte to the end of what was read.
            startByte = endByte
        }
    }
}

func main() {
    // Create and start a new BvrClient.
    client := newBvrClient(new(bytes.Buffer), &http.Client{})
    client.start()
}
