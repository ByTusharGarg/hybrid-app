package service

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"hybrid-app/backend/internal/domain"
	"hybrid-app/backend/internal/repository"
)

type DiscoveryService struct{ baseService }

func NewDiscoveryService(repo repository.Repository, idGen func(string) string) *DiscoveryService {
	return &DiscoveryService{baseService{repo: repo, newID: idGen}}
}

func (s *DiscoveryService) Discover(userID string) ([]domain.DiscoveryProfile, []domain.DiscoveryProfile, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, ErrUserNotFound
	}
	others, err := s.repo.ListOtherCompletedUsers(userID)
	if err != nil {
		return nil, nil, err
	}
	topPicks, feed := buildDiscovery(user, others)
	return topPicks, feed, nil
}

func (s *DiscoveryService) ApplyAction(userID, targetID, action string) (map[string]any, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	target, err := s.repo.FindUserByID(targetID)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, ErrTargetUserNotFound
	}

	now := time.Now().UTC()
	if !sameDay(user.LastLikeQuotaResetAt, now) {
		user.DailyLikeQuota = domain.DefaultDailyLikeCap
		user.LastLikeQuotaResetAt = now
	}

	switch action {
	case "like":
		if user.SubscriptionTier != "premium" && user.DailyLikeQuota <= 0 {
			return map[string]any{"status": "limit_reached", "message": "Daily like limit reached. Refill with Sparks.", "refillCostSparks": domain.LikeRefillCost}, nil
		}
		if user.SubscriptionTier != "premium" {
			user.DailyLikeQuota--
		}
	case "super_like":
		if user.SparksBalance < domain.SuperLikeCost {
			return nil, ErrNotEnoughSparks
		}
		user.SparksBalance -= domain.SuperLikeCost
		if _, err := addTransaction(s.repo, s.newID, user.ID, "super_like", -domain.SuperLikeCost, fmt.Sprintf("Super Like sent to %s", target.Name)); err != nil {
			return nil, err
		}
	case "undo":
		if user.SparksBalance < domain.UndoCost {
			return nil, ErrNotEnoughSparks
		}
		user.SparksBalance -= domain.UndoCost
		if _, err := addTransaction(s.repo, s.newID, user.ID, "undo_swipe", -domain.UndoCost, "Undo last swipe"); err != nil {
			return nil, err
		}
	case "pass":
	default:
		return nil, ErrUnsupportedAction
	}

	user.UpdatedAt = now
	if err := s.repo.SaveUser(user); err != nil {
		return nil, err
	}
	if err := s.repo.AddLike(domain.LikeAction{ID: s.newID("like"), FromUserID: userID, ToUserID: targetID, Action: action, CreatedAt: now}); err != nil {
		return nil, err
	}

	response := map[string]any{"status": "ok", "remainingLikes": user.DailyLikeQuota}
	if action == "like" || action == "super_like" {
		if err := addNotification(s.repo, s.newID, targetID, "likes", "Someone liked you", fmt.Sprintf("%s sent you a %s.", user.Name, strings.ReplaceAll(action, "_", " "))); err != nil {
			return nil, err
		}
		hasMutual, err := s.repo.HasPositiveLike(targetID, userID)
		if err != nil {
			return nil, err
		}
		if hasMutual {
			match, err := s.repo.FindMatchBetween(userID, targetID)
			if err != nil {
				return nil, err
			}
			if match == nil {
				newMatch := domain.Match{ID: s.newID("match"), UserA: userID, UserB: targetID, CreatedAt: now}
				if err := s.repo.SaveMatch(newMatch); err != nil {
					return nil, err
				}
				match = &newMatch
			}
			chat, err := s.repo.FindChatByMatchID(match.ID)
			if err != nil {
				return nil, err
			}
			if chat == nil {
				chat = &domain.ChatSummary{
					ID:             s.newID("chat"),
					MatchID:        match.ID,
					ParticipantIDs: []string{userID, targetID},
					LastMessage:    "It's a match! Start the conversation.",
					LastActiveAt:   now,
				}
				if err := s.repo.SaveChat(chat); err != nil {
					return nil, err
				}
			}
			if err := addNotification(s.repo, s.newID, userID, "messages", "It's a Match!", fmt.Sprintf("You and %s liked each other.", target.Name)); err != nil {
				return nil, err
			}
			if err := addNotification(s.repo, s.newID, targetID, "messages", "It's a Match!", fmt.Sprintf("You and %s liked each other.", user.Name)); err != nil {
				return nil, err
			}
			response["match"] = map[string]any{"id": match.ID, "message": "It's a Match!", "chatId": chat.ID}
		}
	}
	return response, nil
}

func (s *DiscoveryService) Matches(userID string) ([]map[string]any, error) {
	matches, err := s.repo.ListMatchesForUser(userID)
	if err != nil {
		return nil, err
	}
	var result []map[string]any
	for _, match := range matches {
		otherID := match.UserA
		if otherID == userID {
			otherID = match.UserB
		}
		other, err := s.repo.FindUserByID(otherID)
		if err != nil {
			return nil, err
		}
		chat, err := s.repo.FindChatByMatchID(match.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, map[string]any{"id": match.ID, "user": other, "createdAt": match.CreatedAt, "chatExists": chat != nil})
	}
	return result, nil
}

func (s *DiscoveryService) LikesYou(userID string, reveal bool) ([]map[string]any, error) {
	user, err := s.repo.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if reveal && user.SubscriptionTier != "premium" {
		if user.SparksBalance < domain.LikesRevealCost {
			return nil, ErrNotEnoughSparks
		}
		user.SparksBalance -= domain.LikesRevealCost
		user.UpdatedAt = time.Now().UTC()
		if err := s.repo.SaveUser(user); err != nil {
			return nil, err
		}
		if _, err := addTransaction(s.repo, s.newID, user.ID, "likes_reveal", -domain.LikesRevealCost, "Revealed users who liked you"); err != nil {
			return nil, err
		}
	}

	likes, err := s.repo.ListLikesForTarget(userID)
	if err != nil {
		return nil, err
	}
	var results []map[string]any
	for _, like := range likes {
		if like.Action != "like" && like.Action != "super_like" {
			continue
		}
		source, err := s.repo.FindUserByID(like.FromUserID)
		if err != nil {
			return nil, err
		}
		entry := map[string]any{"likedAt": like.CreatedAt, "teaser": source.City}
		if reveal || user.SubscriptionTier == "premium" {
			entry["user"] = source
		}
		results = append(results, entry)
	}
	slices.SortFunc(results, func(a, b map[string]any) int {
		at := a["likedAt"].(time.Time)
		bt := b["likedAt"].(time.Time)
		if at.After(bt) {
			return -1
		}
		if at.Before(bt) {
			return 1
		}
		return 0
	})
	return results, nil
}
