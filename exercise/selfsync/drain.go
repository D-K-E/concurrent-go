package selfsync


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
