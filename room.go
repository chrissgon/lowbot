package lowbot

import "github.com/google/uuid"

type Room struct {
	RoomID uuid.UUID
	Interactions []*Interaction
}