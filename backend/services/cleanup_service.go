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
	Interval time.Duration // How often to run cleanup (default: 1 hour)
}

// DefaultCleanupConfig returns default cleanup configuration
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
		log.Println("Cleanup service is already running")
		return
	}

	s.running = true
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		log.Printf("Cleanup service started (interval: %v)", s.interval)

		// Run cleanup immediately on startup
		s.runCleanup()

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.runCleanup()
			case <-s.stopCh:
				log.Println("Cleanup service stopped")
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
	log.Println("Running scheduled cleanup tasks...")

	// Clean expired blacklisted tokens
	if err := s.cleanExpiredTokens(); err != nil {
		log.Printf("Error cleaning expired tokens: %v", err)
	}

	log.Println("Cleanup tasks completed")
}

// cleanExpiredTokens removes expired tokens from the blacklist
func (s *CleanupService) cleanExpiredTokens() error {
	startTime := time.Now()

	count, err := s.tokenBlacklistRepo.CleanExpiredTokensWithCount()
	if err != nil {
		return err
	}

	if count > 0 {
		log.Printf("Cleaned %d expired tokens in %v", count, time.Since(startTime))
	} else {
		log.Printf("No expired tokens to clean (took %v)", time.Since(startTime))
	}

	return nil
}

// RunNow triggers an immediate cleanup (useful for testing or manual triggers)
func (s *CleanupService) RunNow() {
	s.runCleanup()
}

