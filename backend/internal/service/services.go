package service

import "hybrid-app/backend/internal/repository"

type Services struct {
	Auth       *AuthService
	Onboarding *OnboardingService
	Discovery  *DiscoveryService
	Wallet     *WalletService
	Chat       *ChatService
	Profile    *ProfileService
	Ecosystem  *EcosystemService
}

func New(repo repository.Repository, idGen func(string) string) *Services {
	return &Services{
		Auth:       NewAuthService(repo, idGen),
		Onboarding: NewOnboardingService(repo, idGen),
		Discovery:  NewDiscoveryService(repo, idGen),
		Wallet:     NewWalletService(repo, idGen),
		Chat:       NewChatService(repo, idGen),
		Profile:    NewProfileService(repo),
		Ecosystem:  NewEcosystemService(repo, idGen),
	}
}
