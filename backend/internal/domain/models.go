package domain

import "time"

const (
	WelcomeBonusSparks  = 100
	ReferralJoinBonus   = 50
	ReferralReward      = 100
	DailyLoginReward    = 30
	BoostCost           = 20
	SuperLikeCost       = 10
	UndoCost            = 5
	LikesRevealCost     = 25
	LikeRefillCost      = 15
	DefaultDailyLikeCap = 10
)

type User struct {
	ID                   string         `json:"id"`
	Phone                string         `json:"phone"`
	Name                 string         `json:"name"`
	Bio                  string         `json:"bio"`
	Gender               string         `json:"gender"`
	InterestedIn         []string       `json:"interestedIn"`
	City                 string         `json:"city"`
	DateOfBirth          string         `json:"dateOfBirth"`
	Photos               []string       `json:"photos"`
	Interests            []string       `json:"interests"`
	Questionnaire        map[string]any `json:"questionnaire"`
	CompatibilityTags    []string       `json:"compatibilityTags"`
	ReferralCode         string         `json:"referralCode"`
	ReferredBy           string         `json:"referredBy,omitempty"`
	VerificationStatus   string         `json:"verificationStatus"`
	OnboardingCompleted  bool           `json:"onboardingCompleted"`
	ProfileCompletion    int            `json:"profileCompletion"`
	SparksBalance        int            `json:"sparksBalance"`
	DailyLikeQuota       int            `json:"dailyLikeQuota"`
	LastLikeQuotaResetAt time.Time      `json:"lastLikeQuotaResetAt"`
	LastLoginAt          time.Time      `json:"lastLoginAt"`
	LoginStreak          int            `json:"loginStreak"`
	PendingBanner        string         `json:"pendingBanner,omitempty"`
	SubscriptionTier     string         `json:"subscriptionTier"`
	HideFromCities       []string       `json:"hideFromCities"`
	InstagramHandle      string         `json:"instagramHandle,omitempty"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
}

type ProfileInput struct {
	Name            string   `json:"name"`
	Bio             string   `json:"bio"`
	Gender          string   `json:"gender"`
	InterestedIn    []string `json:"interestedIn"`
	City            string   `json:"city"`
	DateOfBirth     string   `json:"dateOfBirth"`
	Photos          []string `json:"photos"`
	Interests       []string `json:"interests"`
	InstagramHandle string   `json:"instagramHandle"`
}

type OTPChallenge struct {
	ID         string    `json:"id"`
	Phone      string    `json:"phone"`
	Code       string    `json:"code"`
	Purpose    string    `json:"purpose"`
	Referral   string    `json:"referralCode,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	VerifiedAt time.Time `json:"verifiedAt,omitempty"`
}

type WalletTransaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Type        string    `json:"type"`
	Amount      int       `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type DiscoveryProfile struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Age                int      `json:"age"`
	City               string   `json:"city"`
	Bio                string   `json:"bio"`
	Photos             []string `json:"photos"`
	Interests          []string `json:"interests"`
	QuestionnaireTags  []string `json:"questionnaireTags"`
	CompatibilityScore int      `json:"compatibilityScore"`
	VerificationBadge  bool     `json:"verificationBadge"`
	LastActive         string   `json:"lastActive"`
}

type LikeAction struct {
	ID         string    `json:"id"`
	FromUserID string    `json:"fromUserId"`
	ToUserID   string    `json:"toUserId"`
	Action     string    `json:"action"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Match struct {
	ID        string    `json:"id"`
	UserA     string    `json:"userA"`
	UserB     string    `json:"userB"`
	CreatedAt time.Time `json:"createdAt"`
}

type ChatSummary struct {
	ID             string    `json:"id"`
	MatchID        string    `json:"matchId"`
	ParticipantIDs []string  `json:"participantIds"`
	LastMessage    string    `json:"lastMessage"`
	LastActiveAt   time.Time `json:"lastActiveAt"`
	UnreadCount    int       `json:"unreadCount"`
	Pinned         bool      `json:"pinned"`
}

type ChatMessage struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chatId"`
	SenderID  string    `json:"senderId"`
	Type      string    `json:"type"`
	Text      string    `json:"text"`
	GiftID    string    `json:"giftId,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type Gift struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CostSparks  int    `json:"costSparks"`
	Description string `json:"description"`
}

type SparkPackage struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Sparks      int    `json:"sparks"`
	PriceINR    int    `json:"priceInr"`
	Description string `json:"description"`
}

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Category  string    `json:"category"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type ReferralRecord struct {
	ID             string    `json:"id"`
	ReferrerUserID string    `json:"referrerUserId"`
	ReferredUserID string    `json:"referredUserId"`
	RewardSparks   int       `json:"rewardSparks"`
	Status         string    `json:"status"`
	CompletedAt    time.Time `json:"completedAt,omitempty"`
}
