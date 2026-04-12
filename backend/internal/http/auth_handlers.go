package http

import (
	"errors"
	"net/http"
	"strings"

	"hybrid-app/backend/internal/domain"
)

func (s *Server) handleRegisterStart(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Phone        string `json:"phone"`
		ReferralCode string `json:"referralCode"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if strings.TrimSpace(input.Phone) == "" {
		writeError(w, http.StatusBadRequest, errors.New("phone is required"))
		return
	}

	challenge, err := s.services.Auth.StartRegistration(input.Phone, input.ReferralCode)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"challenge": map[string]any{
			"id":         challenge.ID,
			"phone":      challenge.Phone,
			"purpose":    challenge.Purpose,
			"expiresIn":  "5m",
			"mockOTP":    challenge.Code,
			"referralIn": challenge.Referral != "",
		},
		"registration": map[string]any{
			"referralBonusSparks": domain.ReferralJoinBonus,
			"welcomeBonusSparks":  domain.WelcomeBonusSparks,
		},
	})
}

func (s *Server) handleRegisterVerify(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OTPSessionID string `json:"otpSessionId"`
		Code         string `json:"code"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, token, isNew, err := s.services.Auth.FinishRegistration(input.OTPSessionID, input.Code)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"session": map[string]any{
			"token": token,
			"user":  user,
		},
		"registration": map[string]any{
			"isNewUser":          isNew,
			"onboardingState":    onboardingState(user),
			"referralApplied":    user.ReferredBy != "",
			"welcomeBonusSparks": domain.WelcomeBonusSparks,
		},
	})
}

func (s *Server) handleLoginStart(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Phone string `json:"phone"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if strings.TrimSpace(input.Phone) == "" {
		writeError(w, http.StatusBadRequest, errors.New("phone is required"))
		return
	}

	challenge, err := s.services.Auth.StartLogin(input.Phone)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"challenge": map[string]any{
			"id":        challenge.ID,
			"phone":     challenge.Phone,
			"purpose":   challenge.Purpose,
			"expiresIn": "5m",
			"mockOTP":   challenge.Code,
		},
	})
}

func (s *Server) handleLoginVerify(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OTPSessionID string `json:"otpSessionId"`
		Code         string `json:"code"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, token, reward, err := s.services.Auth.FinishLogin(input.OTPSessionID, input.Code)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"session": map[string]any{
			"token": token,
			"user":  user,
		},
		"reward": map[string]any{
			"sparks": reward,
			"reason": "daily_login_streak",
			"banner": user.PendingBanner,
		},
		"onboardingState": onboardingState(user),
	})
}
