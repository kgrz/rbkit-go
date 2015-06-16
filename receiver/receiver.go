package receiver

import (
	"fmt"
	"github.com/vaughan0/go-zmq"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func createAndBindSocket() *zmq.Socket {
	ctx, err := zmq.NewContext()
	checkError(err)
	defer ctx.Close()

	sock, err := ctx.Socket(zmq.Sub)
	checkError(err)
	defer sock.Close()

	e := sock.Bind("tcp://127.0.0.1:5555")
	checkError(e)

	return sock
}

func Receive() {
	sock := createAndBindSocket()
	for {
		parts, err := sock.Recv()
		checkError(err)
		fmt.Sprintf("Recieved %d message parts", len(parts))
	}
}
