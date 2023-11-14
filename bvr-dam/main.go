package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var (
    // Command line parameters
    port string
)

func init() {
    // Process command line arguments
    flag.StringVar(&port, "p", ":3010", "Address to host the server " +
    "at. The default value is \":3010\".")
    flag.Parse()
}

// Interface for the server program --
type BvrDam struct{
    server *http.Server
}

// Returns a new BvrDam instance
func newBvrDam(server *http.Server) *BvrDam {
    return &BvrDam{ server: server }
}

// Handler that writes data to a file --
func (fs *BvrDam) logPostHandler(w http.ResponseWriter, r *http.Request) {
    // Return a 405 Method Not Allowed on anything but a POST request.
    if r.Method != "POST" {
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
    // Get the content length so we know how many bytes to copy
    contentLength, err := strconv.Atoi(r.Header.Get("Content-Length"))
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
    // Get a buffer of the data in the request body
    buf, err := fs.getRequestData(r.Body, int64(contentLength))
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
    }
    fmt.Println("Data received:\n", string(buf.Bytes()))
    // TODO Write the data to the corresponding log
}

// Read data from the request body --
// Return a buffer containing the data or an error.
func (fs *BvrDam) getRequestData(body io.ReadCloser, contentLength int64) (*bytes.Buffer, error) {
    // Initialize buffer to read body into
    buf := new(bytes.Buffer)
    // Copy the received data to a buffer.
    n, err := io.CopyN(buf, body, contentLength)
    if err != nil {
        return nil, err
    }
    fmt.Printf("%d bytes were read over the network\n", n)
    // Return the buffer of data from the request body.
    return buf, nil
}

func main() {
    // Create a new BvrDam.
    dam := newBvrDam(&http.Server { Addr: port })

    // Register the endpoint handlers
    http.HandleFunc("/log", dam.logPostHandler)
    http.HandleFunc("/", http.NotFound)

    // Listen for requests
    fmt.Printf("BvrDam running at localhost%s\n", port)
    log.Fatal(dam.server.ListenAndServe())
}
