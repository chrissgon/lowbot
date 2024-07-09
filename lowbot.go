package lowbot

import (
	"fmt"
	"sync"
)

func StartConsumer(consumer IConsumer, channels []IChannel) {
	var wg sync.WaitGroup

	for _, channel := range channels {
		go func(consumer IConsumer, channel IChannel) {
			listener := channel.GetChannel().Broadcast.Listen()

			for interaction := range listener {
				err := consumer.Run(interaction, channel)
				
				if err != nil {
					printLog(fmt.Sprintf("%v: WhoID:<%v> ERR: %v\n", consumer.GetConsumer().Name, interaction.Sender.WhoID, err))
				}
			}

			channel.GetChannel().Broadcast.Close()
		}(consumer, channel)
	}

	wg.Add(1)
	wg.Wait()
}
