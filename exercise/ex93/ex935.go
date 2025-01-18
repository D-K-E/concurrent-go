package ex93

import (
	"fmt"

	selfsync "github.com/D-K-E/concurrent-go/selfsync"
)

/*
Connect the components developed in exercises 1 to 4 together in a main()
function using the following pseudocode:
*/

func Print[ChannelType any](quit <-chan int,
	input <-chan ChannelType,
) <-chan ChannelType {
	output := make(chan ChannelType)
	go func() {
		defer close(output)
		var msg ChannelType
		isOpen := true
		for isOpen {
			select {
			case msg, isOpen = (<-input):
				if isOpen {
					fmt.Println(msg)
				}
			case <-quit:
				return
			}
		}
	}()
	return output
}

func Drain[ChannelType any](quit <-chan int, input <-chan ChannelType) <-chan ChannelType {
	//
	output := make(chan ChannelType)
	go func() {
		defer close(output)
		var msg ChannelType
		isOpen := true
		for isOpen {
			select {
			case msg, isOpen = (<-input):
				_ = msg
			case <-quit:
				return
			}
		}
	}()
	return output
}

func Ex935Main() {
	quit := make(chan int)
	defer close(quit)
	squares := GenerateSquares(quit)
	taken := selfsync.TakeUntil(func(s int) bool { return s <= 10000 },
		quit, squares)
	//
	printed := Print(quit, taken)
	<-Drain(quit, printed)
}
