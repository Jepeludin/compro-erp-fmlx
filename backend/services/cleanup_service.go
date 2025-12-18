package services

import (
	"context"
	"log"
	"sync"
	"time"

	"ganttpro-backend/repository"
)

// CleanupService handles periodic cleanup tasks
type CleanupService struct {
	tokenBlacklistRepo *repository.TokenBlacklistRepository
	interval           time.Duration
	stopCh             chan struct{}
	wg                 sync.WaitGroup
	running            bool
	mu                 sync.Mutex
}

// CleanupConfig holds configuration for cleanup service
type CleanupConfig struct {
	Interval time.Duration
}

// DefaultCleanupConfig returns default cleanup configuration (1 hour interval)
func DefaultCleanupConfig() CleanupConfig {
	return CleanupConfig{
		Interval: 1 * time.Hour,
	}
}

// NewCleanupService creates a new cleanup service
func NewCleanupService(tokenBlacklistRepo *repository.TokenBlacklistRepository, config CleanupConfig) *CleanupService {
	if config.Interval <= 0 {
		config.Interval = 1 * time.Hour
	}

	return &CleanupService{
		tokenBlacklistRepo: tokenBlacklistRepo,
		interval:           config.Interval,
		stopCh:             make(chan struct{}),
	}
}

// Start begins the background cleanup goroutine
func (s *CleanupService) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		log.Println("[CleanupService] Already running")
		return
	}

	s.running = true
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		log.Printf("[CleanupService] Started (interval: %v)", s.interval)

		// Run cleanup immediately on startup
		s.runCleanup()

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.runCleanup()
			case <-s.stopCh:
				log.Println("[CleanupService] Stopped")
				return
			}
		}
	}()
}

// Stop gracefully stops the cleanup service
func (s *CleanupService) Stop(ctx context.Context) error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = false
	s.mu.Unlock()

	close(s.stopCh)

	// Wait for goroutine to finish or context timeout
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// runCleanup performs all cleanup tasks
func (s *CleanupService) runCleanup() {
	startTime := time.Now()

	// Clean expired blacklisted tokens
	count, err := s.tokenBlacklistRepo.CleanExpiredTokensWithCount()
	if err != nil {
		log.Printf("[CleanupService] Error cleaning expired tokens: %v", err)
		return
	}

	duration := time.Since(startTime)

	if count > 0 {
		log.Printf("[CleanupService] Cleaned %d expired tokens in %v", count, duration)
	}
}

// RunNow triggers an immediate cleanup (useful for testing)
func (s *CleanupService) RunNow() {
	s.runCleanup()
}

// GetStats returns cleanup service statistics
func (s *CleanupService) GetStats() CleanupStats {
	s.mu.Lock()
	defer s.mu.Unlock()

	totalCount, _ := s.tokenBlacklistRepo.GetBlacklistCount()
	expiredCount, _ := s.tokenBlacklistRepo.GetExpiredCount()

	return CleanupStats{
		IsRunning:     s.running,
		Interval:      s.interval,
		TotalTokens:   totalCount,
		ExpiredTokens: expiredCount,
	}
}

// CleanupStats holds statistics about the cleanup service
type CleanupStats struct {
	IsRunning     bool          `json:"is_running"`
	Interval      time.Duration `json:"interval"`
	TotalTokens   int64         `json:"total_tokens"`
	ExpiredTokens int64         `json:"expired_tokens"`
}
