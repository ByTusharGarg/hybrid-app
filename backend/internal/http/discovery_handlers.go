package http

import (
	"net/http"
	"strconv"
)

func (s *Server) handleDiscover(w http.ResponseWriter, _ *http.Request, userID string) {
	topPicks, feed, err := s.services.Discovery.Discover(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"catalog": map[string]any{
			"topPicks": topPicks,
			"profiles": feed,
		},
	})
}

func (s *Server) handleDiscoverAction(w http.ResponseWriter, r *http.Request, userID string) {
	var input struct {
		TargetUserID string `json:"targetUserId"`
		Action       string `json:"action"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	result, err := s.services.Discovery.ApplyAction(userID, input.TargetUserID, input.Action)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleMatches(w http.ResponseWriter, _ *http.Request, userID string) {
	matches, err := s.services.Discovery.Matches(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"matches": matches})
}

func (s *Server) handleLikesYou(w http.ResponseWriter, r *http.Request, userID string) {
	reveal, _ := strconv.ParseBool(r.URL.Query().Get("reveal"))
	likes, err := s.services.Discovery.LikesYou(userID, reveal)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"likes":    likes,
		"revealed": reveal,
	})
}
