package lowbot

type IConsumer interface{
	Run(*Interaction, IChannel)
}

