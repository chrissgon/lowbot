package lowbot

import (
	"fmt"

	"github.com/google/uuid"
)

type RoomGuests map[string]*Guest
type Room struct {
	RoomID       uuid.UUID
	Guests       RoomGuests
	Interactions []*Interaction
	Running      bool
}

func NewRoom(guests RoomGuests) *Room {
	return &Room{
		RoomID:       uuid.New(),
		Guests:       guests,
		Interactions: []*Interaction{},
	}
}

func (room *Room) Start() error {
	for _, guest := range room.Guests {
		err := guest.Channel.Start()

		if err != nil {
			return err
		}

		go room.ListenGuest(guest)
	}

	room.Running = true

	return nil
}

func (room *Room) Stop() error {
	for _, guest := range room.Guests {
		err := guest.Channel.Stop()

		if err != nil {
			return err
		}
	}

	room.Running = false

	return nil
}

func (room *Room) AddGuest(guest *Guest) {
	room.Guests[guest.Who.WhoID] = guest
}

func (room *Room) ListenGuest(guest *Guest) {
	listener := guest.Channel.GetChannel().Broadcast.Listen()

	for interaction := range listener {
		fmt.Println("new interaction", interaction.Parameters.Text, guest.Who.Name)
		err := room.AddInteraction(interaction)

		// TODO: improve how to receive the consumer errors
		if err != nil {
			printLog(fmt.Sprintf("%v: WhoID:<%v> ERR: %v\n", guest.Who.Name, interaction.Sender.WhoID, err))
		}
	}
}

func (room *Room) AddInteraction(interaction *Interaction) error {
	room.Interactions = append(room.Interactions, interaction)

	return room.SendInteractionExcludingSender(interaction)
}

func (room *Room) SendInteractionExcludingSender(interaction *Interaction) (err error) {
	for _, guest := range room.Guests {
		// fmt.Println(interaction.Sender.WhoID, guest.Who.WhoID)
		// if interaction.Destination.WhoID == guest.Who.WhoID {
		// 	continue
		// }

		newInteraction := *interaction
		newInteraction.SetDestination(guest.Who)

		err := SendInteraction(guest.Channel, &newInteraction)

		fmt.Println("passed", guest.Channel.GetChannel().Name, err)
		if err != nil {
			return err
		}
	}

	return nil
}
