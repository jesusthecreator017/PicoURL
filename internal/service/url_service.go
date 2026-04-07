package service

import (
	"context"

	"github.com/jesusthecreator017/PicoURL/internal/shortcode"
	"github.com/jesusthecreator017/PicoURL/internal/store"
	"github.com/jesusthecreator017/PicoURL/internal/utils"
)

type URLService interface {
	Shorten(ctx context.Context, originalURL string) (string, error)
	Resolve(ctx context.Context, shortURL string) (string, error)
	GetStats(ctx context.Context, shortURL string) (int, error)
	Delete(ctx context.Context, shortURL string) error
}

type urlService struct {
	store      store.Store
	codeLength int
}

func NewURLService(store store.Store) *urlService {
	return &urlService{
		store:      store,
		codeLength: 7,
	}
}

func (s *urlService) Shorten(ctx context.Context, originalURL string) (string, error) {
	// Validate URL format and reachability
	if err := utils.ValidateURLFormat(originalURL); err != nil {
		return "", ErrInvalidURL
	}

	if err := utils.ValidateURLReachable(originalURL); err != nil {
		return "", ErrUnreachableURL
	}

	// Generate deterministic short code
	shortURL := shortcode.GenerateShortCode(originalURL, s.codeLength)

	// Check if this short code already exists
	existingURL, err := s.store.GetOriginalURL(ctx, shortURL)
	if err != nil {
		// Not found — save the new URL
		if err := s.store.SaveURL(ctx, shortURL, originalURL); err != nil {
			return "", err
		}
		return shortURL, nil
	}

	// Already exists and maps to the same URL (idempotent)
	if existingURL == originalURL {
		return shortURL, nil
	}

	// Collision — different URL hashed to the same short code
	return "", ErrCollision
}

func (s *urlService) Resolve(ctx context.Context, shortURL string) (string, error) {
	originalURL, err := s.store.GetOriginalURL(ctx, shortURL)
	if err != nil {
		return "", ErrNotFound
	}

	// Fire-and-forget — don't fail the redirect if count increment fails
	_ = s.store.IncrementCount(ctx, shortURL)

	return originalURL, nil
}

func (s *urlService) GetStats(ctx context.Context, shortURL string) (int, error) {
	return s.store.GetCount(ctx, shortURL)
}

func (s *urlService) Delete(ctx context.Context, shortURL string) error {
	return s.store.DeleteURL(ctx, shortURL)
}
