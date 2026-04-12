# Go Backend

This folder now contains a lightweight Go API for the dating app flows described in the product brief.

## Architecture

The backend is now arranged in layers:

- `internal/http`: controllers
- `internal/service`: business logic
- `internal/repository`: repository contracts
- `internal/repository/memory`: current in-memory adapter

Recommended production storage split:

- `PostgreSQL`: users, auth/session metadata, wallet balances, wallet transactions, referrals, likes, matches, purchases, chat room metadata
- `MongoDB`: profile documents, questionnaire answers, discovery projections, chat messages, high-volume activity documents

See [internal/service/architecture.md](/Users/devmunjal/Desktop/per/hybrid-app/backend/internal/service/architecture.md) for the rationale.

## What it covers

- Phone registration and login with mock OTP verification
- Referral code entry during registration
- Onboarding steps:
  - questionnaire
  - gender verification
  - profile setup
- Sparks wallet:
  - welcome bonus
  - referral bonus
  - daily login streak reward
  - boost spend
  - super like spend
  - undo spend
  - like refill
  - gift sending
- Discover feed, top picks, likes, matches, chats
- Activity center and referral center
- Profile fetch and update

## Run

```bash
cd backend
go run ./cmd/api
```

Server default address: `:8080`

Required env vars:

- `DATABASE_URL`
- `MONGODB_URI`
- `MONGODB_DB` optional, defaults to `hybrid_app`

## Docker Quick Start

```bash
cd backend
docker compose up --build
```

Backend:

- API: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/swagger`
- OpenAPI JSON: `http://localhost:8080/openapi.json`

Databases started by compose:

- PostgreSQL: `localhost:5432`
- MongoDB: `localhost:27017`

## Mock OTP

Every OTP challenge returns `123456` so frontend integration is easy during development.

## Main routes

- `POST /api/v1/auth/register/start`
- `POST /api/v1/auth/register/verify-otp`
- `POST /api/v1/auth/login/start`
- `POST /api/v1/auth/login/verify-otp`
- `POST /api/v1/onboarding/questionnaire`
- `POST /api/v1/onboarding/gender-verification`
- `POST /api/v1/onboarding/profile`
- `GET /api/v1/home`
- `GET /api/v1/discover`
- `POST /api/v1/discover/actions`
- `GET /api/v1/matches`
- `GET /api/v1/likes-you`
- `GET /api/v1/chats`
- `GET /api/v1/chats/{chatId}`
- `POST /api/v1/chats/{chatId}/messages`
- `GET /api/v1/activity`
- `GET /api/v1/wallet`
- `POST /api/v1/wallet/daily-login`
- `POST /api/v1/wallet/boost`
- `POST /api/v1/wallet/like-refill`
- `GET /api/v1/referrals`
- `GET /api/v1/me`
- `PATCH /api/v1/me`
