# Bvr Logs
Bvr collects logs at the dam to centralize availability.
## Architecture
bvr logs is composed of two applications written in Go.
- BvrClient expects a path to a log file when it is run. When it starts, it first 
connects to the server and sends the contents of the file. The file is kept
open, and any additions to the file are read and sent to the server as well.
bvr client consumes logs to send to bvr dam. Each time data is sent to the 
server, the size of what is being sent is first written to the connection.
- BvrDam has an endpoint for BvrClient to send logs to. BvrDam listens on port 
3000 for connections. It first checks the size of the incoming data, then 
writes the specified amount to a buffer. At the moment it simply prints out the
received data.
## Usage
## Status
BvrClient and BvrDam are using a basic TCP connection from the standard net
package. This is not safe for any production environments because the
communication is unencrypted. Additionally, there is no recourse for BvrClient
if BvrDam is not online to handle data sent to it.
## ToDo
- Update to using net/http instead of net to increase security and reliability.
- Use Zap for logging for better error messages.
- Add unit tests for BvrClient and BvrDam.
- Integrate with messaging architecture to increase reliability (Kafka?).
