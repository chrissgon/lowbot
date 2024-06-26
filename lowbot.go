package lowbot

import (
	"fmt"
	"sync"
)

func StartConsumer(consumer IConsumer, channels []IChannel) {
	var wg sync.WaitGroup

	for _, channel := range channels {
		go func(consumer IConsumer, channel IChannel) {
			interactions := make(chan *Interaction)

			go channel.Next(interactions)

			for interaction := range interactions {
				err := consumer.Run(interaction, channel)

				if err != nil {
					printLog(fmt.Sprintf("%v: WhoID:<%v> ERR: %v\n", consumer.GetConsumer().Name, interaction.Sender.WhoID, err))
				}
			}

			close(interactions)
		}(consumer, channel)
	}

	wg.Add(1)
	wg.Wait()
}
