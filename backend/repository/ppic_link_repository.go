package repository

import (
	"ganttpro-backend/models"

	"gorm.io/gorm"
)

type PPICLinkRepository struct {
	db *gorm.DB
}

func NewPPICLinkRepository(db *gorm.DB) *PPICLinkRepository {
	return &PPICLinkRepository{db: db}
}

// Create creates a new PPIC link
func (r *PPICLinkRepository) Create(req *models.CreatePPICLinkRequest) (*models.PPICLink, error) {
	link := &models.PPICLink{
		SourceScheduleID: req.SourceScheduleID,
		TargetScheduleID: req.TargetScheduleID,
		LinkType:         req.LinkType,
	}

	if err := r.db.Create(link).Error; err != nil {
		return nil, err
	}

	return link, nil
}

// GetAll returns all PPIC links
func (r *PPICLinkRepository) GetAll() ([]models.PPICLink, error) {
	var links []models.PPICLink
	if err := r.db.Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

// Delete deletes a PPIC link by ID
func (r *PPICLinkRepository) Delete(id int64) error {
	return r.db.Delete(&models.PPICLink{}, id).Error
}

// GetBySourceScheduleID returns all links where the given schedule is the source
func (r *PPICLinkRepository) GetBySourceScheduleID(scheduleID int64) ([]models.PPICLink, error) {
	var links []models.PPICLink
	if err := r.db.Where("source_schedule_id = ?", scheduleID).Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

// GetByTargetScheduleID returns all links where the given schedule is the target
func (r *PPICLinkRepository) GetByTargetScheduleID(scheduleID int64) ([]models.PPICLink, error) {
	var links []models.PPICLink
	if err := r.db.Where("target_schedule_id = ?", scheduleID).Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}
