package service

import (
	"fmt"
	"strings"
	"time"

	"hybrid-app/backend/internal/domain"
	"hybrid-app/backend/internal/repository"
)

type EcosystemService struct{ baseService }

func NewEcosystemService(repo repository.Repository, idGen func(string) string) *EcosystemService {
	return &EcosystemService{baseService{repo: repo, newID: idGen}}
}

func (s *EcosystemService) StartLiveVerification(userID string) (map[string]any, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	user.PhoneVerified = true
	if user.LiveVideoStatus == "" || user.LiveVideoStatus == "not_started" {
		user.LiveVideoStatus = "pending_review"
	}
	user.UpdatedAt = time.Now().UTC()
	if err := s.repo.SaveUser(user); err != nil {
		return nil, err
	}
	return map[string]any{
		"verification": map[string]any{
			"sessionId":      s.newID("live"),
			"provider":       "mock_live_video",
			"status":         user.LiveVideoStatus,
			"expiresIn":      "10m",
			"livenessChecks": []string{"blink", "head_turn", "voice_prompt"},
		},
	}, nil
}

func (s *EcosystemService) CompleteLiveVerification(userID, verdict string) (*domain.User, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	verdict = strings.TrimSpace(strings.ToLower(verdict))
	if verdict == "" {
		verdict = "verified"
	}
	switch verdict {
	case "verified", "review", "rejected":
	default:
		return nil, fmt.Errorf("unsupported live verification verdict")
	}
	user.PhoneVerified = true
	user.LiveVideoStatus = verdict
	user.VerificationStatus = verdict
	if verdict == "verified" && user.ProfileCompletion < 50 {
		user.ProfileCompletion = 50
	}
	user.UpdatedAt = time.Now().UTC()
	if err := s.repo.SaveUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *EcosystemService) VaultStatus(userID string) (map[string]any, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	settings := map[string]any{
		"e2ee":                  true,
		"blockScreenshots":      true,
		"iosRecordingDetection": true,
		"mediaExport":           "restricted",
	}
	for key, value := range user.SafetySettings {
		settings[key] = value
	}
	return map[string]any{
		"vault": map[string]any{
			"enabled":  true,
			"protocol": "vault_v1",
			"settings": settings,
		},
	}, nil
}

func (s *EcosystemService) VouchStatus(userID string) (map[string]any, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return map[string]any{
		"vouch": map[string]any{
			"badge":            user.VouchedBadge,
			"vouchedBy":        user.VouchedBy,
			"premiumActive":    user.PremiumUntil.After(time.Now().UTC()),
			"premiumUntil":     user.PremiumUntil,
			"algorithmicBoost": user.VouchedBadge,
		},
	}, nil
}

func (s *EcosystemService) CreateVouchInvite(userID, inviteeName, inviteePhone string) (map[string]any, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return map[string]any{
		"invite": map[string]any{
			"code":         user.ReferralCode,
			"inviteeName":  strings.TrimSpace(inviteeName),
			"inviteePhone": strings.TrimSpace(inviteePhone),
			"shareMessage": fmt.Sprintf("%s invited you to verify and vouch for him on the app. Use code %s.", user.Name, user.ReferralCode),
		},
	}, nil
}

func (s *EcosystemService) ConfirmVouch(voucherUserID, targetUserID string) (map[string]any, error) {
	voucher, err := s.repo.FindUserByID(voucherUserID)
	if err != nil {
		return nil, err
	}
	if voucher == nil {
		return nil, ErrUserNotFound
	}
	target, err := s.repo.FindUserByID(targetUserID)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, ErrTargetUserNotFound
	}
	if voucher.ID == target.ID {
		return nil, fmt.Errorf("users cannot vouch for themselves")
	}
	if !strings.EqualFold(voucher.Gender, "female") {
		return nil, fmt.Errorf("only verified women can vouch in this flow")
	}
	if voucher.VerificationStatus != "verified" && voucher.LiveVideoStatus != "verified" {
		return nil, fmt.Errorf("voucher must complete verification first")
	}

	now := time.Now().UTC()
	target.VouchedBadge = true
	target.VouchedBy = voucher.ID
	target.SubscriptionTier = "premium"
	if target.PremiumUntil.After(now) {
		target.PremiumUntil = target.PremiumUntil.Add(7 * 24 * time.Hour)
	} else {
		target.PremiumUntil = now.Add(7 * 24 * time.Hour)
	}
	target.PendingBanner = "You were vouched for. Premium unlocked for 1 week."
	target.UpdatedAt = now
	if err := s.repo.SaveUser(target); err != nil {
		return nil, err
	}
	if err := addNotification(s.repo, s.newID, target.ID, "offers", "Vouched badge unlocked", fmt.Sprintf("%s vouched for you. Your profile now has a trust boost.", voucher.Name)); err != nil {
		return nil, err
	}
	return map[string]any{
		"vouch": map[string]any{
			"targetUserId":     target.ID,
			"badge":            target.VouchedBadge,
			"premiumUntil":     target.PremiumUntil,
			"algorithmicBoost": true,
		},
	}, nil
}

func (s *EcosystemService) PortalLink(userID, destination string) (map[string]any, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	token := "portal_" + user.ID + "_" + strings.TrimPrefix(s.newID("session"), "session_")
	if err := s.repo.SaveToken(token, user.ID); err != nil {
		return nil, err
	}
	if destination == "" {
		destination = "ai"
	}
	return map[string]any{
		"portal": map[string]any{
			"token":       token,
			"url":         "https://portal.sparkapp.example/sso?token=" + token + "&destination=" + destination,
			"destination": destination,
			"expiresIn":   "15m",
		},
	}, nil
}

func (s *EcosystemService) ExchangePortalToken(portalToken string) (map[string]any, error) {
	userID, err := s.repo.FindUserIDByToken(strings.TrimSpace(portalToken))
	if err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, ErrUnauthorized
	}
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUnauthorized
	}
	sessionToken := "token_" + user.ID + "_" + strings.TrimPrefix(s.newID("session"), "session_")
	if err := s.repo.SaveToken(sessionToken, user.ID); err != nil {
		return nil, err
	}
	return map[string]any{
		"session": map[string]any{
			"token":       sessionToken,
			"user":        user,
			"surface":     "web",
			"linkedLogin": true,
		},
	}, nil
}

func (s *EcosystemService) Personas(surface string) []map[string]any {
	if strings.TrimSpace(surface) == "" {
		surface = "native"
	}
	return []map[string]any{
		{
			"id":          "wingman_coach",
			"name":        "Wingman Coach",
			"surface":     surface,
			"mode":        "practice",
			"safeForWork": surface != "web",
		},
		{
			"id":          "icebreaker_guide",
			"name":        "Icebreaker Guide",
			"surface":     surface,
			"mode":        "conversation",
			"safeForWork": surface != "web",
		},
	}
}

func (s *EcosystemService) AIChat(userID, surface, personaID, message string, requestImage bool) (map[string]any, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	surface = strings.ToLower(strings.TrimSpace(surface))
	if surface == "" {
		surface = "native"
	}
	if personaID == "" {
		personaID = "wingman_coach"
	}

	adultIntent := isAdultPrompt(message) || requestImage
	if surface == "native" && adultIntent {
		portal, err := s.PortalLink(userID, "ai")
		if err != nil {
			return nil, err
		}
		return map[string]any{
			"personaId": personaID,
			"surface":   surface,
			"access":    "locked",
			"reply":     "This conversation can continue in the secure web portal.",
			"media": map[string]any{
				"blurred": true,
				"locked":  true,
				"cta":     "Link your account on the secure Web Portal to continue privately.",
			},
			"portal": portal["portal"],
		}, nil
	}

	feature := "text"
	if requestImage {
		if adultIntent {
			feature = "explicit_image"
		} else {
			feature = "image"
		}
	}
	totalCost := creditCost(feature, 1)
	if surface == "web" {
		if err := s.consumeCredits(user, feature, 1); err != nil {
			return nil, err
		}
	}

	response := map[string]any{
		"personaId": personaID,
		"surface":   surface,
		"access":    "unlocked",
		"reply":     aiReplyFor(surface, personaID, message, adultIntent),
		"usage": map[string]any{
			"creditsCharged":   ternaryCredits(surface == "web", totalCost),
			"remainingCredits": user.WebCreditsBalance,
		},
	}
	if requestImage {
		response["media"] = map[string]any{
			"generated": true,
			"type":      feature,
			"url":       fmt.Sprintf("https://cdn.sparkapp.example/generated/%s.jpg", personaID),
		}
	}
	return response, nil
}

func (s *EcosystemService) WebWallet(userID string) (map[string]any, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return map[string]any{
		"subscription": map[string]any{
			"tier":           user.SubscriptionTier,
			"priceUSD":       domain.WebSubscriptionCost,
			"monthlyCredits": domain.WebMonthlyCredits,
			"premiumUntil":   user.PremiumUntil,
			"active":         user.PremiumUntil.After(time.Now().UTC()),
		},
		"wallet": map[string]any{
			"credits": user.WebCreditsBalance,
			"pricing": map[string]int{
				"text":          domain.WebTextCreditCost,
				"image":         domain.WebImageCreditCost,
				"explicitImage": domain.WebExplicitImageCost,
			},
		},
	}, nil
}

func (s *EcosystemService) SubscribeWeb(userID, planID string) (map[string]any, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if planID == "" {
		planID = "web_plus_1000"
	}
	now := time.Now().UTC()
	user.WebCreditsBalance += domain.WebMonthlyCredits
	if user.PremiumUntil.After(now) {
		user.PremiumUntil = user.PremiumUntil.Add(30 * 24 * time.Hour)
	} else {
		user.PremiumUntil = now.Add(30 * 24 * time.Hour)
	}
	user.SubscriptionTier = "web_plus"
	user.UpdatedAt = now
	if err := s.repo.SaveUser(user); err != nil {
		return nil, err
	}
	if _, err := addTransaction(s.repo, s.newID, user.ID, "web_subscription_credit", domain.WebMonthlyCredits, "Monthly web subscription credits"); err != nil {
		return nil, err
	}
	return map[string]any{
		"billing": map[string]any{
			"provider":     "stripe",
			"planId":       planID,
			"amountUSD":    domain.WebSubscriptionCost,
			"status":       "active",
			"creditsAdded": domain.WebMonthlyCredits,
			"premiumUntil": user.PremiumUntil,
		},
	}, nil
}

func (s *EcosystemService) ConsumeCredits(userID, feature string, units int) (map[string]any, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if err := s.consumeCredits(user, feature, units); err != nil {
		return nil, err
	}
	return map[string]any{
		"usage": map[string]any{
			"feature":          feature,
			"units":            units,
			"remainingCredits": user.WebCreditsBalance,
		},
	}, nil
}

func (s *EcosystemService) InferencePolicy() map[string]any {
	return map[string]any{
		"rateLimits": map[string]any{
			"textPerMinute":          60,
			"standardImagesPerHour":  40,
			"explicitImagesPerHour":  10,
			"burstQueueConcurrency":  4,
			"gpuLoadBalancingPolicy": "least_loaded_region",
		},
		"routing": map[string]any{
			"nativeSurface": "sfw_only",
			"webSurface":    "full_persona_access",
		},
	}
}

func (s *EcosystemService) consumeCredits(user *domain.User, feature string, units int) error {
	totalCost := creditCost(feature, units)
	if totalCost <= 0 {
		return fmt.Errorf("unsupported credit feature")
	}
	if user.WebCreditsBalance < totalCost {
		return ErrNotEnoughCredits
	}
	user.WebCreditsBalance -= totalCost
	user.UpdatedAt = time.Now().UTC()
	if err := s.repo.SaveUser(user); err != nil {
		return err
	}
	_, err := addTransaction(s.repo, s.newID, user.ID, "web_"+feature, -totalCost, fmt.Sprintf("Consumed %d credits for %s", totalCost, feature))
	return err
}

func creditCost(feature string, units int) int {
	switch strings.ToLower(strings.TrimSpace(feature)) {
	case "text":
		return units * domain.WebTextCreditCost
	case "image":
		return units * domain.WebImageCreditCost
	case "explicit_image":
		return units * domain.WebExplicitImageCost
	default:
		return 0
	}
}

func isAdultPrompt(message string) bool {
	normalized := strings.ToLower(message)
	for _, token := range []string{"adult", "nsfw", "private photo", "nude", "explicit", "unfiltered"} {
		if strings.Contains(normalized, token) {
			return true
		}
	}
	return false
}

func aiReplyFor(surface, personaID, message string, adultIntent bool) string {
	if surface == "web" && adultIntent {
		return "Private mode unlocked. Your request has been queued with the selected persona."
	}
	if strings.Contains(strings.ToLower(message), "first message") {
		return "Try leading with something specific and easy to answer, like a detail from their profile plus a light question."
	}
	if personaID == "icebreaker_guide" {
		return "A strong opener is playful, short, and personal. Ask about something they clearly care about."
	}
	return "Keep it warm, specific, and low-pressure. The best openers sound curious, not rehearsed."
}

func ternaryCredits(condition bool, value int) int {
	if condition {
		return value
	}
	return 0
}
