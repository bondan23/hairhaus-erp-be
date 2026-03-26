package routes

import (
	"hairhaus-pos-be/handlers"
	"hairhaus-pos-be/middleware"
	"hairhaus-pos-be/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth            *handlers.AuthHandler
	Branch          *handlers.BranchHandler
	User            *handlers.UserHandler
	ProductCategory *handlers.ProductCategoryHandler
	Product         *handlers.ProductHandler
	Stylist         *handlers.StylistHandler
	BranchProduct   *handlers.BranchProductHandler
	BranchStylist   *handlers.BranchStylistHandler
	Customer        *handlers.CustomerHandler
	CashDrawer      *handlers.CashDrawerHandler
	Transaction     *handlers.TransactionHandler
	Inventory       *handlers.InventoryHandler
	Affiliate       *handlers.AffiliateHandler
	Salary          *handlers.SalaryHandler
	ExpenseCategory *handlers.ExpenseCategoryHandler
	Expense         *handlers.ExpenseHandler
	Report          *handlers.ReportHandler
}

func SetupRouter(h *Handlers, jwtSecret string) *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "hairhaus-erp"})
	})

	api := r.Group("/api/v1")

	// === Public Routes ===
	auth := api.Group("/auth")
	{
		auth.POST("/login", h.Auth.Login)
	}

	// === Protected Routes ===
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		// --- Branches ---
		branches := protected.Group("/branches")
		{
			branches.GET("", h.Branch.GetAll)
			branches.GET("/:id", h.Branch.GetByID)
			branches.POST("", middleware.RequireRole(models.RoleAdmin), h.Branch.Create)
			branches.PUT("/:id", middleware.RequireRole(models.RoleAdmin), h.Branch.Update)
			branches.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), h.Branch.Delete)

			// Branch Products
			branches.GET("/:id/products", h.BranchProduct.GetByBranch)
			branches.GET("/:id/stylists", h.BranchStylist.GetByBranch)

			// Cash Drawers by branch
			branches.GET("/:id/drawers", h.CashDrawer.GetByBranch)
			branches.GET("/:id/drawers/open", h.CashDrawer.GetOpen)

			// Transactions by branch
			branches.GET("/:id/transactions", h.Transaction.GetByBranch)

			// Expenses by branch
			branches.GET("/:id/expenses", h.Expense.GetByBranch)

			// Salaries by branch
			branches.GET("/:id/salaries", h.Salary.GetByBranch)
		}

		// --- Users ---
		users := protected.Group("/users")
		users.Use(middleware.RequireRole(models.RoleAdmin))
		{
			users.POST("", h.User.Create)
			users.GET("", h.User.GetAll)
			users.GET("/:id", h.User.GetByID)
			users.PUT("/:id", h.User.Update)
			users.DELETE("/:id", h.User.Delete)
		}

		// --- Product Categories ---
		productCategories := protected.Group("/product-categories")
		{
			productCategories.GET("", h.ProductCategory.GetAll)
			productCategories.GET("/:id", h.ProductCategory.GetByID)
			productCategories.POST("", middleware.RequireRole(models.RoleAdmin), h.ProductCategory.Create)
			productCategories.PUT("/:id", middleware.RequireRole(models.RoleAdmin), h.ProductCategory.Update)
			productCategories.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), h.ProductCategory.Delete)
		}

		// --- Products ---
		products := protected.Group("/products")
		{
			products.GET("", h.Product.GetAll)
			products.GET("/:id", h.Product.GetByID)
			products.POST("", middleware.RequireRole(models.RoleAdmin), h.Product.Create)
			products.PUT("/:id", middleware.RequireRole(models.RoleAdmin), h.Product.Update)
			products.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), h.Product.Delete)
		}

		// --- Stylists ---
		stylists := protected.Group("/stylists")
		{
			stylists.GET("", h.Stylist.GetAll)
			stylists.GET("/:id", h.Stylist.GetByID)
			stylists.POST("", middleware.RequireRole(models.RoleAdmin), h.Stylist.Create)
			stylists.PUT("/:id", middleware.RequireRole(models.RoleAdmin), h.Stylist.Update)
			stylists.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), h.Stylist.Delete)
		}

		// --- Branch Products ---
		branchProducts := protected.Group("/branch-products")
		{
			branchProducts.GET("/:id", h.BranchProduct.GetByID)
			branchProducts.POST("", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.BranchProduct.Create)
			branchProducts.PUT("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.BranchProduct.Update)
			branchProducts.DELETE("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.BranchProduct.Delete)
		}

		// --- Branch Stylists ---
		branchStylists := protected.Group("/branch-stylists")
		{
			branchStylists.GET("", h.BranchStylist.GetAll)
			branchStylists.GET("/:id", h.BranchStylist.GetByID)
			branchStylists.POST("", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.BranchStylist.Create)
			branchStylists.PUT("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.BranchStylist.Update)
			branchStylists.DELETE("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.BranchStylist.Delete)
		}

		// --- Customers ---
		customers := protected.Group("/customers")
		{
			customers.GET("", h.Customer.GetAll)
			customers.GET("/deleted", h.Customer.GetAllSoftDelete)
			customers.GET("/:id", h.Customer.GetByID)
			customers.POST("", h.Customer.Create)
			customers.PUT("/:id", h.Customer.Update)
			customers.DELETE("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.Customer.Delete)

			// Identification & Loyalty
			customers.POST("/identify", h.Customer.Identify)
			customers.POST("/loyalty/register", h.Customer.RegisterLoyalty)
			customers.POST("/loyalty/otp", h.Customer.RequestLoyaltyOTP)
			customers.POST("/loyalty/verify", h.Customer.VerifyLoyaltyOTP)
		}

		// --- Cash Drawers ---
		drawers := protected.Group("/drawers")
		{
			drawers.GET("/:id", h.CashDrawer.GetByID)
			drawers.POST("/open", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.CashDrawer.Open)
			drawers.POST("/:id/close", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.CashDrawer.Close)
		}

		// --- Transactions ---
		transactions := protected.Group("/transactions")
		{
			transactions.POST("/save", middleware.RequireRole(models.RoleAdmin, models.RoleManager, models.RoleCashier), h.Transaction.Save)
			transactions.PUT("/draft/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager, models.RoleCashier), h.Transaction.EditDraft)
			transactions.POST("/checkout", middleware.RequireRole(models.RoleAdmin, models.RoleManager, models.RoleCashier), h.Transaction.Checkout)
			transactions.POST("/:id/checkout", middleware.RequireRole(models.RoleAdmin, models.RoleManager, models.RoleCashier), h.Transaction.CheckoutDraft)
			transactions.GET("/:id", h.Transaction.GetByID)
			transactions.PUT("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.Transaction.Edit)
			transactions.POST("/:id/void", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.Transaction.Void)
		}

		// --- Inventory ---
		inventory := protected.Group("/inventory")
		inventory.Use(middleware.RequireRole(models.RoleAdmin, models.RoleManager))
		{
			inventory.POST("/restock", h.Inventory.Restock)
			inventory.POST("/adjust", h.Inventory.Adjust)
			inventory.GET("/movements", h.Inventory.GetMovements)
		}

		// --- Affiliates ---
		affiliates := protected.Group("/affiliates")
		affiliates.Use(middleware.RequireRole(models.RoleAdmin))
		{
			affiliates.POST("", h.Affiliate.Create)
			affiliates.GET("", h.Affiliate.GetAll)
			affiliates.GET("/:id", h.Affiliate.GetByID)
			affiliates.PUT("/:id", h.Affiliate.Update)
			affiliates.DELETE("/:id", h.Affiliate.Delete)
			affiliates.GET("/:id/commissions", h.Affiliate.GetCommissions)
		}

		// --- Salary ---
		salary := protected.Group("/salaries")
		salary.Use(middleware.RequireRole(models.RoleAdmin))
		{
			salary.POST("/generate", h.Salary.Generate)
			salary.GET("/:id", h.Salary.GetByID)
			salary.POST("/:id/paid", h.Salary.MarkPaid)
		}

		// --- Expense Categories ---
		expenseCategories := protected.Group("/expense-categories")
		{
			expenseCategories.GET("", h.ExpenseCategory.GetAll)
			expenseCategories.GET("/:id", h.ExpenseCategory.GetByID)
			expenseCategories.POST("", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.ExpenseCategory.Create)
			expenseCategories.PUT("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.ExpenseCategory.Update)
			expenseCategories.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), h.ExpenseCategory.Delete)
		}

		// --- Expenses ---
		expenses := protected.Group("/expenses")
		{
			expenses.GET("/:id", h.Expense.GetByID)
			expenses.POST("", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.Expense.Create)
			expenses.PUT("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.Expense.Update)
			expenses.DELETE("/:id", middleware.RequireRole(models.RoleAdmin, models.RoleManager), h.Expense.Delete)
		}

		// --- Reports ---
		reports := protected.Group("/reports")
		reports.Use(middleware.RequireRole(models.RoleAdmin, models.RoleManager))
		{
			reports.GET("/financial", h.Report.GetFinancialReport)
			reports.GET("/income-by-type", h.Report.GetRevenueByIncomeType)
		}
	}

	return r
}
