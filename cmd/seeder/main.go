package main

import (
	"fmt"
	"log"
	"time"

	"hairhaus-pos-be/config"
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("🌱 Starting database seeder...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.InitDatabase(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Helper for safe hashing
	hashPassword := func(pw string) string {
		hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
		return string(hash)
	}

	// Helper for pointer ints
	intPtr := func(i int64) *int64 { return &i }

	// --- 1. INDEPENDENT TABLES ---
	fmt.Println("Seeding Branches...")
	branchPusat := models.Branch{Code: "BKS", Name: "HAIRHAUS Bekasi", Address: "Jl. Sudirman No 1", Phone: "08123456789", OutletID: "cm1tntou2000611xf7cpkhb7r", IsActive: true}
	branchCabang := models.Branch{Code: "KNG", Name: "HAIRHAUS Kuningan", Address: "Jl. Thamrin No 2", Phone: "08198765432", OutletID: "cm1tntuli000711xfo7cbgz3k", IsActive: true}
	db.FirstOrCreate(&branchPusat, models.Branch{Code: "BKS"})
	db.FirstOrCreate(&branchCabang, models.Branch{Code: "KNG"})

	fmt.Println("Seeding Product Categories...")
	catHaircut := models.ProductCategory{Code: "HC", Name: "Haircut"}
	catTreatment := models.ProductCategory{Code: "TRT", Name: "Treatment"}
	catRetail := models.ProductCategory{Code: "RTL", Name: "Retail"}
	db.FirstOrCreate(&catHaircut, models.ProductCategory{Code: "HC"})
	db.FirstOrCreate(&catTreatment, models.ProductCategory{Code: "TRT"})
	db.FirstOrCreate(&catRetail, models.ProductCategory{Code: "RTL"})

	fmt.Println("Seeding Stylists...")
	stylistBondan := models.Stylist{Name: "Djarot", IsActive: true}
	stylistDjarot := models.Stylist{Name: "Fika", IsActive: true}
	db.FirstOrCreate(&stylistBondan, models.Stylist{Name: "Djarot"})
	db.FirstOrCreate(&stylistDjarot, models.Stylist{Name: "Fika"})

	fmt.Println("Seeding Customers...")
	customer1 := models.Customer{Name: "Budi", Phone: "081111111", LoyaltyExternalID: "LYL001"}
	customer2 := models.Customer{Name: "Andi", Phone: "082222222", LoyaltyExternalID: "LYL002"}
	db.FirstOrCreate(&customer1, models.Customer{Phone: "081111111"})
	db.FirstOrCreate(&customer2, models.Customer{Phone: "082222222"})

	fmt.Println("Seeding Affiliates...")
	affiliate1 := models.Affiliate{
		LoyaltyMemberID:      "LYL003",
		AffiliateCode:        "AFF-PROMO",
		Name:                 "Promo Partner",
		CommissionType:       "PERCENTAGE",
		CommissionPercentage: 0.1,
		IsActive:             true,
	}
	db.FirstOrCreate(&affiliate1, models.Affiliate{AffiliateCode: "AFF-PROMO"})

	fmt.Println("Seeding Expense Categories...")
	expCatListrik := models.ExpenseCategory{Code: "ELEC", Name: "Electricity"}
	expCatAir := models.ExpenseCategory{Code: "WATER", Name: "Water"}
	db.FirstOrCreate(&expCatListrik, models.ExpenseCategory{Code: "ELEC"})
	db.FirstOrCreate(&expCatAir, models.ExpenseCategory{Code: "WATER"})

	// --- 2. FIRST-LEVEL DEPENDENT ---
	fmt.Println("Seeding Users...")
	users := []models.User{
		{Name: "Admin", PhoneNumber: "082210001000", EmployeeID: "EMP001", Pin: hashPassword("0808"), Role: models.RoleAdmin, BranchID: branchPusat.ID},
		{Name: "Manager Pusat", PhoneNumber: "082299990359", EmployeeID: "EMP002", Pin: hashPassword("0808"), Role: models.RoleManager, BranchID: branchPusat.ID},
		{Name: "Cashier Pusat", PhoneNumber: "082210007020", EmployeeID: "EMP003", Pin: hashPassword("0808"), Role: models.RoleCashier, BranchID: branchPusat.ID},
	}
	for _, u := range users {
		db.FirstOrCreate(&u, models.User{PhoneNumber: u.PhoneNumber})
	}

	fmt.Println("Seeding Products...")
	prodHaircut := models.Product{Name: "Premium Haircut", ProductType: "SERVICE", CategoryID: catHaircut.ID, BasePrice: 100000, CostPrice: 0, IsActive: true}
	prodCreambath := models.Product{Name: "Creambath", ProductType: "SERVICE", CategoryID: catTreatment.ID, BasePrice: 150000, CostPrice: 50000, IsActive: true}
	prodPomade := models.Product{Name: "Hairhaus Pomade", ProductType: "RETAIL", CategoryID: catRetail.ID, BasePrice: 200000, CostPrice: 100000, IsActive: true}
	db.FirstOrCreate(&prodHaircut, models.Product{Name: "Premium Haircut"})
	db.FirstOrCreate(&prodCreambath, models.Product{Name: "Creambath"})
	db.FirstOrCreate(&prodPomade, models.Product{Name: "Hairhaus Pomade"})

	fmt.Println("Seeding Branch Stylists...")
	bs1 := models.BranchStylist{BranchID: branchPusat.ID, StylistID: stylistBondan.ID, HaircutPriceOverride: nil}
	bs2 := models.BranchStylist{BranchID: branchPusat.ID, StylistID: stylistDjarot.ID, HaircutPriceOverride: intPtr(120000)} // Djarot is premium
	db.FirstOrCreate(&bs1, models.BranchStylist{BranchID: branchPusat.ID, StylistID: stylistBondan.ID})
	db.FirstOrCreate(&bs2, models.BranchStylist{BranchID: branchPusat.ID, StylistID: stylistDjarot.ID})

	// --- 3. SECOND-LEVEL DEPENDENT ---
	fmt.Println("Seeding Branch Products...")
	bp1 := models.BranchProduct{BranchID: branchPusat.ID, ProductID: prodHaircut.ID, PriceOverride: nil, Stock: 0}
	bp2 := models.BranchProduct{BranchID: branchPusat.ID, ProductID: prodCreambath.ID, PriceOverride: nil, Stock: 0}
	bp3 := models.BranchProduct{BranchID: branchPusat.ID, ProductID: prodPomade.ID, PriceOverride: intPtr(210000), Stock: 50}
	db.FirstOrCreate(&bp1, models.BranchProduct{BranchID: branchPusat.ID, ProductID: prodHaircut.ID})
	db.FirstOrCreate(&bp2, models.BranchProduct{BranchID: branchPusat.ID, ProductID: prodCreambath.ID})
	db.FirstOrCreate(&bp3, models.BranchProduct{BranchID: branchPusat.ID, ProductID: prodPomade.ID})

	// --- 4. OPERATIONAL / TRANSACTIONS ---
	fmt.Println("Seeding Open Cash Drawer...")
	drawer := models.CashDrawer{
		BranchID:      branchPusat.ID,
		OpenedAt:      time.Now(),
		OpeningAmount: 500000,
		ExpectedCash:  500000,
		Status:        "OPEN",
	}
	// Only seed if no OPEN drawer exists
	var count int64
	db.Model(&models.CashDrawer{}).Where("branch_id = ? AND status = 'OPEN'", branchPusat.ID).Count(&count)
	if count == 0 {
		db.Create(&drawer)
	}

	fmt.Println("Seeding Stock Movement...")
	// We gave the pomade a stock of 50. Let's record the movement for completeness.
	stockMov := models.StockMovement{
		BranchID:    branchPusat.ID,
		ProductID:   prodPomade.ID,
		Change:      50,
		Type:        "RESTOCK",
		ReferenceID: uuid.New(), // dummy ID
		Note:        "Initial DB Seeding",
	}
	var smCount int64
	db.Model(&models.StockMovement{}).Where("branch_id = ? AND product_id = ? AND type = 'RESTOCK'", branchPusat.ID, prodPomade.ID).Count(&smCount)
	if smCount == 0 {
		db.Create(&stockMov)
	}

	fmt.Println("Seeding Demo Expense...")
	exp1 := models.Expense{
		BranchID:    branchPusat.ID,
		CategoryID:  expCatListrik.ID,
		Description: "Listrik Bulan Lalu",
		Amount:      500000,
		ExpenseDate: time.Now(),
	}
	db.FirstOrCreate(&exp1, models.Expense{Description: "Listrik Bulan Lalu"})

	log.Println("✅ Database seeded successfully!")
}
