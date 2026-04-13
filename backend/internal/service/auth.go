package service

import (
	"strings"
	"time"

	"hybrid-app/backend/internal/domain"
	"hybrid-app/backend/internal/repository"
)

type AuthService struct{ baseService }

func NewAuthService(repo repository.Repository, idGen func(string) string) *AuthService {
	return &AuthService{baseService{repo: repo, newID: idGen}}
}

func (s *AuthService) StartRegistration(phone, referral string) (domain.OTPChallenge, error) {
	challenge := domain.OTPChallenge{
		ID:        s.newID("otp"),
		Phone:     phone,
		Code:      "123456",
		Purpose:   "register",
		Referral:  strings.TrimSpace(referral),
		CreatedAt: time.Now().UTC(),
	}
	return challenge, s.repo.CreateOTP(challenge)
}

func (s *AuthService) FinishRegistration(otpSessionID, code string) (*domain.User, string, bool, error) {
	challenge, err := s.repo.GetOTP(otpSessionID)
	if err != nil {
		return nil, "", false, err
	}
	if challenge == nil {
		return nil, "", false, ErrChallengeNotFound
	}
	if challenge.Code != code {
		return nil, "", false, ErrInvalidOTP
	}
	challenge.VerifiedAt = time.Now().UTC()
	if err := s.repo.SaveOTP(challenge); err != nil {
		return nil, "", false, err
	}

	existing, err := s.repo.FindUserByPhone(challenge.Phone)
	if err != nil {
		return nil, "", false, err
	}
	if existing != nil {
		token := s.issueToken(existing.ID)
		return existing, token, false, nil
	}

	referredBy := ""
	if challenge.Referral != "" {
		referrer, err := s.repo.FindUserByReferralCode(challenge.Referral)
		if err != nil {
			return nil, "", false, err
		}
		if referrer != nil {
			referredBy = referrer.ID
		}
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:                   s.newID("user"),
		Phone:                challenge.Phone,
		Name:                 "New User",
		Gender:               "unknown",
		InterestedIn:         []string{},
		City:                 "",
		DateOfBirth:          "",
		Photos:               []string{},
		Interests:            []string{},
		Questionnaire:        map[string]any{},
		CompatibilityTags:    []string{},
		ReferralCode:         strings.ToUpper("SPK" + strings.TrimPrefix(s.newID("code"), "code_")),
		ReferredBy:           referredBy,
		PhoneVerified:        true,
		VerificationStatus:   "pending",
		LiveVideoStatus:      "not_started",
		OnboardingCompleted:  false,
		ProfileCompletion:    10,
		SparksBalance:        0,
		WebCreditsBalance:    0,
		DailyLikeQuota:       domain.DefaultDailyLikeCap,
		LastLikeQuotaResetAt: now,
		SubscriptionTier:     "free",
		SafetySettings: map[string]any{
			"e2ee":                  true,
			"blockScreenshots":      true,
			"iosRecordingDetection": true,
			"mediaExport":           "restricted",
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.SaveUser(user); err != nil {
		return nil, "", false, err
	}
	token := s.issueToken(user.ID)
	return user, token, true, nil
}

func (s *AuthService) StartLogin(phone string) (domain.OTPChallenge, error) {
	challenge := domain.OTPChallenge{
		ID:        s.newID("otp"),
		Phone:     phone,
		Code:      "123456",
		Purpose:   "login",
		CreatedAt: time.Now().UTC(),
	}
	return challenge, s.repo.CreateOTP(challenge)
}

func (s *AuthService) FinishLogin(otpSessionID, code string) (*domain.User, string, int, error) {
	challenge, err := s.repo.GetOTP(otpSessionID)
	if err != nil {
		return nil, "", 0, err
	}
	if challenge == nil {
		return nil, "", 0, ErrChallengeNotFound
	}
	if challenge.Code != code {
		return nil, "", 0, ErrInvalidOTP
	}
	challenge.VerifiedAt = time.Now().UTC()
	if err := s.repo.SaveOTP(challenge); err != nil {
		return nil, "", 0, err
	}

	user, err := s.repo.FindUserByPhone(challenge.Phone)
	if err != nil {
		return nil, "", 0, err
	}
	if user == nil {
		return nil, "", 0, ErrUserNotRegistered
	}
	user.PhoneVerified = true

	reward := 0
	now := time.Now().UTC()
	if !sameDay(user.LastLoginAt, now) {
		if sameDay(user.LastLoginAt.Add(24*time.Hour), now) {
			user.LoginStreak++
		} else {
			user.LoginStreak = 1
		}
		user.LastLoginAt = now
		user.SparksBalance += domain.DailyLoginReward
		user.PendingBanner = "+30 Sparks from your login streak!"
		reward = domain.DailyLoginReward
		if _, err := addTransaction(s.repo, s.newID, user.ID, "credit", reward, "Daily login streak reward"); err != nil {
			return nil, "", 0, err
		}
		if err := addNotification(s.repo, s.newID, user.ID, "sparks", "Daily streak reward", user.PendingBanner); err != nil {
			return nil, "", 0, err
		}
	} else {
		user.PendingBanner = ""
	}
	user.UpdatedAt = now
	if err := s.repo.SaveUser(user); err != nil {
		return nil, "", 0, err
	}
	token := s.issueToken(user.ID)
	return user, token, reward, nil
}

func (s *AuthService) Authenticate(token string) (*domain.User, error) {
	userID, err := s.repo.FindUserIDByToken(token)
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
	now := time.Now().UTC()
	if !sameDay(user.LastLikeQuotaResetAt, now) {
		user.DailyLikeQuota = domain.DefaultDailyLikeCap
		user.LastLikeQuotaResetAt = now
		if err := s.repo.SaveUser(user); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (s *AuthService) issueToken(userID string) string {
	token := "token_" + userID + "_" + strings.TrimPrefix(s.newID("session"), "session_")
	_ = s.repo.SaveToken(token, userID)
	return token
}
