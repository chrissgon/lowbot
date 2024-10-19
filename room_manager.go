package lowbot

import (
	"fmt"

	"github.com/google/uuid"
)

type GuestRoomRelation map[string]uuid.UUID

type RoomManager struct {
	ChannelGuestRoomRelation map[uuid.UUID]GuestRoomRelation
	Rooms                    map[uuid.UUID]*Room
	Err error
}

var roomManager = NewRoomManager()

func NewRoomManager() *RoomManager {
	return &RoomManager{
		ChannelGuestRoomRelation: map[uuid.UUID]GuestRoomRelation{},
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
		channelID := manager.AddChannel(guest.Channel)
		manager.ChannelGuestRoomRelation[channelID][guest.Who.WhoID] = room.RoomID
		manager.Rooms[room.RoomID] = room
	}

	return nil
}
func (manager *RoomManager) AddChannel(channel IChannel) uuid.UUID {
	channelID := channel.GetChannel().ChannelID
	_, exists := manager.ChannelGuestRoomRelation[channelID]

	if !exists {
		manager.ChannelGuestRoomRelation[channelID] = GuestRoomRelation{}
		go manager.ListenChannel(channel)
	}

	return channelID
}

func (manager *RoomManager) ListenChannel(channel IChannel) {
	listener := channel.GetChannel().Broadcast.Listen()

	for interaction := range listener {
		err := manager.AddInteraction(interaction)

		if err != nil {
			manager.Err = err
			printLog(fmt.Sprintf("Room Manager: ListenChannel ERR: %v", err))
		}
	}
}

func (manager *RoomManager) AddInteraction(interaction *Interaction) error {
	// channelID := interaction.Channel.ChannelID
	// guestID := interaction.Sender.WhoID

	// roomID, exists := manager.ChannelGuestRoomRelation[channelID][guestID]

	// if !exists {
	// 	return ERR_UNKNOWN_ROOM
	// }

	// return manager.Rooms[roomID].AddInteraction(interaction)
	return nil
}

func (manager *RoomManager) AddGuest(roomID uuid.UUID, guest *Guest) error {
	room := manager.GetRoom(roomID)

	if room == nil {
		return ERR_UNKNOWN_ROOM
	}

	room.AddGuest(guest)

	channelID := manager.AddChannel(guest.Channel)

	manager.ChannelGuestRoomRelation[channelID][guest.Who.WhoID] = room.RoomID

	return nil
}
