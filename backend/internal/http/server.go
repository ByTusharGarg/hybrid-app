package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"hybrid-app/backend/internal/service"
)

type Server struct {
	services *service.Services
	mux      *http.ServeMux
}

func NewServer(services *service.Services) *Server {
	s := &Server{
		services: services,
		mux:      http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) Router() http.Handler {
	return withCORS(s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("GET /swagger", s.handleSwaggerUI)
	s.mux.HandleFunc("GET /openapi.json", s.handleOpenAPI)
	s.mux.HandleFunc("POST /api/v1/auth/register/start", s.handleRegisterStart)
	s.mux.HandleFunc("POST /api/v1/auth/register/verify-otp", s.handleRegisterVerify)
	s.mux.HandleFunc("POST /api/v1/auth/login/start", s.handleLoginStart)
	s.mux.HandleFunc("POST /api/v1/auth/login/verify-otp", s.handleLoginVerify)
	s.mux.HandleFunc("POST /api/v1/auth/portal/exchange", s.handlePortalExchange)

	s.mux.HandleFunc("GET /api/v1/home", s.authed(s.handleHome))
	s.mux.HandleFunc("POST /api/v1/onboarding/questionnaire", s.authed(s.handleQuestionnaire))
	s.mux.HandleFunc("POST /api/v1/onboarding/gender-verification", s.authed(s.handleGenderVerification))
	s.mux.HandleFunc("POST /api/v1/onboarding/live-verification/start", s.authed(s.handleLiveVerificationStart))
	s.mux.HandleFunc("POST /api/v1/onboarding/live-verification/complete", s.authed(s.handleLiveVerificationComplete))
	s.mux.HandleFunc("POST /api/v1/onboarding/profile", s.authed(s.handleProfileSetup))

	s.mux.HandleFunc("GET /api/v1/discover", s.authed(s.handleDiscover))
	s.mux.HandleFunc("POST /api/v1/discover/actions", s.authed(s.handleDiscoverAction))
	s.mux.HandleFunc("GET /api/v1/matches", s.authed(s.handleMatches))
	s.mux.HandleFunc("GET /api/v1/likes-you", s.authed(s.handleLikesYou))

	s.mux.HandleFunc("GET /api/v1/chats", s.authed(s.handleChats))
	s.mux.HandleFunc("GET /api/v1/chats/", s.authed(s.handleChatMessages))
	s.mux.HandleFunc("POST /api/v1/chats/", s.authed(s.handleSendMessage))

	s.mux.HandleFunc("GET /api/v1/activity", s.authed(s.handleActivity))
	s.mux.HandleFunc("GET /api/v1/wallet", s.authed(s.handleWallet))
	s.mux.HandleFunc("POST /api/v1/wallet/daily-login", s.authed(s.handleDailyLoginReward))
	s.mux.HandleFunc("POST /api/v1/wallet/boost", s.authed(s.handleBoost))
	s.mux.HandleFunc("POST /api/v1/wallet/like-refill", s.authed(s.handleLikeRefill))
	s.mux.HandleFunc("GET /api/v1/referrals", s.authed(s.handleReferrals))
	s.mux.HandleFunc("GET /api/v1/me", s.authed(s.handleMe))
	s.mux.HandleFunc("PATCH /api/v1/me", s.authed(s.handleUpdateMe))
	s.mux.HandleFunc("GET /api/v1/security/vault", s.authed(s.handleVaultStatus))
	s.mux.HandleFunc("GET /api/v1/vouch/status", s.authed(s.handleVouchStatus))
	s.mux.HandleFunc("POST /api/v1/vouch/invite", s.authed(s.handleVouchInvite))
	s.mux.HandleFunc("POST /api/v1/vouch/confirm", s.authed(s.handleVouchConfirm))
	s.mux.HandleFunc("GET /api/v1/ai/personas", s.authed(s.handleAIPersonas))
	s.mux.HandleFunc("POST /api/v1/ai/teaser", s.authed(s.handleAITeaser))
	s.mux.HandleFunc("POST /api/v1/auth/portal/link", s.authed(s.handlePortalLink))
	s.mux.HandleFunc("GET /api/v1/web/wallet", s.authed(s.handleWebWallet))
	s.mux.HandleFunc("POST /api/v1/web/billing/subscribe", s.authed(s.handleWebSubscribe))
	s.mux.HandleFunc("POST /api/v1/web/billing/consume", s.authed(s.handleWebConsume))
	s.mux.HandleFunc("GET /api/v1/web/inference/policy", s.authed(s.handleInferencePolicy))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) authed(next func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := bearerToken(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err)
			return
		}
		user, err := s.services.Auth.Authenticate(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err)
			return
		}
		next(w, r, user.ID)
	}
}

func bearerToken(r *http.Request) (string, error) {
	value := strings.TrimSpace(r.Header.Get("Authorization"))
	if value == "" {
		return "", errors.New("missing authorization header")
	}
	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header")
	}
	return parts[1], nil
}

func decodeJSON(r *http.Request, dest any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dest)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]any{
		"error": err.Error(),
	})
}
