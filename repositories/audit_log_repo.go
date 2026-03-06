package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLogRepository interface {
	Create(log *models.AuditLog) error
	CreateWithTx(tx *gorm.DB, log *models.AuditLog) error
	FindByEntity(entityType string, entityID uuid.UUID) ([]models.AuditLog, error)
	FindAll(offset, limit int) ([]models.AuditLog, int64, error)
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *auditLogRepository) CreateWithTx(tx *gorm.DB, log *models.AuditLog) error {
	return tx.Create(log).Error
}

func (r *auditLogRepository) FindByEntity(entityType string, entityID uuid.UUID) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) FindAll(offset, limit int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64
	r.db.Model(&models.AuditLog{}).Count(&total)
	err := r.db.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error
	return logs, total, err
}
