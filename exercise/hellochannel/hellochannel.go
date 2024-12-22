package hellochannel

// very simple channel example
import "fmt"

func receiverString(messages chan string) {
	msg := ""
	for msg != "STOP" {
		msg = (<-messages) /* the operator `<-` fetches a string from the
				   string channel. The paranthesis are not required but it makes it
		           similar to (*x) operator where we dereference a pointer
		*/
		fmt.Println("received", msg, "!")
	}
}

func HelloChannelMain() {
	//
	msgChannel := make(chan string)
	go receiverString(msgChannel)
	fmt.Println("Sending HELLO")
	msgChannel <- "HELLO"

	fmt.Println("Sending ONE")
	msgChannel <- "ONE"

	fmt.Println("Sending THREE")
	msgChannel <- "THREE"

	fmt.Println("Sending STOP")
	msgChannel <- "STOP"
}
