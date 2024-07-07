package lowbot

import "github.com/google/uuid"

type RoomGuests map[string]*Guest
type Room struct {
	RoomID       uuid.UUID
	Guests       RoomGuests
	Interactions []*Interaction
}

func NewRoom(guests RoomGuests) *Room {
	room := &Room{
		RoomID:       uuid.New(),
		Guests:       guests,
		Interactions: []*Interaction{},
	}

	for _, guest := range guests{
		room.ListenGuest(guest)
	}

	return room
}

func (room *Room) AddGuest(guest *Guest) {
	room.Guests[guest.Who.WhoID] = guest
	room.ListenGuest(guest)
}

func (room *Room) ListenGuest(guest *Guest) {
	interactions := make(chan *Interaction)

	guest.Channel.Next(interactions)

	for interaction := range interactions {
		room.AddInteraction(interaction)
	}

	close(interactions)
}

func (room *Room) AddInteraction(interaction *Interaction) []error {
	room.Interactions = append(room.Interactions, interaction)

	return room.SendInteractionExcludingSender(interaction)
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
