package service

import (
	"fmt"
	"slices"
	"time"

	"hybrid-app/backend/internal/domain"
	"hybrid-app/backend/internal/repository"
)

type OnboardingService struct{ baseService }

func NewOnboardingService(repo repository.Repository, idGen func(string) string) *OnboardingService {
	return &OnboardingService{baseService{repo: repo, newID: idGen}}
}

func (s *OnboardingService) Home(userID string) (map[string]any, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	others, err := s.repo.ListOtherCompletedUsers(userID)
	if err != nil {
		return nil, err
	}
	topPicks, feed := buildDiscovery(user, others)
	return map[string]any{
		"user":               user,
		"pendingBanner":      user.PendingBanner,
		"dailyLikeQuota":     user.DailyLikeQuota,
		"topPicks":           topPicks,
		"discoverFeed":       feed,
		"boostCostSparks":    domain.BoostCost,
		"superLikeCost":      domain.SuperLikeCost,
		"likesRevealCost":    domain.LikesRevealCost,
		"referralJoinBonus":  domain.ReferralJoinBonus,
		"welcomeBonusSparks": domain.WelcomeBonusSparks,
	}, nil
}

func (s *OnboardingService) SaveQuestionnaire(userID string, answers map[string]any) (*domain.User, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	user.Questionnaire = answers
	user.CompatibilityTags = tagsFromQuestionnaire(answers)
	if user.ProfileCompletion < 35 {
		user.ProfileCompletion = 35
	}
	user.UpdatedAt = time.Now().UTC()
	return user, s.repo.SaveUser(user)
}

func (s *OnboardingService) SaveGenderVerification(userID, status string) (*domain.User, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	user.VerificationStatus = status
	if user.ProfileCompletion < 50 {
		user.ProfileCompletion = 50
	}
	user.UpdatedAt = time.Now().UTC()
	return user, s.repo.SaveUser(user)
}

func (s *OnboardingService) CompleteProfile(userID string, input domain.ProfileInput) (*domain.User, []domain.WalletTransaction, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, ErrUserNotFound
	}

	user.Name = input.Name
	user.Bio = input.Bio
	user.Gender = input.Gender
	user.InterestedIn = input.InterestedIn
	user.City = input.City
	user.DateOfBirth = input.DateOfBirth
	user.Photos = input.Photos
	user.Interests = input.Interests
	user.InstagramHandle = input.InstagramHandle
	user.ProfileCompletion = 100
	user.OnboardingCompleted = true
	user.UpdatedAt = time.Now().UTC()

	transactions, err := s.repo.ListTransactions(userID)
	if err != nil {
		return nil, nil, err
	}
	var rewards []domain.WalletTransaction
	if !hasTransactionType(transactions, "welcome_bonus") {
		user.SparksBalance += domain.WelcomeBonusSparks
		txn, err := addTransaction(s.repo, s.newID, userID, "welcome_bonus", domain.WelcomeBonusSparks, "Welcome bonus Sparks")
		if err != nil {
			return nil, nil, err
		}
		rewards = append(rewards, txn)
	}
	if user.ReferredBy != "" && !hasTransactionType(transactions, "referral_join_bonus") {
		user.SparksBalance += domain.ReferralJoinBonus
		txn, err := addTransaction(s.repo, s.newID, userID, "referral_join_bonus", domain.ReferralJoinBonus, "Referral join bonus")
		if err != nil {
			return nil, nil, err
		}
		rewards = append(rewards, txn)

		referrer, err := s.repo.FindUserByID(user.ReferredBy)
		if err != nil {
			return nil, nil, err
		}
		if referrer != nil {
			referrer.SparksBalance += domain.ReferralReward
			referrer.UpdatedAt = time.Now().UTC()
			if err := s.repo.SaveUser(referrer); err != nil {
				return nil, nil, err
			}
			if _, err := addTransaction(s.repo, s.newID, referrer.ID, "referral_reward", domain.ReferralReward, fmt.Sprintf("Referral reward for %s", user.Name)); err != nil {
				return nil, nil, err
			}
			if err := s.repo.AddReferralRecord(domain.ReferralRecord{
				ID:             s.newID("ref"),
				ReferrerUserID: referrer.ID,
				ReferredUserID: user.ID,
				RewardSparks:   domain.ReferralReward,
				Status:         "completed",
				CompletedAt:    time.Now().UTC(),
			}); err != nil {
				return nil, nil, err
			}
			if err := addNotification(s.repo, s.newID, referrer.ID, "sparks", "Referral reward unlocked", fmt.Sprintf("%s completed registration. You earned %d Sparks.", user.Name, domain.ReferralReward)); err != nil {
				return nil, nil, err
			}
		}
	}
	if !hasTransactionType(transactions, "profile_completion_bonus") {
		user.SparksBalance += 20
		txn, err := addTransaction(s.repo, s.newID, userID, "profile_completion_bonus", 20, "Completed profile bonus")
		if err != nil {
			return nil, nil, err
		}
		rewards = append(rewards, txn)
	}

	if err := s.repo.SaveUser(user); err != nil {
		return nil, nil, err
	}
	message := fmt.Sprintf("Congratulations! You earned %d Sparks welcome bonus.", domain.WelcomeBonusSparks)
	if user.ReferredBy != "" {
		message = fmt.Sprintf("Congratulations! You earned %d Sparks welcome bonus + %d referral bonus.", domain.WelcomeBonusSparks, domain.ReferralJoinBonus)
	}
	if err := addNotification(s.repo, s.newID, userID, "offers", "Welcome to Sparks", message); err != nil {
		return nil, nil, err
	}
	return user, rewards, nil
}

func buildDiscovery(user *domain.User, others []*domain.User) ([]domain.DiscoveryProfile, []domain.DiscoveryProfile) {
	var topPicks []domain.DiscoveryProfile
	var feed []domain.DiscoveryProfile
	for _, other := range others {
		score := compatibility(user, other)
		profile := domain.DiscoveryProfile{
			ID:                 other.ID,
			Name:               other.Name,
			Age:                ageFromDOB(other.DateOfBirth),
			City:               other.City,
			Bio:                other.Bio,
			Photos:             other.Photos,
			Interests:          other.Interests,
			QuestionnaireTags:  other.CompatibilityTags,
			CompatibilityScore: score,
			VerificationBadge:  other.VerificationStatus == "verified",
			LastActive:         humanizeLastActive(other.LastLoginAt),
		}
		feed = append(feed, profile)
		if score >= 75 {
			topPicks = append(topPicks, profile)
		}
	}
	slices.SortFunc(feed, func(a, b domain.DiscoveryProfile) int { return b.CompatibilityScore - a.CompatibilityScore })
	slices.SortFunc(topPicks, func(a, b domain.DiscoveryProfile) int { return b.CompatibilityScore - a.CompatibilityScore })
	return topPicks, feed
}
