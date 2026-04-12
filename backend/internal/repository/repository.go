package repository

import "hybrid-app/backend/internal/domain"

type AuthRepository interface {
	CreateOTP(challenge domain.OTPChallenge) error
	GetOTP(id string) (*domain.OTPChallenge, error)
	SaveOTP(challenge *domain.OTPChallenge) error

	SaveToken(token, userID string) error
	FindUserIDByToken(token string) (string, error)

	FindUserByPhone(phone string) (*domain.User, error)
	FindUserByReferralCode(code string) (*domain.User, error)
	SaveUser(user *domain.User) error
	FindUserByID(id string) (*domain.User, error)
}

type OnboardingRepository interface {
	FindUserByID(id string) (*domain.User, error)
	SaveUser(user *domain.User) error
	ListOtherCompletedUsers(excludeUserID string) ([]*domain.User, error)

	AddTransaction(txn domain.WalletTransaction) error
	ListTransactions(userID string) ([]domain.WalletTransaction, error)

	AddNotification(notification domain.Notification) error
	AddReferralRecord(record domain.ReferralRecord) error
}

type DiscoveryRepository interface {
	FindUserByID(id string) (*domain.User, error)
	SaveUser(user *domain.User) error
	ListOtherCompletedUsers(excludeUserID string) ([]*domain.User, error)

	AddTransaction(txn domain.WalletTransaction) error
	ListGifts() ([]domain.Gift, error)

	AddLike(like domain.LikeAction) error
	ListLikesForTarget(userID string) ([]domain.LikeAction, error)
	HasPositiveLike(fromUserID, toUserID string) (bool, error)

	SaveMatch(match domain.Match) error
	FindMatchBetween(userA, userB string) (*domain.Match, error)
	ListMatchesForUser(userID string) ([]domain.Match, error)

	SaveChat(chat *domain.ChatSummary) error
	FindChatByMatchID(matchID string) (*domain.ChatSummary, error)

	AddNotification(notification domain.Notification) error
}

type WalletRepository interface {
	FindUserByID(id string) (*domain.User, error)
	SaveUser(user *domain.User) error

	AddTransaction(txn domain.WalletTransaction) error
	ListTransactions(userID string) ([]domain.WalletTransaction, error)
	ListGifts() ([]domain.Gift, error)
	ListSparkPackages() ([]domain.SparkPackage, error)

	AddNotification(notification domain.Notification) error
	ListNotifications(userID string) ([]domain.Notification, error)
	ListReferralRecords(referrerUserID string) ([]domain.ReferralRecord, error)
}

type ChatRepository interface {
	FindUserByID(id string) (*domain.User, error)
	SaveUser(user *domain.User) error
	ListGifts() ([]domain.Gift, error)
	AddTransaction(txn domain.WalletTransaction) error

	SaveChat(chat *domain.ChatSummary) error
	FindChatByID(chatID string) (*domain.ChatSummary, error)
	ListChatsForUser(userID string) ([]*domain.ChatSummary, error)

	AddMessage(message domain.ChatMessage) error
	ListMessages(chatID string) ([]domain.ChatMessage, error)
}

type ProfileRepository interface {
	FindUserByID(id string) (*domain.User, error)
	SaveUser(user *domain.User) error
}

type Repository interface {
	AuthRepository
	OnboardingRepository
	DiscoveryRepository
	WalletRepository
	ChatRepository
	ProfileRepository
}
