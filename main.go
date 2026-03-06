package main

import (
	"fmt"
	"log"

	"hairhaus-pos-be/clients"
	"hairhaus-pos-be/config"
	"hairhaus-pos-be/docs"
	"hairhaus-pos-be/handlers"
	"hairhaus-pos-be/repositories"
	"hairhaus-pos-be/routes"
	"hairhaus-pos-be/services"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Run database migrations via golang-migrate (SQL files)
	if err := config.RunMigrations(cfg.DB, "migrations"); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Connect database (GORM — used only for queries, not migrations)
	db, err := config.InitDatabase(cfg.DB)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// Connect Redis
	rdb, err := config.InitRedis(cfg.Redis)
	if err != nil {
		log.Println("⚠️  Redis connection failed (non-fatal):", err)
	} else {
		_ = rdb // Redis client available for future caching use
		log.Println("✅ Redis ready for use")
	}

	// Initialize repositories
	branchRepo := repositories.NewBranchRepository(db)
	userRepo := repositories.NewUserRepository(db)
	productCategoryRepo := repositories.NewProductCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	stylistRepo := repositories.NewStylistRepository(db)
	branchProductRepo := repositories.NewBranchProductRepository(db)
	branchStylistRepo := repositories.NewBranchStylistRepository(db)
	customerRepo := repositories.NewCustomerRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	cashDrawerRepo := repositories.NewCashDrawerRepository(db)
	stockMovementRepo := repositories.NewStockMovementRepository(db)
	affiliateRepo := repositories.NewAffiliateRepository(db)
	affiliateCommRepo := repositories.NewAffiliateCommissionRepository(db)
	salaryRepo := repositories.NewSalaryRepository(db)
	expenseCategoryRepo := repositories.NewExpenseCategoryRepository(db)
	expenseRepo := repositories.NewExpenseRepository(db)
	auditLogRepo := repositories.NewAuditLogRepository(db)

	// Initialize secondary clients
	loyaltyClient := clients.NewLoyaltyClient(cfg.LoyaltyService.APIURL, cfg.LoyaltyService.APIKey)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	branchService := services.NewBranchService(branchRepo)
	userService := services.NewUserService(userRepo)
	productCategoryService := services.NewProductCategoryService(productCategoryRepo)
	productService := services.NewProductService(productRepo)
	stylistService := services.NewStylistService(stylistRepo)
	branchProductService := services.NewBranchProductService(branchProductRepo)
	branchStylistService := services.NewBranchStylistService(branchStylistRepo)
	customerService := services.NewCustomerService(customerRepo)
	cashDrawerService := services.NewCashDrawerService(cashDrawerRepo, transactionRepo, auditLogRepo)
	transactionService := services.NewTransactionService(
		transactionRepo, branchRepo, branchProductRepo, branchStylistRepo,
		productRepo, stylistRepo, cashDrawerRepo, affiliateRepo,
		affiliateCommRepo, stockMovementRepo, auditLogRepo, loyaltyClient,
	)
	inventoryService := services.NewInventoryService(branchProductRepo, stockMovementRepo)
	affiliateService := services.NewAffiliateService(affiliateRepo, affiliateCommRepo)
	salaryService := services.NewSalaryService(salaryRepo, transactionRepo, auditLogRepo)
	expenseCategoryService := services.NewExpenseCategoryService(expenseCategoryRepo)
	expenseService := services.NewExpenseService(expenseRepo)
	reportService := services.NewReportService(transactionRepo, expenseRepo)

	// Initialize handlers
	h := &routes.Handlers{
		Auth:            handlers.NewAuthHandler(authService),
		Branch:          handlers.NewBranchHandler(branchService),
		User:            handlers.NewUserHandler(userService),
		ProductCategory: handlers.NewProductCategoryHandler(productCategoryService),
		Product:         handlers.NewProductHandler(productService),
		Stylist:         handlers.NewStylistHandler(stylistService),
		BranchProduct:   handlers.NewBranchProductHandler(branchProductService),
		BranchStylist:   handlers.NewBranchStylistHandler(branchStylistService),
		Customer:        handlers.NewCustomerHandler(customerService),
		CashDrawer:      handlers.NewCashDrawerHandler(cashDrawerService),
		Transaction:     handlers.NewTransactionHandler(transactionService),
		Inventory:       handlers.NewInventoryHandler(inventoryService),
		Affiliate:       handlers.NewAffiliateHandler(affiliateService),
		Salary:          handlers.NewSalaryHandler(salaryService),
		ExpenseCategory: handlers.NewExpenseCategoryHandler(expenseCategoryService),
		Expense:         handlers.NewExpenseHandler(expenseService),
		Report:          handlers.NewReportHandler(reportService),
	}

	// Setup router
	router := routes.SetupRouter(h, cfg.JWT.Secret)

	// Setup Swagger UI at /swagger
	docs.SetupSwagger(router)
	log.Println("📖 Swagger UI available at /swagger")

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("🚀 HAIRHAUS ERP server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
