package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/digitalocean/graphql-api/graph/model"
)

// ChatRoomMessages is the resolver for the chatRoomMessages field.
func (r *chatRoomResolver) ChatRoomMessages(ctx context.Context, obj *model.ChatRoom, first *int, after *string) (*model.ChatRoomMessagesConnection, error) {

	// The cursor is base64 encoded by convention, so we need to decode it first
	var decodedCursor string
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, err
		}
		decodedCursor = string(b)
	}

	// Here we could query the DB to get data, e.g.
	// SELECT * FROM messages WHERE chat_room_id = obj.ID AND timestamp < decodedCursor
	edges := make([]*model.ChatRoomMessagesEdge, *first)
	count := 0
	currentPage := false

	// If no cursor present start from the top
	if decodedCursor == "" {
		currentPage = true
	}
	hasNextPage := false

	// Iterating over the mocked messages to find the current page
	// In real world use-case you should fetch only the required
	// part of data from the database
	for i, v := range r.chatRoomMessages[obj.ID] {
		if v.ID == decodedCursor {
			currentPage = true
		}

		if currentPage && count < *first {
			edges[count] = &model.ChatRoomMessagesEdge{
				Cursor: base64.StdEncoding.EncodeToString([]byte(v.ID)),
				Node:   &v,
			}
			count++
		}

		// If there are any elements left after the current page
		// we indicate that in the response
		if count == *first && i < len(r.chatRoomMessages[obj.ID]) {
			hasNextPage = true
		}
	}

	pageInfo := model.PageInfo{
		StartCursor: mustStringPtr(base64.StdEncoding.EncodeToString([]byte(edges[0].Node.ID))),
		EndCursor:   mustStringPtr(base64.StdEncoding.EncodeToString([]byte(edges[count-1].Node.ID))),
		HasNextPage: hasNextPage,
	}

	mc := model.ChatRoomMessagesConnection{
		Edges:    edges[:count],
		PageInfo: &pageInfo,
	}

	return &mc, nil
}

// AllChatRooms is the resolver for the allChatRooms field.
func (r *queryResolver) AllChatRooms(ctx context.Context) ([]*model.ChatRoom, error) {
	chatRooms := make([]*model.ChatRoom, len(r.chatRooms))
	idx := 0
	for _, chatRoom := range r.chatRooms {
		chatRooms[idx] = &chatRoom
		idx = idx + 1
	}
	return chatRooms, nil
}

// ChatRoom is the resolver for the chatRoom field.
func (r *queryResolver) ChatRoom(ctx context.Context, id string) (*model.ChatRoom, error) {
	if t, ok := r.chatRooms[id]; ok {
		return &t, nil
	}
	return nil, errors.New("chat room not found")
}

// ChatRoom returns ChatRoomResolver implementation.
func (r *Resolver) ChatRoom() ChatRoomResolver { return &chatRoomResolver{r} }

type chatRoomResolver struct{ *Resolver }
