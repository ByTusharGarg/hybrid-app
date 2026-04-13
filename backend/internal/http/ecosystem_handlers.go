package http

import "net/http"

func (s *Server) handleLiveVerificationStart(w http.ResponseWriter, _ *http.Request, userID string) {
	payload, err := s.services.Ecosystem.StartLiveVerification(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleLiveVerificationComplete(w http.ResponseWriter, r *http.Request, userID string) {
	var input struct {
		Verdict string `json:"verdict"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	user, err := s.services.Ecosystem.CompleteLiveVerification(userID, input.Verdict)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user":            user,
		"onboardingState": onboardingState(user),
	})
}

func (s *Server) handleVaultStatus(w http.ResponseWriter, _ *http.Request, userID string) {
	payload, err := s.services.Ecosystem.VaultStatus(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleVouchStatus(w http.ResponseWriter, _ *http.Request, userID string) {
	payload, err := s.services.Ecosystem.VouchStatus(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleVouchInvite(w http.ResponseWriter, r *http.Request, userID string) {
	var input struct {
		InviteeName  string `json:"inviteeName"`
		InviteePhone string `json:"inviteePhone"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	payload, err := s.services.Ecosystem.CreateVouchInvite(userID, input.InviteeName, input.InviteePhone)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleVouchConfirm(w http.ResponseWriter, r *http.Request, userID string) {
	var input struct {
		TargetUserID string `json:"targetUserId"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	payload, err := s.services.Ecosystem.ConfirmVouch(userID, input.TargetUserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleAIPersonas(w http.ResponseWriter, r *http.Request, _ string) {
	writeJSON(w, http.StatusOK, map[string]any{
		"personas": s.services.Ecosystem.Personas(r.URL.Query().Get("surface")),
	})
}

func (s *Server) handleAITeaser(w http.ResponseWriter, r *http.Request, userID string) {
	var input struct {
		Surface      string `json:"surface"`
		PersonaID    string `json:"personaId"`
		Message      string `json:"message"`
		RequestImage bool   `json:"requestImage"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	payload, err := s.services.Ecosystem.AIChat(userID, input.Surface, input.PersonaID, input.Message, input.RequestImage)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handlePortalLink(w http.ResponseWriter, r *http.Request, userID string) {
	var input struct {
		Destination string `json:"destination"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	payload, err := s.services.Ecosystem.PortalLink(userID, input.Destination)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handlePortalExchange(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PortalToken string `json:"portalToken"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	payload, err := s.services.Ecosystem.ExchangePortalToken(input.PortalToken)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleWebWallet(w http.ResponseWriter, _ *http.Request, userID string) {
	payload, err := s.services.Ecosystem.WebWallet(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleWebSubscribe(w http.ResponseWriter, r *http.Request, userID string) {
	var input struct {
		PlanID string `json:"planId"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	payload, err := s.services.Ecosystem.SubscribeWeb(userID, input.PlanID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleWebConsume(w http.ResponseWriter, r *http.Request, userID string) {
	var input struct {
		Feature string `json:"feature"`
		Units   int    `json:"units"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if input.Units <= 0 {
		input.Units = 1
	}
	payload, err := s.services.Ecosystem.ConsumeCredits(userID, input.Feature, input.Units)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, payload)
}

func (s *Server) handleInferencePolicy(w http.ResponseWriter, _ *http.Request, _ string) {
	writeJSON(w, http.StatusOK, s.services.Ecosystem.InferencePolicy())
}
