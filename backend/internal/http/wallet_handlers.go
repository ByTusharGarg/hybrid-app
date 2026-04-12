package http

import (
	"net/http"

	"hybrid-app/backend/internal/domain"
)

func (s *Server) handleActivity(w http.ResponseWriter, _ *http.Request, userID string) {
	activity, err := s.services.Wallet.Activity(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, activity)
}

func (s *Server) handleWallet(w http.ResponseWriter, _ *http.Request, userID string) {
	balance, transactions, gifts, packages, err := s.services.Wallet.Wallet(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"wallet": map[string]any{
			"balance":      balance,
			"transactions": transactions,
		},
		"catalog": map[string]any{
			"gifts":    gifts,
			"packages": packages,
		},
	})
}

func (s *Server) handleDailyLoginReward(w http.ResponseWriter, _ *http.Request, userID string) {
	user, reward, err := s.services.Wallet.DailyLoginReward(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user": user,
		"reward": map[string]any{
			"sparks": reward,
			"reason": "daily_login_streak",
		},
	})
}

func (s *Server) handleBoost(w http.ResponseWriter, _ *http.Request, userID string) {
	user, err := s.services.Wallet.BoostProfile(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user": user,
		"spend": map[string]any{
			"action":       "boost",
			"costSparks":   domain.BoostCost,
			"durationHint": "30m",
		},
	})
}

func (s *Server) handleLikeRefill(w http.ResponseWriter, _ *http.Request, userID string) {
	user, err := s.services.Wallet.RefillLikes(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user": user,
		"spend": map[string]any{
			"action":     "like_refill",
			"costSparks": domain.LikeRefillCost,
		},
	})
}

func (s *Server) handleReferrals(w http.ResponseWriter, _ *http.Request, userID string) {
	center, err := s.services.Wallet.ReferralCenter(userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, center)
}
