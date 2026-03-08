package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AffiliateRepository interface {
	Create(affiliate *models.Affiliate) error
	FindAll(offset, limit int) ([]models.Affiliate, int64, error)
	FindByID(id uuid.UUID) (*models.Affiliate, error)
	FindByCode(code string) (*models.Affiliate, error)
	Update(affiliate *models.Affiliate) error
	Delete(id uuid.UUID) error
}

type affiliateRepository struct {
	db *gorm.DB
}

func NewAffiliateRepository(db *gorm.DB) AffiliateRepository {
	return &affiliateRepository{db: db}
}

func (r *affiliateRepository) Create(affiliate *models.Affiliate) error {
	return r.db.Create(affiliate).Error
}

func (r *affiliateRepository) FindAll(offset, limit int) ([]models.Affiliate, int64, error) {
	var affiliates []models.Affiliate
	var total int64
	r.db.Model(&models.Affiliate{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&affiliates).Error
	return affiliates, total, err
}

func (r *affiliateRepository) FindByID(id uuid.UUID) (*models.Affiliate, error) {
	var affiliate models.Affiliate
	err := r.db.First(&affiliate, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &affiliate, nil
}

func (r *affiliateRepository) FindByCode(code string) (*models.Affiliate, error) {
	var affiliate models.Affiliate
	err := r.db.First(&affiliate, "affiliate_code = ? AND is_active = ?", code, true).Error
	if err != nil {
		return nil, err
	}
	return &affiliate, nil
}

func (r *affiliateRepository) Update(affiliate *models.Affiliate) error {
	return r.db.Save(affiliate).Error
}

func (r *affiliateRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Affiliate{}, "id = ?", id).Error
}

// AffiliateCommissionRepository
type AffiliateCommissionRepository interface {
	CreateWithTx(tx *gorm.DB, ac *models.AffiliateCommission) error
	FindByAffiliateID(affiliateID uuid.UUID, offset, limit int) ([]models.AffiliateCommission, int64, error)
	FindByTransactionID(txnID uuid.UUID) (*models.AffiliateCommission, error)
	UpdateWithTx(tx *gorm.DB, ac *models.AffiliateCommission) error
}

type affiliateCommissionRepository struct {
	db *gorm.DB
}

func NewAffiliateCommissionRepository(db *gorm.DB) AffiliateCommissionRepository {
	return &affiliateCommissionRepository{db: db}
}

func (r *affiliateCommissionRepository) CreateWithTx(tx *gorm.DB, ac *models.AffiliateCommission) error {
	return tx.Create(ac).Error
}

func (r *affiliateCommissionRepository) FindByAffiliateID(affiliateID uuid.UUID, offset, limit int) ([]models.AffiliateCommission, int64, error) {
	var commissions []models.AffiliateCommission
	var total int64
	r.db.Model(&models.AffiliateCommission{}).Where("affiliate_id = ?", affiliateID).Count(&total)
	err := r.db.Preload("Transaction").Where("affiliate_id = ?", affiliateID).
		Offset(offset).Limit(limit).Find(&commissions).Error
	return commissions, total, err
}

func (r *affiliateCommissionRepository) FindByTransactionID(txnID uuid.UUID) (*models.AffiliateCommission, error) {
	var ac models.AffiliateCommission
	err := r.db.Where("transaction_id = ?", txnID).First(&ac).Error
	if err != nil {
		return nil, err
	}
	return &ac, nil
}

func (r *affiliateCommissionRepository) UpdateWithTx(tx *gorm.DB, ac *models.AffiliateCommission) error {
	return tx.Save(ac).Error
}
