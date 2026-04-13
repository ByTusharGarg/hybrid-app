package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"hybrid-app/backend/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func New(ctx context.Context, databaseURL string) (*Repository, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	repo := &Repository{db: db}
	if err := repo.initSchema(ctx); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) initSchema(ctx context.Context) error {
	schema := `
CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	phone TEXT UNIQUE NOT NULL,
	name TEXT NOT NULL,
	bio TEXT NOT NULL,
	gender TEXT NOT NULL,
	interested_in JSONB NOT NULL DEFAULT '[]',
	city TEXT NOT NULL,
	date_of_birth TEXT NOT NULL,
	photos JSONB NOT NULL DEFAULT '[]',
	interests JSONB NOT NULL DEFAULT '[]',
	questionnaire JSONB NOT NULL DEFAULT '{}',
	compatibility_tags JSONB NOT NULL DEFAULT '[]',
	referral_code TEXT UNIQUE NOT NULL,
	referred_by TEXT NOT NULL DEFAULT '',
	phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
	verification_status TEXT NOT NULL,
	live_video_status TEXT NOT NULL DEFAULT 'not_started',
	onboarding_completed BOOLEAN NOT NULL DEFAULT FALSE,
	profile_completion INT NOT NULL DEFAULT 0,
	sparks_balance INT NOT NULL DEFAULT 0,
	web_credits_balance INT NOT NULL DEFAULT 0,
	daily_like_quota INT NOT NULL DEFAULT 10,
	last_like_quota_reset_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	last_login_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	login_streak INT NOT NULL DEFAULT 0,
	pending_banner TEXT NOT NULL DEFAULT '',
	subscription_tier TEXT NOT NULL DEFAULT 'free',
	premium_until TIMESTAMPTZ,
	vouched_badge BOOLEAN NOT NULL DEFAULT FALSE,
	vouched_by TEXT NOT NULL DEFAULT '',
	safety_settings JSONB NOT NULL DEFAULT '{}',
	hide_from_cities JSONB NOT NULL DEFAULT '[]',
	instagram_handle TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS otp_challenges (
	id TEXT PRIMARY KEY,
	phone TEXT NOT NULL,
	code TEXT NOT NULL,
	purpose TEXT NOT NULL,
	referral_code TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMPTZ NOT NULL,
	verified_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS auth_tokens (
	token TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS wallet_transactions (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	type TEXT NOT NULL,
	amount INT NOT NULL,
	description TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL
);
CREATE TABLE IF NOT EXISTS likes (
	id TEXT PRIMARY KEY,
	from_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	to_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	action TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL
);
CREATE TABLE IF NOT EXISTS matches (
	id TEXT PRIMARY KEY,
	user_a TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	user_b TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ NOT NULL
);
CREATE TABLE IF NOT EXISTS chats (
	id TEXT PRIMARY KEY,
	match_id TEXT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
	participant_ids JSONB NOT NULL,
	last_message TEXT NOT NULL,
	last_active_at TIMESTAMPTZ NOT NULL,
	unread_count INT NOT NULL DEFAULT 0,
	pinned BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE TABLE IF NOT EXISTS notifications (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	category TEXT NOT NULL,
	title TEXT NOT NULL,
	message TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL
);
CREATE TABLE IF NOT EXISTS referrals (
	id TEXT PRIMARY KEY,
	referrer_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	referred_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	reward_sparks INT NOT NULL,
	status TEXT NOT NULL,
	completed_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS gifts (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	cost_sparks INT NOT NULL,
	description TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS spark_packages (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	sparks INT NOT NULL,
	price_inr INT NOT NULL,
	description TEXT NOT NULL
);
INSERT INTO gifts (id, name, cost_sparks, description) VALUES
	('gift_rose', 'Virtual Rose', 10, 'A sweet opener for a new match.'),
	('gift_box', 'Premium Gift Box', 50, 'A bigger gesture for active chats.'),
	('gift_ticket', 'Concert Ticket', 100, 'A premium virtual surprise.')
ON CONFLICT (id) DO NOTHING;
INSERT INTO spark_packages (id, name, sparks, price_inr, description) VALUES
	('spark_100', 'Starter Sparks', 100, 99, 'Perfect for boosts and super likes.'),
	('spark_500', 'Power Pack', 500, 399, 'Great value for active daters.'),
	('spark_1200', 'VIP Pack', 1200, 799, 'Heavy-use bundle with extra reach.')
ON CONFLICT (id) DO NOTHING;
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone_verified BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS live_video_status TEXT NOT NULL DEFAULT 'not_started';
ALTER TABLE users ADD COLUMN IF NOT EXISTS web_credits_balance INT NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS premium_until TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS vouched_badge BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS vouched_by TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS safety_settings JSONB NOT NULL DEFAULT '{}';`
	_, err := r.db.ExecContext(ctx, schema)
	return err
}

func (r *Repository) CreateOTP(challenge domain.OTPChallenge) error {
	_, err := r.db.Exec(`INSERT INTO otp_challenges (id, phone, code, purpose, referral_code, created_at, verified_at) VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (id) DO UPDATE SET phone=EXCLUDED.phone, code=EXCLUDED.code, purpose=EXCLUDED.purpose, referral_code=EXCLUDED.referral_code, created_at=EXCLUDED.created_at, verified_at=EXCLUDED.verified_at`,
		challenge.ID, challenge.Phone, challenge.Code, challenge.Purpose, challenge.Referral, challenge.CreatedAt, nullTime(challenge.VerifiedAt))
	return err
}

func (r *Repository) GetOTP(id string) (*domain.OTPChallenge, error) {
	row := r.db.QueryRow(`SELECT id, phone, code, purpose, referral_code, created_at, verified_at FROM otp_challenges WHERE id=$1`, id)
	var c domain.OTPChallenge
	var verified sql.NullTime
	if err := row.Scan(&c.ID, &c.Phone, &c.Code, &c.Purpose, &c.Referral, &c.CreatedAt, &verified); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if verified.Valid {
		c.VerifiedAt = verified.Time
	}
	return &c, nil
}

func (r *Repository) SaveOTP(challenge *domain.OTPChallenge) error { return r.CreateOTP(*challenge) }

func (r *Repository) SaveToken(token, userID string) error {
	_, err := r.db.Exec(`INSERT INTO auth_tokens (token, user_id) VALUES ($1,$2) ON CONFLICT (token) DO UPDATE SET user_id=EXCLUDED.user_id`, token, userID)
	return err
}

func (r *Repository) FindUserIDByToken(token string) (string, error) {
	row := r.db.QueryRow(`SELECT user_id FROM auth_tokens WHERE token=$1`, token)
	var userID string
	if err := row.Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return userID, nil
}

func (r *Repository) FindUserByPhone(phone string) (*domain.User, error) {
	return r.findUser(`SELECT `+userColumns()+` FROM users WHERE phone=$1`, phone)
}

func (r *Repository) FindUserByReferralCode(code string) (*domain.User, error) {
	return r.findUser(`SELECT `+userColumns()+` FROM users WHERE referral_code=$1`, code)
}

func (r *Repository) FindUserByID(id string) (*domain.User, error) {
	return r.findUser(`SELECT `+userColumns()+` FROM users WHERE id=$1`, id)
}

func (r *Repository) SaveUser(user *domain.User) error {
	_, err := r.db.Exec(`INSERT INTO users (
		id, phone, name, bio, gender, interested_in, city, date_of_birth, photos, interests, questionnaire, compatibility_tags,
		referral_code, referred_by, phone_verified, verification_status, live_video_status, onboarding_completed, profile_completion,
		sparks_balance, web_credits_balance, daily_like_quota, last_like_quota_reset_at, last_login_at, login_streak, pending_banner,
		subscription_tier, premium_until, vouched_badge, vouched_by, safety_settings, hide_from_cities, instagram_handle, created_at, updated_at
	) VALUES (
		$1,$2,$3,$4,$5,$6::jsonb,$7,$8,$9::jsonb,$10::jsonb,$11::jsonb,$12::jsonb,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30::jsonb,$31::jsonb,$32,$33,$34
	) ON CONFLICT (id) DO UPDATE SET
		phone=EXCLUDED.phone, name=EXCLUDED.name, bio=EXCLUDED.bio, gender=EXCLUDED.gender,
		interested_in=EXCLUDED.interested_in, city=EXCLUDED.city, date_of_birth=EXCLUDED.date_of_birth,
		photos=EXCLUDED.photos, interests=EXCLUDED.interests, questionnaire=EXCLUDED.questionnaire,
		compatibility_tags=EXCLUDED.compatibility_tags, referral_code=EXCLUDED.referral_code, referred_by=EXCLUDED.referred_by,
		phone_verified=EXCLUDED.phone_verified, verification_status=EXCLUDED.verification_status,
		live_video_status=EXCLUDED.live_video_status, onboarding_completed=EXCLUDED.onboarding_completed,
		profile_completion=EXCLUDED.profile_completion, sparks_balance=EXCLUDED.sparks_balance,
		web_credits_balance=EXCLUDED.web_credits_balance, daily_like_quota=EXCLUDED.daily_like_quota,
		last_like_quota_reset_at=EXCLUDED.last_like_quota_reset_at, last_login_at=EXCLUDED.last_login_at, login_streak=EXCLUDED.login_streak,
		pending_banner=EXCLUDED.pending_banner, subscription_tier=EXCLUDED.subscription_tier, premium_until=EXCLUDED.premium_until,
		vouched_badge=EXCLUDED.vouched_badge, vouched_by=EXCLUDED.vouched_by, safety_settings=EXCLUDED.safety_settings,
		hide_from_cities=EXCLUDED.hide_from_cities, instagram_handle=EXCLUDED.instagram_handle, created_at=EXCLUDED.created_at, updated_at=EXCLUDED.updated_at`,
		user.ID, user.Phone, user.Name, user.Bio, user.Gender,
		mustJSON(user.InterestedIn), user.City, user.DateOfBirth, mustJSON(user.Photos), mustJSON(user.Interests),
		mustJSON(user.Questionnaire), mustJSON(user.CompatibilityTags), user.ReferralCode, user.ReferredBy, user.PhoneVerified,
		user.VerificationStatus, user.LiveVideoStatus, user.OnboardingCompleted, user.ProfileCompletion, user.SparksBalance,
		user.WebCreditsBalance, user.DailyLikeQuota, user.LastLikeQuotaResetAt, user.LastLoginAt, user.LoginStreak, user.PendingBanner,
		user.SubscriptionTier, nullTime(user.PremiumUntil), user.VouchedBadge, user.VouchedBy, mustJSON(user.SafetySettings),
		mustJSON(user.HideFromCities), user.InstagramHandle, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *Repository) ListOtherCompletedUsers(excludeUserID string) ([]*domain.User, error) {
	rows, err := r.db.Query(`SELECT `+userColumns()+` FROM users WHERE onboarding_completed=TRUE AND id <> $1`, excludeUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*domain.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *Repository) AddTransaction(txn domain.WalletTransaction) error {
	_, err := r.db.Exec(`INSERT INTO wallet_transactions (id, user_id, type, amount, description, created_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		txn.ID, txn.UserID, txn.Type, txn.Amount, txn.Description, txn.CreatedAt)
	return err
}

func (r *Repository) ListTransactions(userID string) ([]domain.WalletTransaction, error) {
	rows, err := r.db.Query(`SELECT id, user_id, type, amount, description, created_at FROM wallet_transactions WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var txns []domain.WalletTransaction
	for rows.Next() {
		var txn domain.WalletTransaction
		if err := rows.Scan(&txn.ID, &txn.UserID, &txn.Type, &txn.Amount, &txn.Description, &txn.CreatedAt); err != nil {
			return nil, err
		}
		txns = append(txns, txn)
	}
	return txns, rows.Err()
}

func (r *Repository) AddLike(like domain.LikeAction) error {
	_, err := r.db.Exec(`INSERT INTO likes (id, from_user_id, to_user_id, action, created_at) VALUES ($1,$2,$3,$4,$5)`, like.ID, like.FromUserID, like.ToUserID, like.Action, like.CreatedAt)
	return err
}

func (r *Repository) ListLikesForTarget(userID string) ([]domain.LikeAction, error) {
	rows, err := r.db.Query(`SELECT id, from_user_id, to_user_id, action, created_at FROM likes WHERE to_user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var likes []domain.LikeAction
	for rows.Next() {
		var like domain.LikeAction
		if err := rows.Scan(&like.ID, &like.FromUserID, &like.ToUserID, &like.Action, &like.CreatedAt); err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}
	return likes, rows.Err()
}

func (r *Repository) HasPositiveLike(fromUserID, toUserID string) (bool, error) {
	row := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM likes WHERE from_user_id=$1 AND to_user_id=$2 AND action IN ('like','super_like'))`, fromUserID, toUserID)
	var exists bool
	return exists, row.Scan(&exists)
}

func (r *Repository) SaveMatch(match domain.Match) error {
	_, err := r.db.Exec(`INSERT INTO matches (id, user_a, user_b, created_at) VALUES ($1,$2,$3,$4) ON CONFLICT (id) DO UPDATE SET user_a=EXCLUDED.user_a, user_b=EXCLUDED.user_b, created_at=EXCLUDED.created_at`,
		match.ID, match.UserA, match.UserB, match.CreatedAt)
	return err
}

func (r *Repository) FindMatchBetween(userA, userB string) (*domain.Match, error) {
	row := r.db.QueryRow(`SELECT id, user_a, user_b, created_at FROM matches WHERE (user_a=$1 AND user_b=$2) OR (user_a=$2 AND user_b=$1) LIMIT 1`, userA, userB)
	var match domain.Match
	if err := row.Scan(&match.ID, &match.UserA, &match.UserB, &match.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &match, nil
}

func (r *Repository) ListMatchesForUser(userID string) ([]domain.Match, error) {
	rows, err := r.db.Query(`SELECT id, user_a, user_b, created_at FROM matches WHERE user_a=$1 OR user_b=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var matches []domain.Match
	for rows.Next() {
		var match domain.Match
		if err := rows.Scan(&match.ID, &match.UserA, &match.UserB, &match.CreatedAt); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, rows.Err()
}

func (r *Repository) SaveChat(chat *domain.ChatSummary) error {
	_, err := r.db.Exec(`INSERT INTO chats (id, match_id, participant_ids, last_message, last_active_at, unread_count, pinned)
	VALUES ($1,$2,$3::jsonb,$4,$5,$6,$7)
	ON CONFLICT (id) DO UPDATE SET match_id=EXCLUDED.match_id, participant_ids=EXCLUDED.participant_ids, last_message=EXCLUDED.last_message, last_active_at=EXCLUDED.last_active_at, unread_count=EXCLUDED.unread_count, pinned=EXCLUDED.pinned`,
		chat.ID, chat.MatchID, mustJSON(chat.ParticipantIDs), chat.LastMessage, chat.LastActiveAt, chat.UnreadCount, chat.Pinned)
	return err
}

func (r *Repository) FindChatByID(chatID string) (*domain.ChatSummary, error) {
	return r.findChat(`SELECT id, match_id, participant_ids, last_message, last_active_at, unread_count, pinned FROM chats WHERE id=$1`, chatID)
}

func (r *Repository) FindChatByMatchID(matchID string) (*domain.ChatSummary, error) {
	return r.findChat(`SELECT id, match_id, participant_ids, last_message, last_active_at, unread_count, pinned FROM chats WHERE match_id=$1`, matchID)
}

func (r *Repository) ListChatsForUser(userID string) ([]*domain.ChatSummary, error) {
	rows, err := r.db.Query(`SELECT id, match_id, participant_ids, last_message, last_active_at, unread_count, pinned FROM chats ORDER BY last_active_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chats []*domain.ChatSummary
	for rows.Next() {
		chat, err := scanChat(rows)
		if err != nil {
			return nil, err
		}
		for _, participantID := range chat.ParticipantIDs {
			if participantID == userID {
				chats = append(chats, chat)
				break
			}
		}
	}
	return chats, rows.Err()
}

func (r *Repository) AddNotification(notification domain.Notification) error {
	_, err := r.db.Exec(`INSERT INTO notifications (id, user_id, category, title, message, created_at) VALUES ($1,$2,$3,$4,$5,$6)`, notification.ID, notification.UserID, notification.Category, notification.Title, notification.Message, notification.CreatedAt)
	return err
}

func (r *Repository) ListNotifications(userID string) ([]domain.Notification, error) {
	rows, err := r.db.Query(`SELECT id, user_id, category, title, message, created_at FROM notifications WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []domain.Notification
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Category, &n.Title, &n.Message, &n.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, n)
	}
	return list, rows.Err()
}

func (r *Repository) AddReferralRecord(record domain.ReferralRecord) error {
	_, err := r.db.Exec(`INSERT INTO referrals (id, referrer_user_id, referred_user_id, reward_sparks, status, completed_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		record.ID, record.ReferrerUserID, record.ReferredUserID, record.RewardSparks, record.Status, nullTime(record.CompletedAt))
	return err
}

func (r *Repository) ListReferralRecords(referrerUserID string) ([]domain.ReferralRecord, error) {
	rows, err := r.db.Query(`SELECT id, referrer_user_id, referred_user_id, reward_sparks, status, completed_at FROM referrals WHERE referrer_user_id=$1 ORDER BY completed_at DESC NULLS LAST`, referrerUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []domain.ReferralRecord
	for rows.Next() {
		var record domain.ReferralRecord
		var completed sql.NullTime
		if err := rows.Scan(&record.ID, &record.ReferrerUserID, &record.ReferredUserID, &record.RewardSparks, &record.Status, &completed); err != nil {
			return nil, err
		}
		if completed.Valid {
			record.CompletedAt = completed.Time
		}
		list = append(list, record)
	}
	return list, rows.Err()
}

func (r *Repository) ListGifts() ([]domain.Gift, error) {
	rows, err := r.db.Query(`SELECT id, name, cost_sparks, description FROM gifts ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var gifts []domain.Gift
	for rows.Next() {
		var gift domain.Gift
		if err := rows.Scan(&gift.ID, &gift.Name, &gift.CostSparks, &gift.Description); err != nil {
			return nil, err
		}
		gifts = append(gifts, gift)
	}
	return gifts, rows.Err()
}

func (r *Repository) ListSparkPackages() ([]domain.SparkPackage, error) {
	rows, err := r.db.Query(`SELECT id, name, sparks, price_inr, description FROM spark_packages ORDER BY price_inr`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []domain.SparkPackage
	for rows.Next() {
		var pkg domain.SparkPackage
		if err := rows.Scan(&pkg.ID, &pkg.Name, &pkg.Sparks, &pkg.PriceINR, &pkg.Description); err != nil {
			return nil, err
		}
		list = append(list, pkg)
	}
	return list, rows.Err()
}

func userColumns() string {
	return `id, phone, name, bio, gender, interested_in, city, date_of_birth, photos, interests, questionnaire, compatibility_tags, referral_code, referred_by, phone_verified, verification_status, live_video_status, onboarding_completed, profile_completion, sparks_balance, web_credits_balance, daily_like_quota, last_like_quota_reset_at, last_login_at, login_streak, pending_banner, subscription_tier, premium_until, vouched_badge, vouched_by, safety_settings, hide_from_cities, instagram_handle, created_at, updated_at`
}

func (r *Repository) findUser(query string, arg any) (*domain.User, error) {
	row := r.db.QueryRow(query, arg)
	return scanUser(row)
}

func scanUser(scanner interface{ Scan(...any) error }) (*domain.User, error) {
	var user domain.User
	var interestedIn, photos, interests, questionnaire, compatibilityTags, safetySettings, hideFromCities []byte
	var premiumUntil sql.NullTime
	if err := scanner.Scan(
		&user.ID, &user.Phone, &user.Name, &user.Bio, &user.Gender,
		&interestedIn, &user.City, &user.DateOfBirth, &photos, &interests, &questionnaire, &compatibilityTags,
		&user.ReferralCode, &user.ReferredBy, &user.PhoneVerified, &user.VerificationStatus, &user.LiveVideoStatus,
		&user.OnboardingCompleted, &user.ProfileCompletion, &user.SparksBalance, &user.WebCreditsBalance, &user.DailyLikeQuota,
		&user.LastLikeQuotaResetAt, &user.LastLoginAt, &user.LoginStreak, &user.PendingBanner, &user.SubscriptionTier,
		&premiumUntil, &user.VouchedBadge, &user.VouchedBy, &safetySettings, &hideFromCities, &user.InstagramHandle, &user.CreatedAt, &user.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(interestedIn, &user.InterestedIn); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(photos, &user.Photos); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(interests, &user.Interests); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(questionnaire, &user.Questionnaire); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(compatibilityTags, &user.CompatibilityTags); err != nil {
		return nil, err
	}
	if premiumUntil.Valid {
		user.PremiumUntil = premiumUntil.Time
	}
	if err := json.Unmarshal(safetySettings, &user.SafetySettings); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(hideFromCities, &user.HideFromCities); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) findChat(query string, arg any) (*domain.ChatSummary, error) {
	row := r.db.QueryRow(query, arg)
	return scanChat(row)
}

func scanChat(scanner interface{ Scan(...any) error }) (*domain.ChatSummary, error) {
	var chat domain.ChatSummary
	var participantIDs []byte
	if err := scanner.Scan(&chat.ID, &chat.MatchID, &participantIDs, &chat.LastMessage, &chat.LastActiveAt, &chat.UnreadCount, &chat.Pinned); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(participantIDs, &chat.ParticipantIDs); err != nil {
		return nil, err
	}
	return &chat, nil
}

func mustJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func nullTime(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return t
}
