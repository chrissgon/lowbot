package lowbot

import (
	"github.com/google/uuid"
)

type RoomGuests map[string]*Guest
type Room struct {
	RoomID       uuid.UUID
	Guests       RoomGuests
	Interactions []*Interaction
}

func NewRoom(guests RoomGuests) *Room {
	return &Room{
		RoomID:       uuid.New(),
		Guests:       guests,
		Interactions: []*Interaction{},
	}
}

func (room *Room) AddGuest(guest *Guest) {
	room.Guests[guest.Who.WhoID] = guest
}

func (room *Room) AddInteraction(interaction *Interaction) error {
	room.Interactions = append(room.Interactions, interaction)

	errors := room.SendInteractionExcludingSender(interaction)

	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

func (room *Room) SendInteractionExcludingSender(interaction *Interaction) (errs []error) {
	for _, guest := range room.Guests {
		if interaction.Sender.WhoID == guest.Who.WhoID {
			continue
		}

		newInteraction := *interaction
		newInteraction.SetDestination(guest.Who)

		err := SendInteraction(guest.Channel, &newInteraction)

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}
