package http

import "hybrid-app/backend/internal/domain"

func welcomeCopy(referred bool) string {
	if referred {
		return "Congratulations! You just earned 100 Sparks welcome bonus + 50 referral bonus."
	}
	return "Congratulations! You just earned 100 Sparks welcome bonus."
}

func onboardingState(user *domain.User) map[string]any {
	return map[string]any{
		"completed":         user.OnboardingCompleted,
		"profileCompletion": user.ProfileCompletion,
		"nextRequiredStep":  nextRequiredStep(user),
		"verification":      user.VerificationStatus,
	}
}

func nextRequiredStep(user *domain.User) string {
	if len(user.Questionnaire) == 0 {
		return "questionnaire"
	}
	if user.VerificationStatus != "verified" && user.LiveVideoStatus != "verified" {
		return "live_video_verification"
	}
	if !user.OnboardingCompleted {
		return "profile_setup"
	}
	return "complete"
}

func ternaryInt(condition bool, left, right int) int {
	if condition {
		return left
	}
	return right
}
