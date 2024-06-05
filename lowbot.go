package lowbot

import "sync"

func StartConsumer(consumer IConsumer, channels []IChannel) {
	var wg sync.WaitGroup

	for _, channel := range channels {
		go func(consumer IConsumer, channel IChannel) {
			interactions := make(chan *Interaction)

			go channel.Next(interactions)

			for interaction := range interactions {
				consumer.Run(interaction, channel)
			}

			close(interactions)
		}(consumer, channel)
	}

	wg.Add(1)
	wg.Wait()
}
