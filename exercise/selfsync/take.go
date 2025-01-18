package selfsync

// take or counter based quitting pattern shows how to quit conditionally a
// pipeline.

func Take[ChannelType any](quit chan int,
	counter int, input <-chan ChannelType,
) <-chan ChannelType {
	output := make(chan ChannelType)
	go func() {
		defer close(output)
		isOpen := true
		var msg ChannelType
		for counter > 0 && isOpen {
			select {
			case msg, isOpen = (<-input):
				if isOpen {
					output <- msg
					counter--
				}
			case <-quit:
				return
			}
		}
		if counter == 0 {
			close(quit)
		}
	}()
	return output
}

func TakeUntil[ChannelType any](f func(ChannelType) bool,
	quit chan int, input <-chan ChannelType,
) <-chan ChannelType {
	output := make(chan ChannelType)
	go func() {
		isOpen := true
		var msg ChannelType
		for f(msg) && isOpen {
			select {
			case msg, isOpen = (<-input):
				if isOpen {
					output <- msg
				}
			case <-quit:
				return
			}
		}
		if !f(msg) {
			close(output)
		}
	}()
	return output
}
