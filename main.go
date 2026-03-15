package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		Customer:        handlers.NewCustomerHandler(customerService, loyaltyClient),
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

	// Setup server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 HAIRHAUS ERP server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
