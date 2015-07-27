This repo hosts the Go daemon that parses data received from the Rbkit
server gem. The incoming data is parsed and stored/processed for use by
the front-end built in Electron.

Installing pre-requisites
=========================

After you have a working Go installation on your machine, install the
ZMQ Go library, and msgpack library by running:

```
go get github.com/vaughan0/go-zmq
go get github.com/ugorji/go/codec
```



Running the GC test
===================


In the Rbkit repo, run the `experiments/using_rbkit.rb` daemon. This
runs a sample script and has the rbkit server started.


Clone this repo, run `go run trigger_gc.go` to repeatedly
trigger the GC in the running ruby script.


To visualize if the GC is getting triggered or not, add a `puts
GC.stat(:count)` in the `using_rbkit.rb` file.


Running the `using_rbkit.rb` clone
==================================

To run the clone of `using_rbkit.rb` which will log all the event
messges received from the server to the file `/tmp/rbkit_go.log`, run
the following command:

`go run trigger.bak.go`

I still have to work on the naming


Running the command test
========================

To run command test, run the following command:

`go run examples/commander/commander.go`

This will open a prompt on the terminal to that of the
`rbkit_command_test.rb` on the Rbkit repo. This one won't log any output
from the data socket anywhere.
