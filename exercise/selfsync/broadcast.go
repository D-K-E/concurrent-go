package selfsync

// broadcasting or worker pattern replicates the messages for multiple output
// channels. It is the opposite of FanIn pattern where we merge multiple
// inputs to single output

// Create multiple channels
func CreateAll[ChannelType any](quit <-chan int,
	input <-chan ChannelType, n int,
) []chan ChannelType {
	channels := make([]chan ChannelType, n)
	for i := range channels {
		channels[i] = make(chan ChannelType)
	}
	return channels
}

// close multiple channels
func CloseAll[ChannelType any](channels ...chan ChannelType) {
	for _, out := range channels {
		close(out)
	}
}

// broadcast input channel to multiple output channels
func Broadcast[ChannelType any](quit <-chan int,
	input <-chan ChannelType, n int,
) []chan ChannelType {
	outputs := CreateAll[ChannelType](quit, input, n)
	go func() {
		defer CloseAll(outputs...)
		var channelMsg ChannelType
		isOpen := true
		for isOpen {
			select {
			case channelMsg, isOpen = (<-input):
				if isOpen {
					for _, out := range outputs {
						out <- channelMsg
					}
				}
			case <-quit:
				return
			}
		}
	}()
	return outputs
}
