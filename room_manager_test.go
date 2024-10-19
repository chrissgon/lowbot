package lowbot

// func TestRoomManager_NewRoomManager(t *testing.T) {
// 	have := NewRoomManager()
// 	expect := &RoomManager{
// 		ChannelGuestRoomRelation: map[uuid.UUID]GuestRoomRelation{},
// 		Rooms:                    map[uuid.UUID]*Room{},
// 	}

// 	if !reflect.DeepEqual(expect, have) {
// 		t.Errorf(FormatTestError(expect, have))
// 	}
// }

// func TestRoomManager_GetRoomManager(t *testing.T) {
// 	have := GetRoomManager()
// 	expect := roomManager

// 	if !reflect.DeepEqual(expect, have) {
// 		t.Errorf(FormatTestError(expect, have))
// 	}
// }

// func TestRoomManager_AddChannel(t *testing.T) {
// 	var have any
// 	var expect any

// 	CHANNEL_MOCK = newMockChannel()

// 	manager := NewRoomManager()

// 	channelID := CHANNEL_MOCK.GetChannel().ChannelID
// 	_, exists := manager.ChannelGuestRoomRelation[channelID]

// 	if exists {
// 		t.Errorf(FormatTestError(false, exists))
// 	}

// 	have = manager.AddChannel(CHANNEL_MOCK)
// 	expect = CHANNEL_MOCK.GetChannel().ChannelID

// 	if expect != have {
// 		t.Errorf(FormatTestError(expect, have))
// 	}

// 	_, exists = manager.ChannelGuestRoomRelation[channelID]

// 	if !exists {
// 		t.Errorf(FormatTestError(true, exists))
// 	}

// 	have = manager.ChannelGuestRoomRelation[channelID]
// 	expect = GuestRoomRelation{}

// 	if !reflect.DeepEqual(expect, have) {
// 		t.Errorf(FormatTestError(expect, have))
// 	}

// 	time.Sleep(1 * time.Millisecond)

// 	have = len(CHANNEL_MOCK.GetChannel().Broadcast.listeners)
// 	expect = 1

// 	if expect != have {
// 		t.Errorf(FormatTestError(expect, have))
// 	}
// }

// func TestRoomManager_ListenChannel(t *testing.T) {
// 	CHANNEL_MOCK = newMockChannel()
// 	manager := NewRoomManager()
// 	go manager.ListenChannel(CHANNEL_MOCK)

// 	time.Sleep(1 * time.Millisecond)

// 	have := len(CHANNEL_MOCK.GetChannel().Broadcast.listeners)
// 	expect := 1

// 	if expect != have {
// 		t.Errorf(FormatTestError(expect, have))
// 	}

// 	interaction := NewInteractionMessageText(WHO_MOCK, WHO_MOCK, TEXT)
// 	CHANNEL_MOCK.GetChannel().Broadcast.Send(interaction)

// 	time.Sleep(1 * time.Millisecond)

// 	if !errors.Is(manager.Err, ERR_UNKNOWN_ROOM) {
// 		t.Errorf(FormatTestError(ERR_UNKNOWN_ROOM, manager.Err))
// 	}
// 	room := NewRoom(RoomGuests{
// 		WHO_MOCK.WhoID: NewGuest(WHO_MOCK, CHANNEL_MOCK),
// 	})
// 	manager.AddRoom(room)

// 	CHANNEL_MOCK.GetChannel().Broadcast.Send(interaction)

// 	time.Sleep(1 * time.Millisecond)

// 	have = len(manager.GetRoom(room.RoomID).Interactions)
// 	expect = 1

// 	if expect != have {
// 		t.Errorf(FormatTestError(expect, have))
// 	}
// }

// func TestRoomManager_AddRoom(t *testing.T) {
// 	CHANNEL_MOCK = newMockChannel()
// 	manager := NewRoomManager()

// 	room := NewRoom(RoomGuests{
// 		WHO_MOCK.WhoID: NewGuest(WHO_MOCK, CHANNEL_MOCK),
// 	})

// 	err := manager.AddRoom(room)

// 	if err != nil {
// 		t.Errorf(FormatTestError(nil, err))
// 	}

// 	_, exists := manager.Rooms[room.RoomID]

// 	if !exists {
// 		t.Errorf(FormatTestError(true, exists))
// 	}

// 	channelID := CHANNEL_MOCK.GetChannel().ChannelID
// 	relation, exists := manager.ChannelGuestRoomRelation[channelID]

// 	if !exists {
// 		t.Errorf(FormatTestError(true, exists))
// 	}

// 	roomID, exists := relation[WHO_MOCK.WhoID]

// 	if !exists {
// 		t.Errorf(FormatTestError(true, exists))
// 	}

// 	expect := room.RoomID
// 	have := roomID

// 	if expect != have {
// 		t.Errorf(FormatTestError(expect, have))
// 	}
// }

// func TestRoomManager_AddInteraction(t *testing.T) {
// 	CHANNEL_MOCK = newMockChannel()
// 	manager := NewRoomManager()

// 	interaction := NewInteractionMessageText( WHO_MOCK, WHO_MOCK, TEXT)

// 	err := manager.AddInteraction(interaction)

// 	if !errors.Is(err, ERR_UNKNOWN_ROOM) {
// 		t.Errorf(FormatTestError(ERR_UNKNOWN_ROOM, err))
// 	}

// 	room := NewRoom(RoomGuests{
// 		WHO_MOCK.WhoID: NewGuest(WHO_MOCK, CHANNEL_MOCK),
// 	})
// 	manager.AddRoom(room)

// 	err = manager.AddInteraction(interaction)

// 	if err != nil {
// 		t.Errorf(FormatTestError(nil, err))
// 	}

// 	have := len(manager.Rooms[room.RoomID].Interactions)
// 	expect := 1

// 	if expect != have {
// 		t.Errorf(FormatTestError(expect, have))
// 	}
// }

// func TestRoomManager_GetRoom(t *testing.T) {
// 	CHANNEL_MOCK = newMockChannel()
// 	manager := NewRoomManager()

// 	room := NewRoom(RoomGuests{
// 		WHO_MOCK.WhoID: NewGuest(WHO_MOCK, CHANNEL_MOCK),
// 	})

// 	manager.AddRoom(room)

// 	have := manager.GetRoom(room.RoomID)
// 	expect := room

// 	if !reflect.DeepEqual(expect, have) {
// 		t.Errorf(FormatTestError(expect, have))
// 	}
// }

// func TestRoomManager_AddGuest(t *testing.T) {
// 	CHANNEL_MOCK = newMockChannel()
// 	manager := NewRoomManager()

// 	room := NewRoom(RoomGuests{
// 		WHO_MOCK.WhoID: NewGuest(WHO_MOCK, CHANNEL_MOCK),
// 	})
// 	guest := NewGuest(NewWho("2", "amanda"), CHANNEL_MOCK)

// 	err := manager.AddGuest(uuid.New(), guest)

// 	if !errors.Is(err, ERR_UNKNOWN_ROOM) {
// 		t.Errorf(FormatTestError(ERR_UNKNOWN_ROOM, err))
// 	}

// 	manager.AddRoom(room)
// 	err = manager.AddGuest(room.RoomID, guest)

// 	if err != nil {
// 		t.Errorf(FormatTestError(nil, err))
// 	}

// 	have := len(manager.GetRoom(room.RoomID).Guests)
// 	expect := 2

// 	if expect != have {
// 		t.Errorf(FormatTestError(expect, have))
// 	}

// 	channelID := CHANNEL_MOCK.GetChannel().ChannelID
// 	relation, exists := manager.ChannelGuestRoomRelation[channelID]

// 	if !exists {
// 		t.Errorf(FormatTestError(true, exists))
// 	}

// 	roomID, exists := relation[guest.Who.WhoID]

// 	if !exists {
// 		t.Errorf(FormatTestError(true, exists))
// 	}

// 	if roomID != room.RoomID {
// 		t.Errorf(FormatTestError(room.RoomID, roomID))
// 	}
// }
