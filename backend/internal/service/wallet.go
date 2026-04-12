package service

import (
	"time"

	"hybrid-app/backend/internal/domain"
	"hybrid-app/backend/internal/repository"
)

type WalletService struct{ baseService }

func NewWalletService(repo repository.Repository, idGen func(string) string) *WalletService {
	return &WalletService{baseService{repo: repo, newID: idGen}}
}

func (s *WalletService) Wallet(userID string) (int, []domain.WalletTransaction, []domain.Gift, []domain.SparkPackage, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return 0, nil, nil, nil, err
	}
	if user == nil {
		return 0, nil, nil, nil, ErrUserNotFound
	}
	transactions, err := s.repo.ListTransactions(userID)
	if err != nil {
		return 0, nil, nil, nil, err
	}
	gifts, err := s.repo.ListGifts()
	if err != nil {
		return 0, nil, nil, nil, err
	}
	packages, err := s.repo.ListSparkPackages()
	if err != nil {
		return 0, nil, nil, nil, err
	}
	return user.SparksBalance, transactions, gifts, packages, nil
}

func (s *WalletService) DailyLoginReward(userID string) (*domain.User, int, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, 0, err
	}
	if user == nil {
		return nil, 0, ErrUserNotFound
	}
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
			return nil, 0, err
		}
	}
	user.UpdatedAt = now
	return user, reward, s.repo.SaveUser(user)
}

func (s *WalletService) BoostProfile(userID string) (*domain.User, error) {
	return s.spend(userID, domain.BoostCost, "boost_spend", "Profile boost for 30 minutes")
}

func (s *WalletService) RefillLikes(userID string) (*domain.User, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if user.SparksBalance < domain.LikeRefillCost {
		return nil, ErrNotEnoughSparks
	}
	user.SparksBalance -= domain.LikeRefillCost
	user.DailyLikeQuota = domain.DefaultDailyLikeCap
	user.UpdatedAt = time.Now().UTC()
	if err := s.repo.SaveUser(user); err != nil {
		return nil, err
	}
	if _, err := addTransaction(s.repo, s.newID, userID, "like_refill", -domain.LikeRefillCost, "Refilled daily likes"); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *WalletService) Activity(userID string) (map[string][]domain.Notification, error) {
	if user, err := s.repo.FindUserByID(userID); err != nil {
		return nil, err
	} else if user == nil {
		return nil, ErrUserNotFound
	}
	notifications, err := s.repo.ListNotifications(userID)
	if err != nil {
		return nil, err
	}
	grouped := map[string][]domain.Notification{"messages": {}, "sparks": {}, "offers": {}, "likes": {}}
	for _, notification := range notifications {
		grouped[notification.Category] = append(grouped[notification.Category], notification)
	}
	return grouped, nil
}

func (s *WalletService) ReferralCenter(userID string) (map[string]any, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	records, err := s.repo.ListReferralRecords(userID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"code":         user.ReferralCode,
		"shareLink":    "https://sparkapp.example/invite/" + user.ReferralCode,
		"shareMessage": "Join with my code " + user.ReferralCode + " and earn Sparks.",
		"history":      records,
		"totalEarned":  len(records) * domain.ReferralReward,
	}, nil
}

func (s *WalletService) spend(userID string, amount int, txnType, description string) (*domain.User, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if user.SparksBalance < amount {
		return nil, ErrNotEnoughSparks
	}
	user.SparksBalance -= amount
	user.UpdatedAt = time.Now().UTC()
	if err := s.repo.SaveUser(user); err != nil {
		return nil, err
	}
	if _, err := addTransaction(s.repo, s.newID, userID, txnType, -amount, description); err != nil {
		return nil, err
	}
	return user, nil
}
