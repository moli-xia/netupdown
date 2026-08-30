package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/moli-xia/netupdown/internal/model"
	"github.com/moli-xia/netupdown/internal/pkg/apperr"
	"github.com/moli-xia/netupdown/internal/repo"
	"gorm.io/gorm"
)

type Service struct {
	store    *repo.Store
	key      []byte
	issuer   string
	mu       sync.Mutex
	failures map[string]*loginAttempt
}
type loginAttempt struct {
	Count       int
	LockedUntil time.Time
}
type Claims struct {
	Role int8 `json:"role"`
	jwt.RegisteredClaims
}

type PreviewClaims struct {
	AppID uint64 `json:"app_id"`
	Slug  string `json:"slug"`
	jwt.RegisteredClaims
}

func New(store *repo.Store, key []byte) *Service {
	return &Service{store: store, key: key, issuer: "netupdown", failures: map[string]*loginAttempt{}}
}
func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, &argon2id.Params{Memory: 64 * 1024, Iterations: 3, Parallelism: 2, SaltLength: 16, KeyLength: 32})
}
func (s *Service) Authenticate(ctx context.Context, username, password, ip, ua string) (model.User, string, string, error) {
	var user model.User
	if s.isLocked(username) {
		return user, "", "", apperr.New(10009, 429, "账号暂时锁定，请 15 分钟后重试")
	}
	if err := s.store.DB.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		s.recordFailure(username)
		return user, "", "", apperr.Unauthorized
	}
	ok, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil || !ok || user.Status != 1 {
		s.recordFailure(username)
		return user, "", "", apperr.Unauthorized
	}
	s.clearFailures(username)
	now := time.Now().UTC()
	user.LastLoginAt = &now
	user.LastLoginIP = ip
	_ = s.store.DB.WithContext(ctx).Model(&user).Updates(map[string]any{"last_login_at": now, "last_login_ip": ip}).Error
	access, err := s.AccessToken(user)
	if err != nil {
		return user, "", "", err
	}
	refresh, err := s.NewRefresh(ctx, user.ID, ip, ua)
	return user, access, refresh, err
}
func (s *Service) isLocked(username string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	a := s.failures[username]
	if a == nil {
		return false
	}
	if !a.LockedUntil.IsZero() && time.Now().Before(a.LockedUntil) {
		return true
	}
	if !a.LockedUntil.IsZero() {
		delete(s.failures, username)
	}
	return false
}
func (s *Service) recordFailure(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a := s.failures[username]
	if a == nil {
		a = &loginAttempt{}
		s.failures[username] = a
	}
	a.Count++
	if a.Count >= 10 {
		a.LockedUntil = time.Now().Add(15 * time.Minute)
	}
}
func (s *Service) clearFailures(username string) {
	s.mu.Lock()
	delete(s.failures, username)
	s.mu.Unlock()
}
func (s *Service) AccessToken(user model.User) (string, error) {
	now := time.Now().UTC()
	claims := Claims{Role: user.Role, RegisteredClaims: jwt.RegisteredClaims{Issuer: s.issuer, Subject: strconv.FormatUint(user.ID, 10), ID: randomString(16), IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Hour))}}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.key)
}
func (s *Service) Parse(token string) (*Claims, error) {
	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.key, nil
	}, jwt.WithIssuer(s.issuer))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, apperr.Expired
		}
		return nil, apperr.Unauthorized
	}
	if !parsed.Valid {
		return nil, apperr.Unauthorized
	}
	return claims, nil
}

func (s *Service) PreviewToken(appID uint64, slug string) (string, error) {
	now := time.Now().UTC()
	claims := PreviewClaims{
		AppID: appID,
		Slug:  slug,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Audience:  []string{"preview"},
			Subject:   strconv.FormatUint(appID, 10),
			ID:        randomString(16),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.key)
}

func (s *Service) ParsePreview(token, slug string) (*PreviewClaims, error) {
	claims := &PreviewClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.key, nil
	}, jwt.WithIssuer(s.issuer), jwt.WithAudience("preview"))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, apperr.Expired
		}
		return nil, apperr.Unauthorized
	}
	if !parsed.Valid || claims.Slug != slug || claims.AppID == 0 {
		return nil, apperr.Unauthorized
	}
	return claims, nil
}
func (s *Service) NewRefresh(ctx context.Context, userID uint64, ip, ua string) (string, error) {
	raw := randomString(32)
	h := sha256.Sum256([]byte(raw))
	row := model.UserToken{UserID: userID, TokenHash: fmt.Sprintf("%x", h[:]), IP: ip, UA: ua, ExpiresAt: time.Now().UTC().Add(30 * 24 * time.Hour)}
	return raw, s.store.DB.WithContext(ctx).Create(&row).Error
}
func (s *Service) Rotate(ctx context.Context, raw, ip, ua string) (model.User, string, string, error) {
	var user model.User
	h := sha256.Sum256([]byte(raw))
	hash := fmt.Sprintf("%x", h[:])
	var row model.UserToken
	err := s.store.DB.WithContext(ctx).Where("token_hash = ?", hash).First(&row).Error
	if err != nil {
		return user, "", "", apperr.Unauthorized
	}
	if row.RevokedAt != nil {
		_ = s.RevokeAll(ctx, row.UserID)
		return user, "", "", apperr.Unauthorized
	}
	if row.ExpiresAt.Before(time.Now().UTC()) {
		return user, "", "", apperr.Unauthorized
	}
	if err := s.store.DB.WithContext(ctx).First(&user, row.UserID).Error; err != nil {
		return user, "", "", apperr.Unauthorized
	}
	var access, next string
	err = s.store.Tx(func(tx *repo.Store) error {
		now := time.Now().UTC()
		if err := tx.DB.WithContext(ctx).Model(&row).Update("revoked_at", now).Error; err != nil {
			return err
		}
		temp := New(tx, s.key)
		var e error
		next, e = temp.NewRefresh(ctx, user.ID, ip, ua)
		if e != nil {
			return e
		}
		access, e = temp.AccessToken(user)
		return e
	})
	return user, access, next, err
}
func (s *Service) Revoke(ctx context.Context, raw string) error {
	h := sha256.Sum256([]byte(raw))
	now := time.Now().UTC()
	return s.store.DB.WithContext(ctx).Model(&model.UserToken{}).Where("token_hash = ?", fmt.Sprintf("%x", h[:])).Update("revoked_at", now).Error
}
func (s *Service) RevokeAll(ctx context.Context, userID uint64) error {
	now := time.Now().UTC()
	return s.store.DB.WithContext(ctx).Model(&model.UserToken{}).Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", now).Error
}
func (s *Service) User(ctx context.Context, id uint64) (model.User, error) {
	var u model.User
	err := s.store.DB.WithContext(ctx).First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return u, apperr.NotFound
	}
	return u, err
}
func (s *Service) UpdateProfile(ctx context.Context, id uint64, nickname, email, avatar string) (model.User, error) {
	var user model.User
	if err := s.store.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return user, apperr.NotFound
	}
	if err := s.store.DB.WithContext(ctx).Model(&user).Updates(map[string]any{"nickname": nickname, "email": email, "avatar": avatar}).Error; err != nil {
		return user, err
	}
	return s.User(ctx, id)
}
func (s *Service) ChangePassword(ctx context.Context, id uint64, oldPassword, newPassword string) error {
	if len([]rune(newPassword)) < 12 {
		return apperr.BadRequest
	}
	var user model.User
	if err := s.store.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return apperr.NotFound
	}
	ok, err := argon2id.ComparePasswordAndHash(oldPassword, user.PasswordHash)
	if err != nil || !ok {
		return apperr.Unauthorized
	}
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := s.store.DB.WithContext(ctx).Model(&user).Update("password_hash", hash).Error; err != nil {
		return err
	}
	return s.RevokeAll(ctx, id)
}
func (s *Service) Sessions(ctx context.Context, userID uint64) ([]model.UserToken, error) {
	var rows []model.UserToken
	err := s.store.DB.WithContext(ctx).Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now().UTC()).Order("created_at DESC").Find(&rows).Error
	for i := range rows {
		rows[i].TokenHash = ""
	}
	return rows, err
}
func (s *Service) RevokeSession(ctx context.Context, userID, sessionID uint64) error {
	now := time.Now().UTC()
	result := s.store.DB.WithContext(ctx).Model(&model.UserToken{}).Where("id = ? AND user_id = ?", sessionID, userID).Update("revoked_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperr.NotFound
	}
	return nil
}
func randomString(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
