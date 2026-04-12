package http

import (
	"net/http"

	"hybrid-app/backend/internal/domain"
)

func (s *Server) handleHome(w http.ResponseWriter, _ *http.Request, userID string) {
	home, err := s.services.Onboarding.Home(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, home)
}

func (s *Server) handleQuestionnaire(w http.ResponseWriter, r *http.Request, userID string) {
	var input struct {
		Answers map[string]any `json:"answers"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := s.services.Onboarding.SaveQuestionnaire(userID, input.Answers)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user":              user,
		"onboardingState":   onboardingState(user),
		"compatibilityTags": user.CompatibilityTags,
	})
}

func (s *Server) handleGenderVerification(w http.ResponseWriter, r *http.Request, userID string) {
	var input struct {
		Status string `json:"status"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if input.Status == "" {
		input.Status = "verified"
	}

	user, err := s.services.Onboarding.SaveGenderVerification(userID, input.Status)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user":            user,
		"onboardingState": onboardingState(user),
	})
}

func (s *Server) handleProfileSetup(w http.ResponseWriter, r *http.Request, userID string) {
	var input domain.ProfileInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, rewards, err := s.services.Onboarding.CompleteProfile(userID, input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user": user,
		"onboarding": map[string]any{
			"state":      onboardingState(user),
			"completed":  user.OnboardingCompleted,
			"completion": user.ProfileCompletion,
		},
		"walletRewards": rewards,
		"welcome": map[string]any{
			"message":             welcomeCopy(user.ReferredBy != ""),
			"welcomeBonusSparks":  domain.WelcomeBonusSparks,
			"referralBonusSparks": ternaryInt(user.ReferredBy != "", domain.ReferralJoinBonus, 0),
		},
	})
}
