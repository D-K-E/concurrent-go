package selfsync

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
