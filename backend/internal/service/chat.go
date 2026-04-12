package service

import (
	"fmt"
	"slices"
	"time"

	"hybrid-app/backend/internal/domain"
	"hybrid-app/backend/internal/repository"
)

type ChatService struct{ baseService }

func NewChatService(repo repository.Repository, idGen func(string) string) *ChatService {
	return &ChatService{baseService{repo: repo, newID: idGen}}
}

func (s *ChatService) Chats(userID string) ([]map[string]any, error) {
	chats, err := s.repo.ListChatsForUser(userID)
	if err != nil {
		return nil, err
	}
	var result []map[string]any
	for _, chat := range chats {
		otherID := chat.ParticipantIDs[0]
		if otherID == userID && len(chat.ParticipantIDs) > 1 {
			otherID = chat.ParticipantIDs[1]
		}
		participant, err := s.repo.FindUserByID(otherID)
		if err != nil {
			return nil, err
		}
		result = append(result, map[string]any{
			"id":           chat.ID,
			"matchId":      chat.MatchID,
			"participant":  participant,
			"lastMessage":  chat.LastMessage,
			"lastActiveAt": chat.LastActiveAt,
			"unreadCount":  chat.UnreadCount,
			"pinned":       chat.Pinned,
		})
	}
	slices.SortFunc(result, func(a, b map[string]any) int {
		at := a["lastActiveAt"].(time.Time)
		bt := b["lastActiveAt"].(time.Time)
		if at.After(bt) {
			return -1
		}
		if at.Before(bt) {
			return 1
		}
		return 0
	})
	return result, nil
}

func (s *ChatService) ChatMessages(userID, chatID string) ([]domain.ChatMessage, error) {
	chat, err := s.repo.FindChatByID(chatID)
	if err != nil {
		return nil, err
	}
	if chat == nil || !slices.Contains(chat.ParticipantIDs, userID) {
		return nil, ErrChatNotFound
	}
	return s.repo.ListMessages(chatID)
}

func (s *ChatService) SendMessage(userID, chatID, messageType, text, giftID string) (*domain.ChatMessage, error) {
	chat, err := s.repo.FindChatByID(chatID)
	if err != nil {
		return nil, err
	}
	if chat == nil || !slices.Contains(chat.ParticipantIDs, userID) {
		return nil, ErrChatNotFound
	}

	if messageType == "gift" {
		gifts, err := s.repo.ListGifts()
		if err != nil {
			return nil, err
		}
		var selected *domain.Gift
		for _, gift := range gifts {
			if gift.ID == giftID {
				g := gift
				selected = &g
				break
			}
		}
		if selected == nil {
			return nil, fmt.Errorf("gift not found")
		}
		user, err := s.repo.FindUserByID(userID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, ErrUserNotFound
		}
		if user.SparksBalance < selected.CostSparks {
			return nil, ErrNotEnoughSparks
		}
		user.SparksBalance -= selected.CostSparks
		user.UpdatedAt = time.Now().UTC()
		if err := s.repo.SaveUser(user); err != nil {
			return nil, err
		}
		if _, err := addTransaction(s.repo, s.newID, user.ID, "gift_send", -selected.CostSparks, fmt.Sprintf("Sent %s", selected.Name)); err != nil {
			return nil, err
		}
		if text == "" {
			text = fmt.Sprintf("Sent a %s", selected.Name)
		}
	}

	msg := domain.ChatMessage{
		ID:        s.newID("msg"),
		ChatID:    chatID,
		SenderID:  userID,
		Type:      messageType,
		Text:      text,
		GiftID:    giftID,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.AddMessage(msg); err != nil {
		return nil, err
	}
	chat.LastMessage = text
	chat.LastActiveAt = msg.CreatedAt
	chat.UnreadCount++
	if err := s.repo.SaveChat(chat); err != nil {
		return nil, err
	}
	return &msg, nil
}
