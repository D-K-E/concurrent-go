package selfsync

import "sync"

// fan in function merges multiple channel outputs into a single one. The
// necessary prerequisite for this is to make sure that we close the common
// output channel after we had done dealing with the incoming channels
func FanIn[ChannelType any](quit <-chan int,
	incomeChannels ...<-chan ChannelType,
) chan ChannelType {
	wg := sync.WaitGroup{}
	wg.Add(len(incomeChannels))
	output := make(chan ChannelType)
	for index, channel := range incomeChannels {
		go func(chnl <-chan ChannelType, indx int) {
			defer wg.Done()
			for j := range chnl {
				select {
				case output <- j:
				case <-quit:
					return
				}
			}
		}(channel, index)
	}

	// now wait for everything to complete
	go func() {
		wg.Wait()
		close(output)
	}()
	return output
}
