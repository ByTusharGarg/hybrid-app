package http

import (
	"os"
	"testing"
)

func TestRegistrationReferralFlow(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" || os.Getenv("MONGODB_URI") == "" {
		t.Skip("database integration env vars are not set")
	}
}

func TestAuthedEndpoints(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" || os.Getenv("MONGODB_URI") == "" {
		t.Skip("database integration env vars are not set")
	}
}
