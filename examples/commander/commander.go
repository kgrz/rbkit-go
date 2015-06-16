// Prompt user for commands, print them out. This has to be done in a
// separate thread
package main

import "fmt"
import "github.com/code-mancers/rbkit-go/receiver"

var optionDict = map[int]string{
	1: "start_memory_profile",
	2: "stop_memory_profile",
	3: "objectspace_snapshot",
	4: "trigger_gc",
	5: "handshake",
}

func askForOption(option chan int) {
	var input int
	fmt.Println("Enter a selection for the command: ")
	fmt.Println("1. Start Memory Profile")
	fmt.Println("2. Stop Memory Profile")
	fmt.Println("3. Objectspace Snapshot")
	fmt.Println("4. Trigger GC")
	fmt.Println("5. Handshake")
	fmt.Println("Hit Ctrl+C to stop")
	fmt.Scanln(&input)
	option <- input
}

func getOption(option chan int) {
	enteredValue := <-option
	if enteredValue < 1 || enteredValue > 5 {
		fmt.Println("\nInvalid option\n\n")
	} else {
		fmt.Printf("\nThe option you've selected is %d\n", enteredValue)
		fmt.Printf("\nOr, in other words: %s\n\n", optionDict[enteredValue])
	}
}

func main() {
	option := make(chan int)
	// Repeatedly ask for an option
	for {
		go askForOption(option)
		getOption(option)
		receiver.Receive()
	}
}
