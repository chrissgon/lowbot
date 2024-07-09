package lowbot

import (
	"fmt"

	"github.com/google/uuid"
)

// type GuestRoomRelation

type RoomManager struct {
	ChannelGuestRoomRelation map[uuid.UUID]map[string]uuid.UUID
	Rooms                    map[uuid.UUID]*Room
	// Channels                 []IChannel
}

var roomManager = NewRoomManager()

func NewRoomManager() *RoomManager {
	return &RoomManager{
		ChannelGuestRoomRelation: map[uuid.UUID]map[string]uuid.UUID{},
		Rooms:                    map[uuid.UUID]*Room{},
	}
}

func GetRoomManager() *RoomManager {
	return roomManager
}

func (manager *RoomManager) GetRoom(roomID uuid.UUID) *Room {
	return manager.Rooms[roomID]
}

func (manager *RoomManager) AddRoom(room *Room) error {
	for _, guest := range room.Guests {
		fmt.Println("guest", guest.Who.WhoID)
		channelID := guest.Channel.GetChannel().ChannelID
		_, exists := manager.ChannelGuestRoomRelation[channelID]

		if !exists {
			manager.ChannelGuestRoomRelation[channelID] = map[string]uuid.UUID{}
			// manager.Channels = append(manager.Channels, guest.Channel)
			go manager.ListenChannel(guest.Channel)
		}

		manager.ChannelGuestRoomRelation[channelID][guest.Who.WhoID] = room.RoomID
		manager.Rooms[room.RoomID] = room
	}

	return nil
}

func (manager *RoomManager) ListenChannel(channel IChannel) {
	listener := channel.GetChannel().Broadcast.Listen()

	for interaction := range listener {
		manager.AddInteraction(interaction)
	}

	channel.GetChannel().Broadcast.Close()
}

func (manager *RoomManager) AddInteraction(interaction *Interaction) error {
	channelID := interaction.Channel.ChannelID
	guestID := interaction.Sender.WhoID

	roomID, exists := manager.ChannelGuestRoomRelation[channelID][guestID]

	if !exists {
		return ERR_UNKNOWN_ROOM
	}

	return manager.Rooms[roomID].AddInteraction(interaction)
}

func (manager *RoomManager) AddGuest(roomID uuid.UUID, guest *Guest) error {
	room := manager.GetRoom(roomID)

	if room == nil {
		return ERR_UNKNOWN_ROOM
	}

	room.AddGuest(guest)

	channelID := guest.Channel.GetChannel().ChannelID

	_, exists := manager.ChannelGuestRoomRelation[channelID]

	if !exists {
		manager.ChannelGuestRoomRelation[channelID] = map[string]uuid.UUID{}
		go manager.ListenChannel(guest.Channel)
	}

	manager.ChannelGuestRoomRelation[channelID][guest.Who.WhoID] = room.RoomID

	return nil
}
