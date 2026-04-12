package http

import (
	"net/http"

	"hybrid-app/backend/internal/domain"
)

func (s *Server) handleMe(w http.ResponseWriter, _ *http.Request, userID string) {
	user, err := s.services.Profile.Me(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (s *Server) handleUpdateMe(w http.ResponseWriter, r *http.Request, userID string) {
	var input struct {
		domain.ProfileInput
		HideFromCities []string `json:"hideFromCities"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := s.services.Profile.UpdateMe(userID, input.ProfileInput, input.HideFromCities)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}
