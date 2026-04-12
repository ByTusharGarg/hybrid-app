package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"hybrid-app/backend/internal/domain"
	apiserver "hybrid-app/backend/internal/http"
	"hybrid-app/backend/internal/repository"
	mongorepo "hybrid-app/backend/internal/repository/mongo"
	postgresrepo "hybrid-app/backend/internal/repository/postgres"
	"hybrid-app/backend/internal/service"
)

type App struct {
	services *service.Services
	router   http.Handler
}

func New() (*App, error) {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	mongoURI := strings.TrimSpace(os.Getenv("MONGODB_URI"))
	mongoDBName := strings.TrimSpace(os.Getenv("MONGODB_DB"))
	if mongoDBName == "" {
		mongoDBName = "hybrid_app"
	}
	if databaseURL == "" || mongoURI == "" {
		return nil, errors.New("DATABASE_URL and MONGODB_URI are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgRepo, err := postgresrepo.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	mongoRepo, err := mongorepo.New(ctx, mongoURI, mongoDBName)
	if err != nil {
		return nil, err
	}

	repo := &repository.Composite{
		AuthRepo:       pgRepo,
		OnboardingRepo: pgRepo,
		DiscoveryRepo:  pgRepo,
		WalletRepo:     pgRepo,
		ChatRepo:       nil,
		ProfileRepo:    pgRepo,
	}
	repo.ChatRepo = &chatRepository{pg: pgRepo, mongo: mongoRepo}

	services := service.New(repo, newIDGenerator())
	server := apiserver.NewServer(services)

	return &App{
		services: services,
		router:   server.Router(),
	}, nil
}

func (a *App) Router() http.Handler {
	return a.router
}

func (a *App) Services() *service.Services {
	return a.services
}

type chatRepository struct {
	pg interface {
		repository.ProfileRepository
		repository.WalletRepository
		SaveChat(chat *domain.ChatSummary) error
		FindChatByID(chatID string) (*domain.ChatSummary, error)
		ListChatsForUser(userID string) ([]*domain.ChatSummary, error)
	}
	mongo interface {
		AddMessage(message domain.ChatMessage) error
		ListMessages(chatID string) ([]domain.ChatMessage, error)
	}
}

func (c *chatRepository) FindUserByID(id string) (*domain.User, error) { return c.pg.FindUserByID(id) }
func (c *chatRepository) SaveUser(user *domain.User) error             { return c.pg.SaveUser(user) }
func (c *chatRepository) ListGifts() ([]domain.Gift, error)            { return c.pg.ListGifts() }
func (c *chatRepository) AddTransaction(txn domain.WalletTransaction) error {
	return c.pg.AddTransaction(txn)
}
func (c *chatRepository) SaveChat(chat *domain.ChatSummary) error { return c.pg.SaveChat(chat) }
func (c *chatRepository) FindChatByID(chatID string) (*domain.ChatSummary, error) {
	return c.pg.FindChatByID(chatID)
}
func (c *chatRepository) ListChatsForUser(userID string) ([]*domain.ChatSummary, error) {
	return c.pg.ListChatsForUser(userID)
}
func (c *chatRepository) AddMessage(message domain.ChatMessage) error {
	return c.mongo.AddMessage(message)
}
func (c *chatRepository) ListMessages(chatID string) ([]domain.ChatMessage, error) {
	return c.mongo.ListMessages(chatID)
}
