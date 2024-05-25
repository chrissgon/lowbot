package lowbot

func StartConsumer(consumer IConsumer, channel IChannel) {
	interactions := make(chan *Interaction)

	go channel.Next(interactions)

	for interaction := range interactions {
		consumer.Run(interaction, channel)
	}
}
