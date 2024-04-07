package graph

import (
	"fmt"
	"strconv"

	"github.com/digitalocean/graphql-api/graph/model"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	todos            []*model.Todo
	chatRooms        map[string]model.ChatRoom
	chatRoomMessages map[string][]model.ChatRoomMessage
}

func NewResolver() Config {
	const nChatRooms = 20
	const nMessagesPerChatRoom = 100
	r := Resolver{}
	r.chatRooms = make(map[string]model.ChatRoom, nChatRooms)
	r.chatRoomMessages = make(map[string][]model.ChatRoomMessage, nChatRooms)

	for i := 0; i < nChatRooms; i++ {
		id := strconv.Itoa(i + 1)
		mockChatRoom := model.ChatRoom{
			ID:   id,
			Name: fmt.Sprintf("ChatRoom %d", i),
		}
		r.chatRooms[id] = mockChatRoom
		r.chatRoomMessages[id] = make([]model.ChatRoomMessage, nMessagesPerChatRoom)

		// Generate messages for the ChatRoom
		for k := 0; k < nMessagesPerChatRoom; k++ {
			msgId := strconv.Itoa(k + 1)
			text := fmt.Sprintf("Message %d", k)

			mockMessage := model.ChatRoomMessage{
				ID:   msgId,
				Text: &text,
			}

			r.chatRoomMessages[id][k] = mockMessage
		}
	}

	return Config{
		Resolvers: &r,
	}
}
