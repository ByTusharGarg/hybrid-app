package repository

import "hybrid-app/backend/internal/domain"

type Composite struct {
	AuthRepo       AuthRepository
	OnboardingRepo OnboardingRepository
	DiscoveryRepo  DiscoveryRepository
	WalletRepo     WalletRepository
	ChatRepo       ChatRepository
	ProfileRepo    ProfileRepository
}

func (c *Composite) CreateOTP(challenge domain.OTPChallenge) error {
	return c.AuthRepo.CreateOTP(challenge)
}

func (c *Composite) GetOTP(id string) (*domain.OTPChallenge, error) {
	return c.AuthRepo.GetOTP(id)
}

func (c *Composite) SaveOTP(challenge *domain.OTPChallenge) error {
	return c.AuthRepo.SaveOTP(challenge)
}

func (c *Composite) SaveToken(token, userID string) error {
	return c.AuthRepo.SaveToken(token, userID)
}

func (c *Composite) FindUserIDByToken(token string) (string, error) {
	return c.AuthRepo.FindUserIDByToken(token)
}

func (c *Composite) FindUserByPhone(phone string) (*domain.User, error) {
	return c.AuthRepo.FindUserByPhone(phone)
}

func (c *Composite) FindUserByReferralCode(code string) (*domain.User, error) {
	return c.AuthRepo.FindUserByReferralCode(code)
}

func (c *Composite) SaveUser(user *domain.User) error {
	return c.AuthRepo.SaveUser(user)
}

func (c *Composite) FindUserByID(id string) (*domain.User, error) {
	return c.AuthRepo.FindUserByID(id)
}

func (c *Composite) ListOtherCompletedUsers(excludeUserID string) ([]*domain.User, error) {
	return c.OnboardingRepo.ListOtherCompletedUsers(excludeUserID)
}

func (c *Composite) AddTransaction(txn domain.WalletTransaction) error {
	return c.OnboardingRepo.AddTransaction(txn)
}

func (c *Composite) ListTransactions(userID string) ([]domain.WalletTransaction, error) {
	return c.OnboardingRepo.ListTransactions(userID)
}

func (c *Composite) ListGifts() ([]domain.Gift, error) {
	return c.DiscoveryRepo.ListGifts()
}

func (c *Composite) ListSparkPackages() ([]domain.SparkPackage, error) {
	return c.WalletRepo.ListSparkPackages()
}

func (c *Composite) AddLike(like domain.LikeAction) error {
	return c.DiscoveryRepo.AddLike(like)
}

func (c *Composite) ListLikesForTarget(userID string) ([]domain.LikeAction, error) {
	return c.DiscoveryRepo.ListLikesForTarget(userID)
}

func (c *Composite) HasPositiveLike(fromUserID, toUserID string) (bool, error) {
	return c.DiscoveryRepo.HasPositiveLike(fromUserID, toUserID)
}

func (c *Composite) SaveMatch(match domain.Match) error {
	return c.DiscoveryRepo.SaveMatch(match)
}

func (c *Composite) FindMatchBetween(userA, userB string) (*domain.Match, error) {
	return c.DiscoveryRepo.FindMatchBetween(userA, userB)
}

func (c *Composite) ListMatchesForUser(userID string) ([]domain.Match, error) {
	return c.DiscoveryRepo.ListMatchesForUser(userID)
}

func (c *Composite) SaveChat(chat *domain.ChatSummary) error {
	return c.ChatRepo.SaveChat(chat)
}

func (c *Composite) FindChatByID(chatID string) (*domain.ChatSummary, error) {
	return c.ChatRepo.FindChatByID(chatID)
}

func (c *Composite) FindChatByMatchID(matchID string) (*domain.ChatSummary, error) {
	return c.DiscoveryRepo.FindChatByMatchID(matchID)
}

func (c *Composite) ListChatsForUser(userID string) ([]*domain.ChatSummary, error) {
	return c.ChatRepo.ListChatsForUser(userID)
}

func (c *Composite) AddMessage(message domain.ChatMessage) error {
	return c.ChatRepo.AddMessage(message)
}

func (c *Composite) ListMessages(chatID string) ([]domain.ChatMessage, error) {
	return c.ChatRepo.ListMessages(chatID)
}

func (c *Composite) AddNotification(notification domain.Notification) error {
	return c.WalletRepo.AddNotification(notification)
}

func (c *Composite) ListNotifications(userID string) ([]domain.Notification, error) {
	return c.WalletRepo.ListNotifications(userID)
}

func (c *Composite) AddReferralRecord(record domain.ReferralRecord) error {
	return c.OnboardingRepo.AddReferralRecord(record)
}

func (c *Composite) ListReferralRecords(referrerUserID string) ([]domain.ReferralRecord, error) {
	return c.WalletRepo.ListReferralRecords(referrerUserID)
}
