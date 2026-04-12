package http

import (
	"net/http"
	"strings"
)

func (s *Server) handleChats(w http.ResponseWriter, _ *http.Request, userID string) {
	chats, err := s.services.Chat.Chats(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"chats": chats})
}

func (s *Server) handleChatMessages(w http.ResponseWriter, r *http.Request, userID string) {
	chatID, ok := chatIDFromPath(r.URL.Path)
	if !ok || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	messages, err := s.services.Chat.ChatMessages(userID, chatID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"messages": messages})
}

func (s *Server) handleSendMessage(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	chatID, ok := chatIDFromPath(strings.TrimSuffix(r.URL.Path, "/messages"))
	if !ok || !strings.HasSuffix(r.URL.Path, "/messages") {
		http.NotFound(w, r)
		return
	}

	var input struct {
		Type   string `json:"type"`
		Text   string `json:"text"`
		GiftID string `json:"giftId"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if input.Type == "" {
		input.Type = "text"
	}

	message, err := s.services.Chat.SendMessage(userID, chatID, input.Type, input.Text, input.GiftID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, message)
}

func chatIDFromPath(path string) (string, bool) {
	prefix := "/api/v1/chats/"
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}
	rest := strings.TrimPrefix(path, prefix)
	parts := strings.Split(strings.Trim(rest, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		return "", false
	}
	return parts[0], true
}
