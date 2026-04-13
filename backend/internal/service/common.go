package service

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"hybrid-app/backend/internal/domain"
	"hybrid-app/backend/internal/repository"
)

var (
	ErrInvalidOTP         = errors.New("invalid otp")
	ErrChallengeNotFound  = errors.New("otp challenge not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrNotEnoughSparks    = errors.New("not enough sparks")
	ErrUnauthorized       = errors.New("invalid token")
	ErrChatNotFound       = errors.New("chat not found")
	ErrTargetUserNotFound = errors.New("target user not found")
	ErrUserNotRegistered  = errors.New("user not registered")
	ErrUnsupportedAction  = errors.New("unsupported action")
	ErrNotEnoughCredits   = errors.New("not enough web credits")
)

type baseService struct {
	repo  repository.Repository
	newID func(string) string
}

func sameDay(a, b time.Time) bool {
	if a.IsZero() {
		return false
	}
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func addNotification(repo repository.Repository, idGen func(string) string, userID, category, title, message string) error {
	return repo.AddNotification(domain.Notification{
		ID:        idGen("notif"),
		UserID:    userID,
		Category:  category,
		Title:     title,
		Message:   message,
		CreatedAt: time.Now().UTC(),
	})
}

func addTransaction(repo repository.Repository, idGen func(string) string, userID, txnType string, amount int, description string) (domain.WalletTransaction, error) {
	txn := domain.WalletTransaction{
		ID:          idGen("txn"),
		UserID:      userID,
		Type:        txnType,
		Amount:      amount,
		Description: description,
		CreatedAt:   time.Now().UTC(),
	}
	return txn, repo.AddTransaction(txn)
}

func tagsFromQuestionnaire(answers map[string]any) []string {
	tags := []string{}
	for _, key := range []string{"relationship_goal", "weekend", "energy", "love_language"} {
		if value, ok := answers[key].(string); ok && value != "" {
			tags = append(tags, value)
		}
	}
	return tags
}

func hasTransactionType(txns []domain.WalletTransaction, txnType string) bool {
	for _, txn := range txns {
		if txn.Type == txnType {
			return true
		}
	}
	return false
}

func compatibility(a, b *domain.User) int {
	score := 50
	for _, interest := range a.Interests {
		if slices.Contains(b.Interests, interest) {
			score += 10
		}
	}
	for _, tag := range a.CompatibilityTags {
		if slices.Contains(b.CompatibilityTags, tag) {
			score += 8
		}
	}
	if a.City == b.City {
		score += 6
	}
	if score > 98 {
		score = 98
	}
	return score
}

func ageFromDOB(dob string) int {
	if dob == "" {
		return 0
	}
	parsed, err := time.Parse("2006-01-02", dob)
	if err != nil {
		return 0
	}
	now := time.Now().UTC()
	age := now.Year() - parsed.Year()
	if now.YearDay() < parsed.YearDay() {
		age--
	}
	return age
}

func humanizeLastActive(t time.Time) string {
	if t.IsZero() {
		return "new here"
	}
	diff := time.Since(t)
	if diff < time.Hour {
		return "active now"
	}
	if diff < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	}
	return fmt.Sprintf("%dd ago", int(diff.Hours()/24))
}
