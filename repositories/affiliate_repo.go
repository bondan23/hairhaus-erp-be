package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AffiliateRepository struct {
	db *gorm.DB
}

func NewAffiliateRepository(db *gorm.DB) *AffiliateRepository {
	return &AffiliateRepository{db: db}
}

func (r *AffiliateRepository) Create(affiliate *models.Affiliate) error {
	return r.db.Create(affiliate).Error
}

func (r *AffiliateRepository) FindAll(offset, limit int) ([]models.Affiliate, int64, error) {
	var affiliates []models.Affiliate
	var total int64
	r.db.Model(&models.Affiliate{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&affiliates).Error
	return affiliates, total, err
}

func (r *AffiliateRepository) FindByID(id uuid.UUID) (*models.Affiliate, error) {
	var affiliate models.Affiliate
	err := r.db.First(&affiliate, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &affiliate, nil
}

func (r *AffiliateRepository) FindByCode(code string) (*models.Affiliate, error) {
	var affiliate models.Affiliate
	err := r.db.First(&affiliate, "affiliate_code = ? AND is_active = ?", code, true).Error
	if err != nil {
		return nil, err
	}
	return &affiliate, nil
}

func (r *AffiliateRepository) Update(affiliate *models.Affiliate) error {
	return r.db.Save(affiliate).Error
}

func (r *AffiliateRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Affiliate{}, "id = ?", id).Error
}

// AffiliateCommissionRepository
type AffiliateCommissionRepository struct {
	db *gorm.DB
}

func NewAffiliateCommissionRepository(db *gorm.DB) *AffiliateCommissionRepository {
	return &AffiliateCommissionRepository{db: db}
}

func (r *AffiliateCommissionRepository) CreateWithTx(tx *gorm.DB, ac *models.AffiliateCommission) error {
	return tx.Create(ac).Error
}

func (r *AffiliateCommissionRepository) FindByAffiliateID(affiliateID uuid.UUID, offset, limit int) ([]models.AffiliateCommission, int64, error) {
	var commissions []models.AffiliateCommission
	var total int64
	r.db.Model(&models.AffiliateCommission{}).Where("affiliate_id = ?", affiliateID).Count(&total)
	err := r.db.Preload("Transaction").Where("affiliate_id = ?", affiliateID).
		Offset(offset).Limit(limit).Find(&commissions).Error
	return commissions, total, err
}
