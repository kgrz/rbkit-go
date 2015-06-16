This repo hosts the Go daemon that parses data received from the Rbkit
server gem. The incoming data is parsed and stored/processed for use by
the front-end built in Electron.

The communication from the Electron app to the Go app is achieved via
IPC.
