package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BranchStylistRepository struct {
	db *gorm.DB
}

func NewBranchStylistRepository(db *gorm.DB) *BranchStylistRepository {
	return &BranchStylistRepository{db: db}
}

func (r *BranchStylistRepository) Create(bs *models.BranchStylist) error {
	return r.db.Create(bs).Error
}

func (r *BranchStylistRepository) FindByBranchID(branchID uuid.UUID, offset, limit int) ([]models.BranchStylist, int64, error) {
	var bss []models.BranchStylist
	var total int64
	r.db.Model(&models.BranchStylist{}).Where("branch_id = ?", branchID).Count(&total)
	err := r.db.Preload("Stylist").Where("branch_id = ?", branchID).
		Offset(offset).Limit(limit).Find(&bss).Error
	return bss, total, err
}

func (r *BranchStylistRepository) FindByID(id uuid.UUID) (*models.BranchStylist, error) {
	var bs models.BranchStylist
	err := r.db.Preload("Stylist").First(&bs, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &bs, nil
}

func (r *BranchStylistRepository) FindByBranchAndStylist(branchID, stylistID uuid.UUID) (*models.BranchStylist, error) {
	var bs models.BranchStylist
	err := r.db.Preload("Stylist").
		First(&bs, "branch_id = ? AND stylist_id = ?", branchID, stylistID).Error
	if err != nil {
		return nil, err
	}
	return &bs, nil
}

func (r *BranchStylistRepository) Update(bs *models.BranchStylist) error {
	return r.db.Save(bs).Error
}

func (r *BranchStylistRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.BranchStylist{}, "id = ?", id).Error
}
