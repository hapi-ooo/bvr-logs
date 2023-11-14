# Bvr Logs
Bvr collects logs at the dam to centralize availability.
## Architecture
Bvr Logs is composed of two applications written in Go.
- BvrClient expects a path to a log file when it is run. When it starts, it first 
connects to the server and sends the contents of the file. The file is kept
open, and any additions to the file are read and sent to the server as well.
- BvrDam has an endpoint for BvrClient to send logs to. BvrDam listens on port 
3000 for connections. It first checks the size of the incoming data, then 
writes that specified amount of data to a buffer. At the moment it simply
prints out the received data.
## Usage
## Status
BvrClient always sends the entire file when it starts up. This can lead to redundant
logs if logs are  cleared between starts. BvrClient has the destination URL hardcoded
so the implementation would require edits to deploy. BvrDam does not authenticate
requests in any meaningful way, and it still just prints the data it receives to
standard output.
## ToDo
- JWT Authentication for BvrDam resources
- BvrDam dynamically generate output file
- Flag for BvrClient to skip the initial post request (skips full file)
- Flag for BvrClient to set the BvrDam URL
- Use Zap for logging for better error messages.
- Add unit tests for BvrClient and BvrDam.
- Integrate with messaging architecture to increase reliability (Kafka?).
