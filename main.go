package main

import (
	"fmt"
	"log"

	"github.com/fleimkeipa/kubernetes-api/config"
	"github.com/fleimkeipa/kubernetes-api/controller"
	_ "github.com/fleimkeipa/kubernetes-api/docs" // which is the generated folder after swag init
	"github.com/fleimkeipa/kubernetes-api/pkg"
	"github.com/fleimkeipa/kubernetes-api/repositories"
	"github.com/fleimkeipa/kubernetes-api/uc"
	"github.com/fleimkeipa/kubernetes-api/util"

	"github.com/go-pg/pg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

func main() {
	// Load environment configuration
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Start the application
	serveApplication()
}

func serveApplication() {
	// Create a new Echo instance
	var e = echo.New()

	// Configure Echo settings
	configureEcho(e)

	// Configure CORS middleware
	configureCORS(e)

	// Configure the logger
	var sugar = configureLogger(e)
	defer sugar.Sync() // Clean up logger at the end

	// Initialize Kubernetes client
	var kubClient = initKubernetes()

	// Initialize PostgreSQL client
	var dbClient = initDB()

	// Create Event handlers and related components
	var eventRepo = repositories.NewEventRepository(dbClient)
	var eventUC = uc.NewEventUC(eventRepo)
	var eventHandler = controller.NewEventHandler(eventUC)

	// Create Pod handlers and related components
	var podRepo = repositories.NewPodRepository(kubClient)
	var podUC = uc.NewPodUC(podRepo, eventUC)
	var podHandlers = controller.NewPodHandler(podUC)

	// Create Namespace handlers and related components
	var namespaceRepo = repositories.NewNamespaceRepository(kubClient)
	var namespaceUC = uc.NewNamespaceUC(namespaceRepo)
	var namespaceHandlers = controller.NewNamespaceHandler(namespaceUC)

	// Create Deployment handlers and related components
	var deploymentRepo = repositories.NewDeploymentInterfaces(kubClient)
	var deploymentUC = uc.NewDeploymentUC(deploymentRepo)
	var deploymentHandlers = controller.NewDeploymentHandler(deploymentUC)

	// Create user handlers and related components
	var userRepo = repositories.NewUserRepository(dbClient)
	var userUC = uc.NewUserUC(userRepo)
	var userHandlers = controller.NewUserHandlers(userUC)

	// Create Auth handlers and related components
	var authHandlers = controller.NewAuthHandlers(userUC)

	// Define authentication routes and handlers
	var authRoutes = e.Group("/auth")
	authRoutes.POST("/login", authHandlers.Login)

	var googleAuthHandler = controller.NewGoogleAuthHandler(userUC)
	var oauthRoutes = authRoutes.Group("")
	oauthRoutes.GET("/google_login", googleAuthHandler.GoogleLogin)
	oauthRoutes.GET("/google_callback", googleAuthHandler.GoogleCallback)

	var githubAuthHandler = controller.NewGithubAuthHandler(userUC)
	oauthRoutes.GET("/github_login", githubAuthHandler.GithubLogin)
	oauthRoutes.GET("/github_callback", githubAuthHandler.GithubCallback)

	// Add JWT authentication and authorization middleware
	var restrictedRoutes = e.Group("")
	restrictedRoutes.Use(util.JWTAuth)
	restrictedRoutes.Use(util.JWTAuthViewer)

	// Define user routes
	var usersRoutes = restrictedRoutes.Group("/users")
	usersRoutes.GET("", userHandlers.List)
	usersRoutes.GET("/:id", userHandlers.GetByID)
	usersRoutes.POST("", userHandlers.CreateUser)
	usersRoutes.PUT("/:id", userHandlers.UpdateUser)
	usersRoutes.DELETE("/:id", userHandlers.DeleteUser)

	// Define pod routes
	var podsRoutes = restrictedRoutes.Group("/pods")
	podsRoutes.GET("", podHandlers.List)
	podsRoutes.GET("/:id", podHandlers.GetByNameOrUID)
	podsRoutes.POST("", podHandlers.Create)
	podsRoutes.PUT("/:id", podHandlers.Update)
	podsRoutes.DELETE("/:id", podHandlers.Delete)

	// Define namespace routes
	var namespacesRoutes = restrictedRoutes.Group("/namespaces")
	namespacesRoutes.GET("", namespaceHandlers.List)
	namespacesRoutes.GET("/:id", namespaceHandlers.GetByNameOrUID)
	namespacesRoutes.POST("", namespaceHandlers.Create)
	namespacesRoutes.PUT("/:id", namespaceHandlers.Update)
	namespacesRoutes.DELETE("/:id", namespaceHandlers.Delete)

	// Define deployment routes
	var deploymentsRoutes = restrictedRoutes.Group("/deployments")
	deploymentsRoutes.GET("", deploymentHandlers.List)
	deploymentsRoutes.GET("/:id", deploymentHandlers.GetByNameOrUID)
	deploymentsRoutes.POST("", deploymentHandlers.Create)
	deploymentsRoutes.PUT("/:id", deploymentHandlers.Update)
	deploymentsRoutes.DELETE("/:id", deploymentHandlers.Delete)

	// Define event routes
	var eventsRoutes = restrictedRoutes.Group("/events")
	eventsRoutes.GET("", eventHandler.List)
	eventsRoutes.GET("/:id", eventHandler.GetByID)

	// Start the Echo application
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", viper.GetInt("api_service.port"))))
}

// Configures the Echo instance
func configureEcho(e *echo.Echo) {
	e.HideBanner = true
	e.HidePort = true

	// Add Swagger documentation route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Add Recover middleware
	e.Use(middleware.Recover())
}

// Configures CORS settings
func configureCORS(e *echo.Echo) {
	var corsConfig = middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{viper.GetString("ui_service.allow_origin")},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	})

	e.Use(corsConfig)
}

// Configures the logger and adds it as middleware
func configureLogger(e *echo.Echo) *zap.SugaredLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	e.Use(pkg.ZapLogger(logger))

	var sugar = logger.Sugar()
	var loggerHandler = controller.NewLogger(sugar)
	e.Use(loggerHandler.LoggerMiddleware)

	return sugar
}

// Initializes the Kubernetes client
func initKubernetes() *kubernetes.Clientset {
	client, err := pkg.NewKubernetesClient()
	if err != nil {
		log.Fatalf("Failed to initialize Kubernetes client: %v", err)
	}

	log.Println("Kubernetes client initialized successfully")
	return client
}

// Initializes the PostgreSQL client
func initDB() *pg.DB {
	var db = pkg.NewPSQLClient()
	if db == nil {
		log.Fatal("Failed to initialize PostgreSQL client")
	}

	log.Println("PostgreSQL client initialized successfully")
	return db
}
